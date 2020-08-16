# YuubariGo

Yet anther KanColle Proxy tool.

No cache, no sniffing, no acceleration, just retrying upon unexcepted EOF network error. 

**Use at your own risk.**

# Build

## Command line version

```bash
go build yuubari_go/cli
```

## SysTray version

### Windows

```bash
go build -o YuubariGo.exe -ldflags "-H=windowsgui" yuubari_go/systray
```

### Linux

build dependencies: 
- gtk3
- libappindicator3

have those dependencies ready, then

```bash
go build -o YuubariGo yuubari_go/systray
```

### Mac OS

![swole](https://files.catbox.moe/7he6r1.png)

not tested due to having no mac device, contributions are welcome

# Usage

```
  -debug
        enable debug log
  -interval int
        retry interval (seconds) (default 5)
  -port int
        listen port (default 8099)
  -proxy string
        backend proxy url
  -retry int
        max retry times (default 3)
```

# Credits

- github.com/getlantern/systray
- github.com/elazarl/goproxy