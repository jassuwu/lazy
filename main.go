package main

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"

	assert "github.com/jassuwu/lazyenv/internal/assert"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type LEConfigSrc struct {
	Path string `json:"path"`
	Cmd  string `json:"cmd"`
}

type LEConfigDest struct {
	Paths      []string          `json:"paths"`
	EnvMapping map[string]string `json:"envMapping"`
}

type LazyEnvConfig struct {
	Src  LEConfigSrc  `json:"src"`
	Dest LEConfigDest `json:"dest"`
}

func main() {
	buf, err := os.ReadFile("lazyenv.config.json")
	assert.Nil(err, "config file couldn't be opened.")

	config := LazyEnvConfig{}
	err = json.Unmarshal(buf, &config)
	assert.Nil(err, "couldn't unmarshal json")

	fmt.Println("LazyEnvConfig: ", config.Dest.Paths)
}
