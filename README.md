# gtmx

Go-Tmux, easy tmux runner.

## Install

```
$ go install github.com/meinside/gtmx@latest
```

## Usage

### 1. Start a new session

```bash
# will start a new session named as your hostname
$ gtmx
```

### 2. Resume or switch to a session

```bash
# will resume, or switch to a session with given name
$ gtmx [SESSION_NAME]

# if session name is not given, it will be your hostname
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

it will be printed on the screen.

#### start a session defined in the config file

```bash
$ gtmx [SESSION_NAME]
```

### 4. List predefined/running sessions

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

## License

MIT
