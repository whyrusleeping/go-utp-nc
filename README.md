# go-utp-nc

A basic implementation of netcat using utp instead of tcp

## Installation
```
go get github.com/whyrusleeping/go-utp-nc
```

## Usage

Listen on a given addr/port

```bash
$ go-utp-nc -l 127.0.0.1 5555
```

Dial to a given addr/port
```bash
$ go-utp-nc 127.0.0.1 5555
```

### License
MIT
