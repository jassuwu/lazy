package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jassuwu/p2pme/internal/assert"
	"github.com/jassuwu/p2pme/internal/flow"
	"github.com/jassuwu/p2pme/internal/utils"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type P2PMEConfigSrc struct {
	Dir      string `json:"dir"`
	FileName string `json:"fileName"`
	Cmd      string `json:"cmd"`
}

type P2PMEConfigDest struct {
	Paths      []string          `json:"paths"`
	EnvMapping map[string]string `json:"envMapping"`
}

type P2PMEConfig struct {
	Src  P2PMEConfigSrc  `json:"src"`
	Dest P2PMEConfigDest `json:"dest"`
}

var validCommands = map[string]bool{
	"run":     true,
	"copy":    true,
	"help":    true,
	"drycopy": true,
}

var validOptions = map[string]bool{
	"config": true,
}

var defaultConfigPath = utils.ExpandTilde("~/.config/p2pme.config.json")

const helpMessage = `p2pme: cli tool for syncing contract addresses to your .env files

  yo, this tool makes life easy - it grabs contract addresses from deployment
  and updates all your .env files automatically. pretty neat, right?

usage:
  p2pme <command> [options]

commands:
  run     run the source command and update all your .env files in one go
  copy    just update the .env files without running any commands
  drycopy just print the changes that would be made to the .env files
  help    show this message (you're looking at it now)

options:
  --config string   path to the config file (default "~/.config/p2pme.config.json")

examples:
  p2pme run                                # do everything in one shot
  p2pme copy                               # just update the env files
  p2pme drycopy                            # print the changes that would be made to the .env files
  p2pme help                               # what you're reading right now
  p2pme run --config ./my-config.json      # use a custom config file

config file (config.json):
  {
    "src": {
      "dir": "~/path/to/source",       # where to find your stuff
      "fileName": "addresses.json",    # the file with your contract addresses
      "cmd": "command to run"          # what command to run before copying
    },
    "dest": {
      "paths": ["~/path/to/.env"],     # env files to update
      "envMapping": {                  # how to map keys to env vars
        "sourceKey": "ENV_VAR_NAME"
      }
    }
  }`

func executeCommand(cmdString string) {
	flow.Action("running command")

	cmdArgs := strings.Split(cmdString, " ")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	assert.Nil(err, "errored when running the command")
	flow.Success("command completed")
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Println(helpMessage)
		os.Exit(0)
	}

	if !validCommands[os.Args[1]] {
		flow.Error(fmt.Sprintf("unknown command '%s'", os.Args[1]))
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	if len(os.Args) > 2 && !validOptions[os.Args[2]] {
		flow.Error(fmt.Sprintf("unknown option '%s'", os.Args[2]))
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	flow.Section("configuration")
	var configPath string
	flag.StringVar(&configPath, "config", defaultConfigPath, "path to the config file")
	flag.Parse()

	buf, err := os.ReadFile(configPath)
	assert.Nil(err, "config file couldn't be opened.")

	config := P2PMEConfig{}
	err = json.Unmarshal(buf, &config)
	assert.Nil(err, "couldn't unmarshal config json")
	flow.Success("config loaded")

	// Execute the command if "run" was specified
	if os.Args[1] == "run" {
		flow.Section("source command")
		executeCommand(config.Src.Cmd)
	}

	// Copy operations for both "run" and "copy" commands
	flow.Section("syncing addresses")

	// Read source file
	filePath := utils.ExpandTilde(config.Src.Dir + "/" + config.Src.FileName)
	flow.Info(fmt.Sprintf("source: %s", filePath))

	src, err := os.ReadFile(filePath)
	assert.Nil(err, "src contract addresses file couldn't be opened.")

	contractAddresses := make(map[string]string)
	err = json.Unmarshal(src, &contractAddresses)
	assert.Nil(err, "couldn't unmarshal src json")
	flow.Success(fmt.Sprintf("%d addresses found", len(contractAddresses)))

	// Update destination files
	flow.Section("updating files")

	for _, path := range config.Dest.Paths {
		expandedPath := utils.ExpandTilde(path)
		flow.FileAction("read", expandedPath)

		dest, err := os.ReadFile(expandedPath)
		assert.Nil(err, "dest contract addresses file couldn't be opened.")

		newFileLines := make([]string, 0)

		entries := strings.Split(string(dest), "\n")
	ENTRIES:
		for _, entry := range entries {
			for _, v := range config.Dest.EnvMapping {
				if strings.HasPrefix(entry, v) {
					newFileLines = append(newFileLines, "# "+entry)
					continue ENTRIES
				}
			}
			newFileLines = append(newFileLines, entry)
		}

		// adding metadata and new env vars
		newFileLines = append(newFileLines, "")
		network := utils.ExtractArg(config.Src.Cmd, "--network")
		if network != "" {
			newFileLines = append(newFileLines, "# "+network)
		}

		timestamp := time.Now().Local().Format("02 Jan 2006 15:04:05")
		newFileLines = append(newFileLines, "# updated: "+timestamp)

		for k, v := range config.Dest.EnvMapping {
			if addr, ok := contractAddresses[k]; ok {
				newFileLines = append(newFileLines, v+"="+addr)
			} else {
				flow.Warn(fmt.Sprintf("missing address for key: %s", k))
			}
		}
		newFileLines = append(newFileLines, "")

		if os.Args[1] == "drycopy" {
			flow.FileAction("write", expandedPath)
			fmt.Println(strings.Join(newFileLines, "\n"))
			flow.Success("drycopy completed")
		} else {
			flow.FileAction("write", expandedPath)
			envFile, err := os.Create(expandedPath)
			assert.Nil(err, "couldn't 'CREATE' envFile.")
			defer envFile.Close()
			envFile.WriteString(strings.Join(newFileLines, "\n"))
			flow.Success(fmt.Sprintf("updated %s", path))
		}
	}

	flow.Done()
}
