# gecco

[![CircleCI](https://circleci.com/gh/ytakahashi/gecco.svg?style=shield&circle-token=262d744aaef58be4ebb8ba75d6d4f2d8c0b2c14a)](https://circleci.com/gh/ytakahashi/gecco)
[![codecov](https://codecov.io/gh/ytakahashi/gecco/branch/master/graph/badge.svg)](https://codecov.io/gh/ytakahashi/gecco)
[![Go Report Card](https://goreportcard.com/badge/github.com/ytakahashi/gecco)](https://goreportcard.com/report/github.com/ytakahashi/gecco)

Gecce is a command line tool to operate AWS EC2.

![gecco](https://user-images.githubusercontent.com/26239560/45939940-1c811a00-c011-11e8-9f85-90c7d76d6733.gif)

## Requirement

Aws credential should be provided with shared config file (`~/.aws/config`) or shared credentials file (`~/.aws/credentials`).

## Usage

```shell
$ gecco  -h
A Command Line Tool To Oprtate AWS EC2.

Usage:
  gecco [command]

Available Commands:
  connect     connect to EC2 instance
  help        Help about any command
  list        lists EC2 instances
  start       start specified EC2 instance
  stop        stop specified EC2 instance
```

### list

`gecco list` lists EC2 instances.  
Listed information contains instanceId, instance size, instance status, and tags.

#### options (list)

- `--status`: only lists instances of specified status.
- `--tagKey`, `--tagValue`: only lists instances which have specified tag.

### start/stop

`gecco start` command starts/stops EC2 instance.  
Target EC2 instance is provided by option `--target` or `-i/--interactive`.  

#### options (start/stop)

- `--target`: connect to the specified instance.
- `--interactive`/`-i`: select an EC2 instance to connect interactively. When using this option, [config file](#config-file) is required.

### connect

`gecco connect` connects to EC2 instances using `ssm start-session` command.  
Target EC2 instance is provided by option `--target` or `-i/--interactive`.  

**Note:** Following requirements should be satisfied when using this command.

- AWS cli version 1.16.12 or above is installed.
- EC2 instance to be connected is running the SSM Agent version 2.3.12 or above.

#### options (connect)

- `--target`: connect to the specified instance.
- `--interactive`/`-i`: select an EC2 instance to connect interactively. When using this option, [config file](#config-file) is required.

### config file

Config file (`gecco.yml` or `gecco.toml`) should be placed at `~/.config/` directory.
Please specify available interactive filtering command (e.g., [fzf](https://github.com/junegunn/fzf) and [peco](https://github.com/peco/peco)) as `InteractiveFilterCommand`.

Example of config file:

- `~/.config/gecco.yml`

  ```yaml
  InteractiveFilterCommand: fzf
  ```

- `~/.config/gecco.toml`

  ```toml
  InteractiveFilterCommand = "peco"
  ```
