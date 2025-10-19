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
go get github.com/swaggo/swag@latest 
go get github.com/swaggo/files
go get github.com/swaggo/gin-swagger

# jwt
go get github.com/golang-jwt/jwt/v5

# langchaingo
go get github.com/tmc/langchaingo/prompts
go get github.com/tmc/langchaingo/llms/openai
go get github.com/joho/godotenv

# other
go get github.com/pkg/errors
go get github.com/google/uuid
```

## Run

```bash
# install dependence
go mode tidy

# create db py docker
docker compose up -d

# run debug mode
# Execute vscode debug run "Attach to Go Process"
make debug
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