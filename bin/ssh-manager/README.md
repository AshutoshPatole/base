## Simple SSH Manager

Simple SSH Manager is a CLI utility designed for efficiently managing a large number of servers. It creates and organizes servers into groups for easy handling and maintenance. The configuration is stored inside the `.ssm.json` file under the `$HOME` directory.

## To Add a Server:

```bash
$ ssm add server serverName -u userName -g groupName -e dev -a test-machine
```

**Options:**

* `-u`: User name for connecting and configuring the machine for the first time.
* `-g`: Group name for the server. If the group does not exist and is mentioned for the first time, `ssm` will automatically create the group and add the server under the group.
* `-e`: Environment for the server. The allowed environments are:
    * `prd`: Production
    * `ppd`: Pre-production
    * `sit`: System Integration Testing
    * `uat`: User Acceptance Testing
    * `dev`: Development machine
* `-a`: Alias for the machine

## To List Servers and Groups

**To list groups:**

```bash
$ ssm list group
```

**To list servers and group**
```bash
$ ssm list server -g groupName
```

**To list servers related to a specific environment**
```bash
$ ssm list server -g groupName -e uat
```

## To connect to a server:

```bash
$ ssm connect chennai
? Select server  [Use arrows to move, type to filter]
> skillset (ppd)
  doc ash (ppd)
  docash01 (prd)
  pim (dev)

```
You can use arrows or type a subset of the machine name and it will filter the results.

Or you can use `-e` to filter based on environmnets

## To reverse copy files
for copying files from remote machine to local machine use `reverse-copy` or `rcp` in short.

```bash
$ ssm rcp -f "/path/to/remote-file.txt" 
? Select server  [Use arrows to move, type to filter]
> skillset (ppd)
  doc ash (ppd)
  docash01 (prd)
  pim (dev)

The file is dowloaded at C:\User\shadow\Dowloads\remote-file.txt
```

## Bulk importing
You can import the configurations using YAML format. To check the template use 

```bash
$ ssm import template

- name: groupname
  user: username
  env:
    - name: dev|uat|sit|ppd|prd
      servers:
        - hostname: example1.com
          alias: dev engine
        - hostname: example2.com
          alias: dev service
        - hostname: example3.com
          alias: prod engine
- name: groupname2
  user: anotheruser
  env:
    - name: dev|uat|sit|ppd|prd
      servers:
        - hostname: example4.com
          alias: dev engine
        - hostname: example5.com
          alias: dev service
        - hostname: example6.com
          alias: prod engine

File saved at C:\Users\shadow\ssh-import-template.yaml
```

To import your file use
```bash
$ ssm import -f C:\path\import.yaml

# or for importing a specific group use 
$ ssm import -f C:\path\import.yaml -g groupName
```

# for any help

```bash
$ ssm help

A CLI utility for managing SSH servers and keys efficiently:

This utility is designed to streamline SSH key handling by categorizing servers into groups.
To get started, users can run the 'ssm import template' command to obtain a predefined template for server configuration

Usage:
  ssm [command]

Available Commands:
  add          Add servers and groups in configurations
  completion   Generate the autocompletion script for the specified shell
  connect      Connect to the servers
  help         Help about any command
  import       Import file
  list         List groups and servers
  reverse-copy Download file from remote machine

Flags:
      --config string   config file (default is $HOME/.ssh-manager.json)
  -h, --help            help for ssm

Use "ssm [command] --help" for more information about a command.

```


todo: 
- [ ] Add backup and migrate subcommands
- [ ] Allow custom download locations in rcp
- [ ] Update subcommand for servers 
- [ ] Implement copy command (scp like)