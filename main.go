package main

import (
	"encoding/json"
	"flag"
	"log"
	"monitoring/checks"
	"monitoring/spec"
	"os"
	"sync"
)

func main() {
	file := flag.String("config", "./config-sample.json", "Configuration file")
	flag.Parse()

	content, err := os.ReadFile(*file)
	if err != nil {
		log.Panic(err)
	}

	checksFile := spec.CheckFile{}
	err = json.Unmarshal(content, &checksFile)
	if err != nil {
		log.Panic(err)
	}

	wg := sync.WaitGroup{}
	for _, v := range checksFile.Checks {
		handler := checks.NewHandler(v)
		go handler.Run()
		wg.Add(1)
	}

	wg.Wait() // indefinitely
}
