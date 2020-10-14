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
	flag.StringVar(&dir, "dir", "./generated/", "the direrctory to generate the service")
	flag.StringVar(&config, "config", "./config.yaml", "the config file to use")
}

func main() {
	flag.Parse()

	cfg, err := parser.ReadYAMLCfg(config)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err = parser.HandleCfg(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = gen.Srv(dir, cfg)
	if err != nil {
		log.Fatal(err)
	}
}
