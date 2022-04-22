package main

import (
	"fmt"
	"log"
	"path/filepath"

	conf_reader "github.com/bohdanhlovatskyi/goindex.git/conf_reader"
	iutil "github.com/bohdanhlovatskyi/goindex.git/indexation_util"
)

const CONFIG_PATH = "index.toml"

func main() {
	conf := conf_reader.Get_config(CONFIG_PATH)

	path_q := make(chan string, 100)
	data_q := make(chan iutil.RawFileSource, 100)

	go iutil.TraverseDirs(conf.Indir, path_q)
	go iutil.Reader(path_q, data_q)

	for elm := range data_q {
		var data string
		var err error
		rdata, path := elm.Data, elm.Path
		if filepath.Ext(path) == ".zip" {
			data, err = iutil.ProcessRawSource(rdata)
			if err != nil {
				log.Println("oculd not unzip data:", err)
			}
		} else {
			data = string(rdata)
		}

		fmt.Println(data)
	}
}
