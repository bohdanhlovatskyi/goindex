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
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go iutil.Indexer(data_q, merge_q, wg)
	}

	wg2 := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go iutil.Merger(merge_q, wg2)
	}

	wg.Wait()
	// drop poisson pill here for the merger set of functions
	merge_q <- make(map[string]int)

	elapsed := time.Since(start)
	fmt.Println("indexers end time: ", elapsed)

	wg2.Wait()
	if len(merge_q) != 2 {
		panic("mergers did not succed")
	}
	m := <-merge_q
	_ = <-merge_q // take the poisson pill out
	close(merge_q)

	elapsed = time.Since(start)
	fmt.Println("mergers end time: ", elapsed)
	sstart := time.Now()

	iutil.WriteMap(m, conf.Out_by_a, conf.Out_by_n)
	elapsed = time.Since(sstart)
	fmt.Println("write taken: ", elapsed)

	elapsed = time.Since(start)
	fmt.Println("Time taken: ", elapsed)
}
