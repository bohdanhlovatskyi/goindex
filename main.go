package main

import (
	"fmt"
	"sync"

	"github.com/bohdanhlovatskyi/goindex.git/conf_reader"
	iutil "github.com/bohdanhlovatskyi/goindex.git/indexation_util"
)

const CONFIG_PATH = "index.toml"

func main() {
	conf := conf_reader.Get_config(CONFIG_PATH)

	path_q := make(chan string, 100)
	data_q := make(chan iutil.RawFileSource, 100)
	merge_q := make(chan map[string]int, 100)

	go iutil.TraverseDirs(conf.Indir, path_q)
	go iutil.Reader(path_q, data_q)

	wg := &sync.WaitGroup{}

	for i := 0; i < conf.Indexing_threads; i++ {
		wg.Add(1)
		go iutil.Indexer(data_q, merge_q, wg)
	}

	wg.Wait()
	close(merge_q)

	for elm := range merge_q {
		fmt.Println(elm)
	}
}
