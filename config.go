package main

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Delegators []string            `yaml:"Delegators"`
	Bakers     []string            `yaml:"Bakers"`
	Whitelist  map[string][]string `yaml:"Whitelist"`
	Aliases    map[string]string   `yaml:"Aliases"`
}

func loadConfig(file string) *config {
	c := config{}

	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln("Unable to read file: ", file, err)
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalln("Unable to parse yaml file: ", file, err)
	}
	return &c
}
