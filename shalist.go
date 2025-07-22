package main

// in-memory version

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"crypto/sha256"
	"io"

	b64 "encoding/base64"
)

var dupes = map[string]uint32{}

func reload(file *os.File) {
	fmt.Println("Reading existing file...")

	// junk for rest of fn

	// read each line and ensure it's a valid JSON document
	scanner := bufio.NewScanner(file)
	lines := 1
	for scanner.Scan() {
		if json.Valid([]byte(scanner.Text())) {
			fmt.Print(".")
		} else {
			fmt.Printf("(%d)", lines)
		}
		lines++
	}

	fmt.Println()
}

func GetTreeSize(startpath string) (int64, error) {
	entries, err := os.ReadDir(startpath)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, entry := range entries {
		if entry.IsDir() {
			size, err := GetTreeSize(path.Join(startpath, entry.Name()))
			if err != nil {
				return 0, err
			}
			total += size
		} else {
			info, err := entry.Info()
			if err != nil {
				return 0, err
			}
			total += info.Size()
		}
	}
	return total, nil
}

// Compute SHA256 for a given filename, returning byte array x 32
func GetSha256OfFile(fn string) ([]byte, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func abort(rc int, reason string) {
	fmt.Println(reason)
	os.Exit(rc)
}

func WalkTree(startpath string) (int64, error) {
	// uses the "new" (1.16) os.ReadPath functionality
	entries, err := os.ReadDir(startpath)
	if err != nil {
		abort(1, "Unrecoverable failure to read directory")
	}
	var total int64
	for _, entry := range entries {
		if !entry.IsDir() {
			// we ignore symlinks
			if !entry.Type().IsRegular() {
				continue
			}
			// emit file data
			name := path.Join(startpath, entry.Name())
			info, err := entry.Info()
			if err != nil {
				abort(2, "Internal error #2")
			}

			size := info.Size()
			unixtime := info.ModTime().Unix()
			// mode := info.Mode() // looks like '-rwxr-xr-x', alsoi synonymous to entry.Type().Perm()

			sha, _ := GetSha256OfFile(name)
			shab64 := b64.StdEncoding.EncodeToString(sha)
			if len(shab64) != 44 {
				abort(3, "Internal error #3: "+name)
			}
			if shab64[43:] != "=" {
				abort(4, "Internal error #4: "+name)
			} else {
				shab64 = shab64[0:43]
			}
			dupes[shab64] = dupes[shab64] + 1

			fmt.Printf("%s%x%04x :%s", shab64, unixtime, size, name)
			fmt.Println()
		}
		if entry.IsDir() {
			size, err := WalkTree(path.Join(startpath, entry.Name()))
			if err != nil {
				return 0, err
			}
			total += size
		} else {
			info, err := entry.Info()
			if err != nil {
				return 0, err
			}
			total += info.Size()
		}
	}
	return total, nil
}

func main() {

	// Current format ('sha file manager' files):
	// shalist ../existing.sfm  -- reads existing file, then indexes current dir
	// shalist                  -- indexes current dir

	// check for reload of existing file
	if len(os.Args) == 2 {
		filename := os.Args[1]
		fmt.Printf("SHAfile:   %s", filename)

		// open logfile (if possible)
		shafile, err := os.Open(filename)
		if err != nil {
			fmt.Println("Unable to open file " + filename)
			os.Exit(1)
		}
		defer shafile.Close()

		reload(shafile)
	}

	// Estimate size
	_, _ = WalkTree(".")

	// This directory reader uses the new os.ReadDir (req 1.16)
	// https://benhoyt.com/writings/go-readdir/

	for id, times := range dupes {
		if times > 1 {
			fmt.Println("# " + id)
		}
	}

}
