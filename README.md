# gecco

[![CircleCI](https://circleci.com/gh/ytakahashi/gecco.svg?style=shield&circle-token=262d744aaef58be4ebb8ba75d6d4f2d8c0b2c14a)](https://circleci.com/gh/ytakahashi/gecco)

Gecce is a command line tool to operate AWS EC2.

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
```

### list

`gecco list` lists EC2 instances.  
Listed information contains instanceId, instance size, instance status, and tags.

#### options (list)

- `--status`: only lists instances of specified status.
- `--tagKey`, `--tagValue`: only lists instances which have specified tag.

### connect

`gecco connect` connects to EC2 instances using `ssm start-session` command.  

**Note:** Following requirements should be satisfied when using this command.

- AWS cli version 1.16.12 or above should be installed.
- EC2 instances to be connected must be running the SSM Agent version 2.3.12 or above.

#### options (connect)

- `--target`: connect to the specified instance.
- `--interactive`/`-i`: if this option is specified, select an EC2 instance to connect interactively. When using this option, config file (`gecco.yml` or `gecco.toml`) is required at `~/.config/` directory and available interactive filtering command ([fzf](https://github.com/junegunn/fzf) or [peco](https://github.com/peco/peco)) should be provided.

Example of config file:

- `~/.config/gecco.yml`

  ```yaml
  InteractiveFilterCommand: fzf
  ```

- `~/.config/gecco.toml`

  ```toml
  InteractiveFilterCommand = "peco"
  ```
