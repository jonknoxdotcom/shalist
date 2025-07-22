# shalist
SHA256 file manager *WIP*

This is the WIP version before integration with cobra.

## File format

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
