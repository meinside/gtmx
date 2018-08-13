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

### 2. resume a session

```bash
# will resume a session with given name
$ gtmx [SESSION_NAME]

# if session name is not given, it will be your hostname
$ gtmx
```

### 3. start a predefined session

#### create a new config file

You can predefine sessions in your config file. (at `~/.gtmx.conf`)

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

## License

Copyright (c) 2018 Sungjin Han

MIT License

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

