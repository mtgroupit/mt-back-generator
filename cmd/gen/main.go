package main

import (
	"flag"
	"log"

	"github.com/mtgroupit/mt-back-generator/internal/gen"
	"github.com/mtgroupit/mt-back-generator/internal/parser"
)

var (
	dir    string
	config string
)

func init() {
	flag.StringVar(&dir, "dir", "./", "dir for generated servicce")
	flag.StringVar(&config, "config", "./config.yaml", "config file")
}

func main() {
	flag.Parse()

	cfg, err := parser.Cfg(config)
	if err != nil {
		log.Fatal(err)
	}

	err = gen.Srv(dir, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
