package main

import (
	"fmt"

	conf_reader "github.com/bohdanhlovatskyi/goindex.git/conf_reader"
	iutil "github.com/bohdanhlovatskyi/goindex.git/indexation_util"
)

const CONFIG_PATH = "index.toml"

func main() {
	conf := conf_reader.Get_config(CONFIG_PATH)

	path_q := make(chan string, 100)

	go iutil.TraverseDirs(conf.Indir, path_q)
	for path := range path_q {
		fmt.Println(path, len(path))
	}
}
