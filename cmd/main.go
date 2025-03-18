package main

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
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

var Config = &LazyEnvConfig{
	Src: LEConfigSrc{
		Path: "~/repos/p2pdotme/contracts-v3",
		Cmd:  "bunx --bun hardhat run ./scripts/deployContracts.ts --network baseSepolia",
	},
	Dest: LEConfigDest{
		Paths: []string{
			"~/repos/p2pdotme/broker-ui",
			"~/repos/p2pdotme/user-app",
			"~/repos/p2pdotme/merchant-app",
		},
		EnvMapping: map[string]string{
			"usdc":              "VITE_USDT_CONTRACT_ADDRESS",
			"config":            "VITE_CONFIG_CONTRACT",
			"merchantRegistry":  "VITE_MERCHANT_REGISTRY",
			"orderProcessor":    "VITE_ORDER_PROCESSOR",
			"reputationManager": "VITE_REPUTATION_MANAGER",
		},
	},
}

// NEXT: Setup assertions to panic incase of skill issues

func main() {
	s, err := json.MarshalToString(Config)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(s)
	}

	fmt.Println("---")

	config := &LazyEnvConfig{}
	err = json.UnmarshalFromString(s, config)
	if err != nil {
		panic(err)
	} else {
		fmt.Println(config.Src.Path)
	}
}
