# YuubariGo

Yet anther KanColle Proxy tool.

No cache, no sniffing, just retrying upon network error. 

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

### Linux/Mac OS

not tested, contribution is welcome

# Usage

```
  -debug
        enable debug log
  -interval int
        retry interval (seconds) (default 5)
  -port int
        listen port (default 8099)
  -retry int
        max retry times (default 3)
```
