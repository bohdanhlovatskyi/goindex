package indexation_util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RawFileSource struct {
	Data []byte
	Path string
}

type FileSource struct {
	Data  string
	Psath string
}

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

func __read_binary(path string) ([]byte, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// calculate the bytes size
	var size int64 = info.Size()
	bytes := make([]byte, size)

	// read into buffer
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)

	return bytes, nil
}

func Reader(path_q chan string, data_q chan RawFileSource) {

	for path := range path_q {
		// poisson pill, this can be done slsighly less ugly
		// via some struct with context (method that cab be applied to it)
		// this will also allow to cancel work at any time
		if path == "" {
			data_q <- RawFileSource{[]byte(""), ""}
			break
		}

		data, err := __read_binary(path)
		if err != nil {
			log.Println("could not read file: ", path, err)
			continue
		}

		data_q <- RawFileSource{data, path}
	}

	close(data_q)
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func ProcessRawSource(raw []byte) (string, error) {
	b := bytes.NewReader(raw)
	reader, err := zip.NewReader(b, int64(len(raw)))
	if err != nil {
		return "", err
	}

	// Read all the files from zip archive
	var sb strings.Builder
	for _, zipFile := range reader.File {
		data, err := readZipFile(zipFile)
		if err != nil {
			log.Println(err)
			continue
		}

		sb.WriteString(string(data) + "\n")
	}

	return sb.String(), nil
}
