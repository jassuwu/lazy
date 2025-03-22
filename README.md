# lazy (OLD && OUTDATED README)

chill, dumb, fast cli for syncing contract addresses non destructively to your `.env` files.

## what is this?

lazyenv makes life easy - it grabs contract addresses from deployment and updates all your .env files automatically. no more manual copy-pasting. it's completely non-destructive - your existing env vars stay commented while new ones are added, with metadata.

## requirements

- go installed on your machine, and the GOPATH and GOBIN properly setup.
- **important**: your contract deployment process must save addresses to a `contract-addresses.json` file (or whatever JSON you specify in config)

## install

```bash
go install github.com/jassuwu/lazyenv@latest
```

## features

- automatic environment variable syncing
- runs your deployment command for you
- updates multiple .env files at once
- simple configuration
- **non-destructive** - comments out old values rather than deleting them

## usage

```bash
lazyenv <command> [options]
```

### commands

- `run` - run the source command and update all your .env files in one go
- `copy` - just update the .env files without running any commands
- `drycopy` - show what changes would be made without actually making them
- `help` - show help message

### options

- `--config string` - path to config file (default "~/.config/lazyenv/config.json")

## examples

```bash
lazyenv run                              # do everything in one shot
lazyenv copy                             # just update the env files
lazyenv drycopy                          # preview changes
lazyenv help                             # show help message
lazyenv run --config ./my-config.json    # use custom config
```

## configuration

drop a config file at `~/.config/lazyenv/config.json`:

```json
{
  "src": {
    "dir": "~/path/to/source",       // where to find your stuff
    "fileName": "addresses.json",    // file with contract addresses
    "cmd": "command to run"          // command to run before copying
  },
  "dest": {
    "paths": ["~/path/to/.env"],     // env files to update
    "envMapping": {                  // how to map keys to env vars
      "sourceKey": "ENV_VAR_NAME"
    }
  }
}
```

check out [example.config.json](./example.config.json) for a complete example.

## workflow

1. set up your config file
2. run `lazyenv run`
3. chill while your contract addresses are synced
4. profit

## why?

ok yeah, the manual process is error prone, slower, blah blah blah. but why? really i just wanted wanted to write `go` cli