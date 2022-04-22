package indexation_util

import (
	"log"
	"os"
	"path/filepath"
)

func TraverseDirs(path string, path_q chan string) {
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() && info.Size() != 0 {
				path_q <- path
			}

			return nil
		})

	if err != nil {
		log.Println(err)
	}

	path_q <- ""
	close(path_q)
}
