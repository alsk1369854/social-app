

## Setup
- go: 1.24.5

### Install cli
```bash
go install github.com/air-verse/air@latest
go install github.com/go-delve/delve/cmd/dlv@latest
# go get github.com/swaggo/swag@latest 
go install github.com/swaggo/swag/cmd/swag@latest 
```

#### Add cli
```bash
# add to .bashrc or .zshrc
export PATH=$PATH:$(go env GOPATH)/bin
```

### Create `.air.toml`
```bash
air init

# update `.air.toml`
[build]
  cmd = "make debug-file"
  full_bin = "dlv exec ./tmp/main --listen=127.0.0.1:12345 --headless=true --api-version=2 --accept-multiclient --continue --log -- "
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs", "tmp", "postgres-data"]
```

### Create `.vscode/launch.json`
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
# env loader
go get github.com/joho/godotenv

# gin
go get -u github.com/gin-gonic/gin

# gorm
go get -u gorm.io/gorm
go get gorm.io/driver/postgres
go get gorm.io/driver/sqlite

# swagger
go get github.com/swaggo/files
go get github.com/swaggo/gin-swagger

# other
go get github.com/pkg/errors
go get github.com/google/uuid
```

### Import Swagger lib
```go
import (
    swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "<your-mod-name>/docs"
)

func main() {
    ....
    // http://localhost:8080/swagger/index.html
	engin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    ....
}
```


### Debug run
```bash
make debug

# Execute vscode debug run "Attach to Go Process"
```


## Test
```bash
go test -cover ./... -v
```

### coverage
```bash
go test --coverprofile=./tmp/coverage.out ./...

go tool cover -func=./tmp/coverage.out

go tool cover -html=./tmp/coverage.out -o ./tmp/coverage.html
```