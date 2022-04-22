package main

import (
	"fmt"

	conf_reader "github.com/bohdanhlovatskyi/goindex.git/conf_reader"
)

func main() {
	const CONFIG_PATH = "index.toml"

	conf := conf_reader.Get_config(CONFIG_PATH)
	fmt.Println(conf)
}
