package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/jassuwu/lazyenv/internal/assert"
	"github.com/jassuwu/lazyenv/internal/utils"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type LEConfigSrc struct {
	Dir      string `json:"dir"`
	FileName string `json:"fileName"`
	Cmd      string `json:"cmd"`
}

type LEConfigDest struct {
	Paths      []string          `json:"paths"`
	EnvMapping map[string]string `json:"envMapping"`
}

type LazyEnvConfig struct {
	Src  LEConfigSrc  `json:"src"`
	Dest LEConfigDest `json:"dest"`
}

var validCommands = map[string]bool{
	"run":  true,
	"copy": true,
	"help": true,
}

const helpMessage = `lazyenv: chill, dumb, fast cli tool for syncing contract addresses to your .env files

  yo, this tool makes life easy - it grabs contract addresses from deployment
  and updates all your .env files automatically. pretty neat, right?
  just drop a 'lazyenv.config.json' in your current directory and we're good to go.

usage:
  lazyenv <command>

commands:
  run     run the source command and update all your .env files in one go
  copy    just update the .env files without running any commands
  help    show this message (you're looking at it now)

examples:
  lazyenv run             # do everything in one shot
  lazyenv copy            # just update the env files
  lazyenv help            # what you're reading right now

config file (lazyenv.config.json):
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
	cmdArgs := strings.Split(cmdString, " ")
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	assert.Nil(err, "errored when running the command")
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		fmt.Println(helpMessage)
		os.Exit(0)
	}
	if !validCommands[os.Args[1]] {
		fmt.Printf("Error: unknown command '%s'\n\n", os.Args[1])
		fmt.Println(helpMessage)
		os.Exit(1)
	}

	buf, err := os.ReadFile("lazyenv.config.json")
	assert.Nil(err, "config file couldn't be opened.")

	config := LazyEnvConfig{}
	err = json.Unmarshal(buf, &config)
	assert.Nil(err, "couldn't unmarshal config json")
	fmt.Println("- read the configuration")

	if os.Args[1] == "run" {
		executeCommand(config.Src.Cmd)
	}

	// Go doesn't understand tilde. So need to expand just in case.
	src, err := os.ReadFile(utils.ExpandTilde(config.Src.Dir + "/" + config.Src.FileName))
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

		newFileLines := make([]string, 0)

		entries := strings.Split(string(dest), "\n") // ["FOO=BAR", "BOO=BAZ", ""]
	ENTRIES:
		for _, entry := range entries {
			// commenting out the existing matching envVars
			for _, v := range config.Dest.EnvMapping {
				if strings.HasPrefix(entry, v) {
					newFileLines = append(newFileLines, "# "+entry)
					continue ENTRIES
				}
			}
			newFileLines = append(newFileLines, entry)
		}

		// adding one line for each of the new envVar
		newFileLines = append(newFileLines, "")
		newFileLines = append(newFileLines, "# "+utils.ExtractArg(config.Src.Cmd, "--network"))
		newFileLines = append(
			newFileLines,
			"# "+
				time.Now().Local().Format("02 January, 2006 15:04:05 IST"),
		)
		for k, v := range config.Dest.EnvMapping {
			newFileLines = append(newFileLines, v+"="+contractAddresses[k])
		}
		newFileLines = append(newFileLines, "")

		envFile, err := os.Create(utils.ExpandTilde(path))
		assert.Nil(err, "couldn't 'CREATE' envFile.")
		defer envFile.Close()
		envFile.WriteString(strings.Join(newFileLines, "\n"))
	}
}
