# gtmx

`gtmx` is an easy [tmux](https://github.com/tmux/tmux/wiki) runner, short for 'go-tmux'.

## Install

```bash
# install the latest binary
$ go install github.com/meinside/gtmx@latest
```

## Usage

### 0. Help

For some help messages,

```bash
$ gtmx -h

# or
$ gtmx --help
```

### 1. Start a new session

```bash
# will start a new session named as your `hostname`:
$ gtmx
```

### 2. Resume or switch to a session

```bash
# will resume, or switch to a session with the given name:
$ gtmx [SESSION_NAME]

# if session name is not given, it will be your `hostname`:
$ gtmx
```

### 3. Start a predefined session

#### create a new config file

You can predefine sessions in your config file. (at `$XDG_CONFIG_HOME/gtmx/config.json`)

If you need a sample config file,

```bash
$ gtmx -g

# or
$ gtmx --gen-config
```

will print the sample config file (in JSON format) to stdout.

#### start a session defined in the config file

```bash
$ gtmx [SESSION_NAME_IN_CONFIG]
```

### 4. List predefined and/or running sessions

```bash
$ gtmx -l

# or
$ gtmx --list
```

### 5. Terminate current session

```bash
$ gtmx -q

# or
$ gtmx --quit
```

### 999. Print the version

```bash
$ gtmx -V

# or
$ gtmx --version
```

will print the version string to stdout.

## License

MIT
