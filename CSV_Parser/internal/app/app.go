package app

import (
	"github.com/tumbleweedd/intership/CSV_Parser/internal/csv"
	"os"
	"path/filepath"
	"sync"
)

const (
	rabbitDSN = "amqp://guest:guest@localhost:5672/"
	rootDir   = "../data"
)

func Run() {
	files, err := filePathWalkDir(rootDir)
	if err != nil {
		panic(err)
	}

	wg := new(sync.WaitGroup)

	done := make(chan struct{})
	defer close(done)

	for _, file := range files {
		wg.Add(1)
		go csv.ParseCSV(file, rabbitDSN, wg, done)
	}

	wg.Wait()
}

func filePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
