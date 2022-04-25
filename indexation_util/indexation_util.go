package indexation_util

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

func processRawSource(raw []byte) (string, error) {
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

func Indexer(data_q chan RawFileSource, merger_q chan map[string]int, wg *sync.WaitGroup) {

	var data string
	var err error

	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	c := cases.Fold()

	for elm := range data_q {
		m := make(map[string]int)

		if filepath.Ext(elm.Path) == ".zip" {
			data, err = processRawSource(elm.Data)
			if err != nil {
				log.Println("oculd not unzip data:", err)
			}
		} else {
			data = string(elm.Data)
		}

		data, _, _ := transform.String(t, data)
		data = c.String(data)
		for _, word := range strings.Fields(data) {
			m[word] += 1
		}

		merger_q <- m
	}

	wg.Done()
}

// this one is reeeeally not in the nature of go :(, though
// I am to stupid to do this in the right way
func Merger(q chan map[string]int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		f := <-q
		if len(f) == 0 {
			q <- f
			if len(q) <= 2 {
				break
			}
			continue
		}

		s := <-q
		if len(s) == 0 {
			q <- f
			q <- s
			if len(q) <= 2 {
				break
			}
			continue
		}

		if len(f) > len(s) {
			for k, v := range s {
				f[k] += v
			}
			q <- f
		} else {
			for k, v := range f {
				s[k] += v
			}
			q <- s
		}
	}
}

func WriteMap(res map[string]int, res_a, res_n string) {
	keys := make([]string, 0, len(res))
	for k := range res {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	f, err := os.Create(res_a)
	if err != nil {
		log.Fatal("could not create a file")
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, k := range keys {
		w.Write([]byte(fmt.Sprintln(k, res[k])))
	}
	w.Flush()

	sort.Slice(keys, func(i, j int) bool { return res[keys[i]] > res[keys[j]] })

	ff, err := os.Create(res_n)
	if err != nil {
		log.Fatal("could not create a file")
	}
	defer ff.Close()

	w = bufio.NewWriter(ff)
	for _, k := range keys {
		w.Write([]byte(fmt.Sprintln(k, res[k])))
	}
	w.Flush()
}
