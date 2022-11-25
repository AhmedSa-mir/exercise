package shred

import (
	"log"
	"os"
	"io/fs"
	"math"
	"fmt"

	"example/pkg/fileutils"
)

type shredConfig struct {
	Iterations int
	BufMaxSize int64
}

// Shred overwrites a file at the given path 3 times and then deletes it.
func Shred(path string) error { 
	ok := fs.ValidPath(path)
	if !ok {
		return fmt.Errorf("error: invalid path: %s. Path must not contain a '.' or '..' or the empty string or start/end with a slash.", path)
	}

	var invalidModes fs.FileMode = fs.ModeNamedPipe |
					 fs.ModeSocket | fs.ModeDevice | fs.ModeCharDevice |
					 fs.ModeIrregular
	invalid, err := fileutils.CheckFileType(path, invalidModes)
	if err != nil {
		return fmt.Errorf("file error: %s", err)
	}
	if invalid {
		return fmt.Errorf("shred: error: invalid file type")
	}

	f, fileSize, err := fileutils.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error opening file: %s", err)
	}
    defer func() {
        if err := f.Close(); err != nil {
            log.Fatalf("error closing file: %s", err)
        }
    }()

	cfg := shredConfig {
		Iterations: 3,
		BufMaxSize: math.MaxUint32,
	}

	for i := 0; i < cfg.Iterations; i++ {
		log.Printf("shred: iteration #%d", i+1)
		var idx int64 = 0
		var bufSize int64 = cfg.BufMaxSize
		var chunksCount int64 = fileSize/cfg.BufMaxSize
		for i := int64(0); i <= chunksCount; i++ {
			// last chunk (remaining bytes)
			if i == int64(chunksCount) {
				bufSize = fileSize - i
			}
			err := fileutils.WriteRandomBytes(f, bufSize, idx)
			if err != nil {
				return fmt.Errorf("error writing to file: %s", err)
			}
			idx += cfg.BufMaxSize
		}
	}

	err = os.Remove(path)
    if err != nil {
		return fmt.Errorf("error deleting file: %s", err)
    }

	log.Printf("shred: success: %s has been shredded and deleted.", path)
	return nil
}
