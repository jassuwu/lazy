package main

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"

	"github.com/jassuwu/lazyenv/internal/assert"
	"github.com/jassuwu/lazyenv/internal/utils"
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
	assert.Nil(err, "couldn't unmarshal config json")
	fmt.Println("- read the configuration")

	// Go doesn't understand tilde. So need to expand just in case.
	src, err := os.ReadFile(utils.ExpandTilde(config.Src.Path))
	assert.Nil(err, "src contract addresses file couldn't be opened.")

	contractAddresses := make(map[string]string)
	err = json.Unmarshal(src, &contractAddresses)
	assert.Nil(err, "couldn't unmarshal src json")
	fmt.Println("- read the src contract-addresses.json")

	fmt.Println("- reading the destination paths:")
	for _, path := range config.Dest.Paths {
		fmt.Println("-- reading:", path)
		dest, err := os.ReadFile(utils.ExpandTilde(path))
		assert.Nil(err, "dest contract addresses file couldn't be opened.")

		fmt.Println(string(dest))
	}
}
