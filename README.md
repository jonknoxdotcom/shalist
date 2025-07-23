# shalist
SHA256 file manager *WIP*

Tool for handing assets in a verifiable manner as part of a broader management strategy. Can be used to de-clutter filespaces, and - as part of a security process - be used to check for sensitive data spillage.

This is the WIP version before integration with cobra.

## What it does
* Maintains a signature file describing the list of files in a file tree based on hash (sha256)
* Allows updates to the signature file at low computational cost
* Allows comparisons of different trees via their signature files
* Produces `bash`-style commands to delete duplicated data and directory structures

## Is and Isn't
Is:
* `shalist` is used to manage/delete duplicate files
* `shalist` is used to generate a cryptographic record of files
* `shalist` is a specialised bill-of-materials generator

Isn't:
* `shalist` isn't an archiving or backup utility
* `shalist` isn't a file copying tool


## Uses
* Assert or validate assets present (cryptographic bill of materials)
* Use to clear 'dead' disks without fear of erasing unique data
* Identify and remove duplicate data
* Detect spillage of files from one tree to another

## Case studies
1. Migrate from one work machine to another - desire to clear up old machine without data loss
2. Use of memory stick or burner machine containing a subset of working files - need to consolidate any changes files
3. Hard disk has duplicate copies of files in different directories - want to rationalise without loss
4. Clean up unstructured or chaotic snapshots of different machines across different media without loss
5. Generate a cryptographically verified list of files present on a target machine/disk

## Operation
The command consists of a verb and none or more optional arguments.

There are fundamentally three types of operations:
* generating and maintaining signature files (in presence of the source material)
``` 
shalist generate
shalist generate -p /volume4
shalist generate -p /volume4/
shalist generate -p accounts/ -p receipts/ -p invoices/
shalist generate new.jsf
shalist update existing.jsf
shalist update existing.jsf -a P
shalist update existing.jsf -a G
shalist update existing.jsf -a K
shalist verify existing.jsf
shalist verify existing.jsf -h -m -s
```

* splicing and dicing files from a signature file into smaller ones, or combining signature files
```
shalist extract bigfile.jsf remtree.jsf "REMTREE/"
shalist crop remtree.jsf "REMTREE/"
shalist extract bigfile.jsf remtree.jsf "REMTREE/" -crop
shalist graft bigfile.jsf subtree.jsf "SUBTREE/"
```

* analyze scripts, or generate scripts containing `bash`-style command to allow deletion of duplicate data
```
shalist info file.jsf
shalist csv file.jsf
shalist biggest file.jsf
shalist biggest file.jsf -n 20
shalist find file.jsf e8faee25618bc95b5954196ba7f2a3251c04b9cc12394cf7eec545bbc2c15a4d
shalist find file.jsf 6PruJWGLyVtZVBlrp/KjJRwEucwSOUz37sVFu8LBWk0
shalist whereis file.jsf wanted.doc
shalist duplicates file.jsf
shalist duplicates file.jsf -rm
shalist duplicates file.jsf -rm -n 10
shalist compare main.jsf lesser.jsf -rm -rd
```

Every command name can be shortened to 3-letters (i.e. `gen`, `upd`, `big`, `dup`...).

## Detailed command descriptions

### 1. Generate - creating new JSF file
Produces a JSF file (to STDOUT) for the current working directory.  Short form `gen`. Can use `--path` or `-p` to select a path other than cwd.  Optional filename for output.
```
shalist generate
shalist generate > myfiles.jsf
shalist generate myfiles.jsf
shalist generate -p Desktop/
shalist gen -p /mnt/thumbdrive oldthumb.jsf
```

The following example specifies multiple paths that will be form the signature file:
```
shalist generate -p accounts/ -p receipts/ -p invoices/ fin.jsf
```
In this case, the paths and collected, sorted, then indexed one by one.  So the single composite output JSF file `fin.jsf` will contain `accounts/...` then `invoices/...` then `receipts/...` records.


### 2. Update an existing JSF file

```
shalist update file.jsf
```
The command `update` can be shortened to `upd`.

### 3. Compare
```
shalist compare
```

### 4. Describe

### 5. Difference

### 6. Extract - remove subtree from signature file

### 7. Merge - merge two signature files

### 8. Remove - remove known files from target (bash script)

### 9. Duplicates - find all duplicate files (bash script)

### 10. 

## File format
* JSF files are line-per-file collections of file descriptions
* Each line contain identifying information consisting of file hash, last modify time/date, and size
* They are in strict ASCII (byte) order of the filename element.  This corresponds to locale specification `LC_COLLATE=C `.
* The specification allow the insertion of extra metadata called annotations between the identification block and filename

### SHA part:  (43x b64 ch)

* SHA256 generates a 256-bit hash - being 8x 32-bit words (i.e. 32 bytes)
* In hexadecimal, this would be represented by 64 characters (of 0-9,a-f).
* As base64 uses a radix-64 encoder (6 bits), SHA256 can be represented by 43 chars.
* A base64 encoding of a SHA256 would be 43ch followed by a '='.
* This trailing '=' is ommitted in the file (only 43 chars stored).

### Epoch time part: (8x hex ch)

* The epoch time is a second-resolution file modify time stored as 4 bytes.
* This is stored as 8 hex characters, probably beginning 68 or 69 (in 2025).

### File size part: (4+ hex ch)

* The file size is store in hex, with a minimum length of 4 hexadecimal chars.
* For a simple JSF file, this tends to make all filenames for <64k files line up.
* This provides a visual cue for visual reading of the file to find large files. 

### Annotations

* None or more annotation records.
* Annotation records contain no spaces and do not begin with ':'.

### Filename (to EOLN)
* Filename, prefixed by a ':'.
* Embedded control characters represented in hex, e.g. `\0x0d`
* Backslash represented by `\\`.
* All other characters (including UTF8) in plaintext.
