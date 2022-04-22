package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/bohdanhlovatskyi/goindex.git/conf_reader"
	iutil "github.com/bohdanhlovatskyi/goindex.git/indexation_util"
)

const CONFIG_PATH = "index.toml"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	conf := conf_reader.Get_config(CONFIG_PATH)

	path_q := make(chan string, 1000)
	data_q := make(chan iutil.RawFileSource, 1000)
	merge_q := make(chan map[string]int, 1000)

	start := time.Now()
	go iutil.TraverseDirs(conf.Indir, path_q)
	go iutil.Reader(path_q, data_q)

	wg := &sync.WaitGroup{}

	for i := 0; i < 65; i++ {
		wg.Add(1)
		go iutil.Indexer(data_q, merge_q, wg)
	}

	wg.Wait()
	close(merge_q)

	elapsed := time.Since(start)
	fmt.Println("indexers taken: ", elapsed)

	// merge the maps
	mc := make(chan map[string]int)
	go func() {
		m := make(map[string]int)
		for lm := range merge_q {
			for k, v := range lm {
				m[k] += v
			}
		}

		mc <- m
	}()

	m := <-mc

	iutil.WriteMap(m, conf.Out_by_a, conf.Out_by_n)
	elapsed = time.Since(start)
	fmt.Println("Time taken: ", elapsed)
}
