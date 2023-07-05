package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Chernichniy/zkpCourses/internal/app/apiserver"
)

var (
	configPATH string
)

func init() {
	flag.StringVar(&configPATH, "config-path", "config/apiserver.toml", "path to config file")
}

func main() {
	flag.Parse()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPATH, config)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
