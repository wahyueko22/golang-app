# echo-xorm

My own toy example with

- HTTP server: [labstack-echo](https://gihtub.com/labstack/echo)
- Database-driver: [go-sqlite3](https://github.com/mattn/go-sqlite3)
- ORM: [xorm](https://github.com/go-xorm/xorm)
- Authorization: [JSON Web Tokens](https://github.com/dgrijalva/jwt-go)


# Installation
## Prerequisites

golang 1.11+ (because of go modules)

## Installation process
```bash
GO111MODULE=on; go install
```

## Configuration
Config file is located at

`$GOPATH/src/github.com/corvinusz/echo-xorm/config/echo-xorm-config.toml`
Feel free to change.

By default application expects to find the configuration file at

`/usr/local/etc/echo-xorm-config.toml`

You can also point the path to the config file with the flag '-config'

## Database
Currently using:
- *sqlite3*-database, located as '/tmp/echo-xorm.sqlite.db' (change it in config)
- ORM [xorm](https://github.com/go-xorm/xorm)

# Application Run
```bash
echo-xorm -h # shows application flags
echo-xorm -config=$GOPATH/src/github.com/corvinusz/echo-xorm/config/echo-xorm-config.toml # runs app with default cfg
```

## Health check
```bash
curl http://localhost:11111/version
```

It should return something like

`{"result":"OK","version":"0.0.1","server_time":1501286982}`

Of course you feel free to use any of applicable applications such as:
- [insomnia](https://insomnia.rest/)
- [postman](https://www.getpostman.com/)
- etc..

or you browser plugin like:
- [chromerestclient](https://advancedrestclient.com/)
- [restclient](https://addons.mozilla.org/ru/firefox/addon/restclient/)
- etc...

# Testing

## Unit-tests for handlers (not yet completed)
```bash
cd app/server
GO111MODULE=on; go test -coverprofile=t.coverprofile ./...
go tool cover -html t.coverprofile
```

## BDD-style tests
Implemented in BDD-style with:
- Test Framework: [gomega](https://github.com/onsi/gomega)
- HTTP-Client: [Go-resty](https://github.com/go-resty/resty)

```bash
cd $GOPATH/src/github.com/corvinusz/echo-xorm/bddtests
GO111MODULE=on; go test
```

Test parameters defined in file:

`$GOPATH/src/github.com/corvinusz/echo-xorm/bddtests/test-config/echo-xorm-test-config.toml`

#License

MIT
