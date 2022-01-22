# Airtel Donge Automation

Automates and implements the following:

- View Stats of the dongel
- Notify if the dongel is low in charge
- Notify if the dongel is fully charged
- View SMS
- Delete SMS

## Requirements

- go >= 1.16


## Install

Follow the following steps in terminal

- `make dep`
- `cp .env.example .env`
- Update the .env file with relevent values
- `make build`

## Run

In order to run, either set the env variables in command line itself or update `.env` properly.

```
./dongel <command> <flags...> <inputs..>

Usage:
  dongel [flags]
  dongel [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  exec        Handle Few execs
  help        Help about any command
  login       Login to Dongel
  logout      Logout from Dongel
  sms         Check SMS
  stats       Get Stats

Flags:
  -a, --alert                 Alert for low battery
      --alert-charge string   Charge value below which alert will be triggered
  -d, --delete                delete available sms
  -h, --help                  help for dongel
  -o, --overcharge-alert      Alert for over charging
      --password string       Password to login with. Default will be taken from the config / env
      --username string       Username to login with. Default will be taken from the config / env
  -v, --view-sms              View SMS
```

Use `dongel [command] --help` for more information about a command.
