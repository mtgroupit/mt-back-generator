package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mtgroupit/mt-back-generator/internal/gen"
	"github.com/mtgroupit/mt-back-generator/internal/parser"
)

var (
	dir    string
	config string
)

const (
	defaultCfgFileName = "config.yaml"
)

func init() {
	flag.StringVar(&dir, "dir", "./generated/", "the direrctory to generate the service")
	flag.StringVar(&config, "config", defaultCfgFileName, "the config file to use")
}

func main() {
	flag.Parse()

	cfg, err := parser.ReadYAMLCfg(config)
	if err != nil {
		log.Fatal(err)
		return
	}

	cfg, err = parser.HandleCfg(cfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	prevCfgPath := filepath.Join(dir, cfg.Name, defaultCfgFileName)

	// read previous version config if exists
	prevCfg, err := parser.ReadYAMLCfg(prevCfgPath)
	if err == nil {
		prevCfg, err = parser.HandleCfg(prevCfg)
		if err != nil {
			log.Fatal(err)
			return
		}
	}

	err = gen.Srv(dir, cfg, prevCfg)
	if err != nil {
		log.Fatal(err)
		return
	}

	// save configuration for futher version generation
	_, err = copyFile(config, prevCfgPath)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
