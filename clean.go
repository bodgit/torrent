package torrent

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
)

func (t *Torrent) Clean(base string, logger *log.Logger, dryRun bool) (int, error) {
	if len(t.Info.Files) == 0 {
		return 0, errors.New("single file torrent")
	}

	dir := path.Join(base, t.Info.Name)

	if _, err := os.Stat(dir); err != nil {
		return 0, err
	}

	files := make(map[string]struct{})

	for _, file := range t.Info.Files {
		files[path.Join(append([]string{dir}, file.Path...)...)] = struct{}{}
	}

	deleted := 0

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if _, ok := files[path]; !ok {
			logger.Println("Deleting", path)
			if !dryRun {
				if err := os.Remove(path); err != nil {
					return err
				}
			}
			deleted++
		}
		return nil
	}); err != nil {
		return 0, err
	}

	return deleted, nil
}
