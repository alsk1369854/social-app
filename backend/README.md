

## Setup
- go: 1.24.5

### Install cli
```bash
go install github.com/air-verse/air@latest
go install github.com/go-delve/delve/cmd/dlv@latest
```

#### Add cli
```bash
# add to .bashrc or .zshrc
export PATH=$PATH:$(go env GOPATH)/bin
```

## Create `.air.toml`
```bash
air init

# update `.air.toml`
[build]
    cmd = "make debug-main-file"
    full_bin = "dlv exec ./tmp/main --listen=127.0.0.1:12345 --headless=true --api-version=2 --accept-multiclient --continue --log -- "
```

## Create `.vscode/launch.json`
```json
{
    "version": "0.2.0",
    "configurations": [

        {
            "name": "Attach to Go Process",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "host": "127.0.0.1",
            "port": 12345,
        }
    ]
}
```

### Install libs
```bash
go get -u gorm.io/gorm
```


### Debug run
```bash
make debug

# Execute vscode debug run "Attach to Go Process"
```
