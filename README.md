# gtmx

Go-Tmux, easy tmux runner.

## Install

```
$ go get -u github.com/meinside/gtmx
```

## Usage

### 1. start a new session

```bash
# will start a new session named as your hostname
$ gtmx
```

### 2. resume or switch to a session

```bash
# will resume, or switch to a session with given name
$ gtmx [SESSION_NAME]

# if session name is not given, it will be your hostname
$ gtmx
```

### 3. start a predefined session

#### create a new config file

You can predefine sessions in your config file. (at `$XDG_CONFIG_HOME/gtmx/config.json`)

If you need a sample config file,

```bash
$ gtmx -g

# or
$ gtmx --gen-config
```

it will be printed on the screen.

#### start a session defined in the config file

```bash
$ gtmx [SESSION_NAME]
```

### 4. list predefined/running sessions

```bash
$ gtmx -l

# or
$ gtmx --list
```

### 5. terminate this session

```bash
$ gtmx -q

# or
$ gtmx --quit
```

## License

MIT
