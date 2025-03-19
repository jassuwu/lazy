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

const helpMessage = `lazyenv: cli to set .env files with latest contract address vars in P2P.me.
  has to be configured with a JSON config file like 'lazyenv.config.json'.

 Usage: lazyenv <command>

 Commands:
   run    Execute the given command in the source directory and then do the copy.
   copy   Copy environment variables from 'contract-addresses.json' to the '.env' files in the destination

 Example:
   lazyenv run
   lazyenv copy`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(helpMessage)
		os.Exit(0)
	}
	buf, err := os.ReadFile("lazyenv.config.json")
	assert.Nil(err, "config file couldn't be opened.")

	config := LazyEnvConfig{}
	err = json.Unmarshal(buf, &config)
	assert.Nil(err, "couldn't unmarshal config json")
	fmt.Println("- read the configuration")

	srcCmdArgs := strings.Split(config.Src.Cmd, " ")
	if os.Args[1] == "run" {
		cmd := exec.Command(srcCmdArgs[0], srcCmdArgs[1:]...)
		cmd.Dir = utils.ExpandTilde(config.Src.Dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		assert.Nil(err, "errored when running the command in the src directory")
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
		newFileLines = append(newFileLines, "# "+utils.ExtractArg(srcCmdArgs, "--network"))
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
