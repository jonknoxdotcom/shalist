package main

// in-memory version

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

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

func WalkTree(startpath string) (int64, error) {
	entries, err := os.ReadDir(startpath)
	if err != nil {
		return 0, err
	}
	var total int64
	for index, entry := range entries {
		if !entry.IsDir() {
			fmt.Println(index, path.Join(startpath, entry.Name()))
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
	size, _ := WalkTree(".")
	fmt.Println(size)

	// This directory reader uses the new os.ReadDir (req 1.16)
	// https://benhoyt.com/writings/go-readdir/

}
