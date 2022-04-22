package conf_reader

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

type Config struct {
	Indir            string
	Out_by_a         string
	Out_by_n         string
	Indexing_threads int
	Merging_threads  int
}

func Get_config(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("could not open file with config")
	}

	conf := Config{}
	toml.Unmarshal(data, &conf)

	valiadte_conf(conf)

	return conf
}

func valiadte_conf(conf Config) {
	_, err := os.Stat(conf.Indir)
	if os.IsNotExist(err) {
		log.Fatal("Indir file does not exist.")
	}

	if conf.Indexing_threads < 0 || conf.Merging_threads < 0 {
		log.Fatal("Invalid range of threds")
	}
}
