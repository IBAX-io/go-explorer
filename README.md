# IBAX Blockchain explorer
[![Go Reference](https://pkg.go.dev/badge/github.com/IBAX-io/go-explorer.svg)](https://pkg.go.dev/github.com/IBAX-io/go-explorer)
[![Go Report Card](https://goreportcard.com/badge/github.com/IBAX-io/go-explorer)](https://goreportcard.com/report/github.com/IBAX-io/go-explorer)

Golang client for the [scan.ibax.network](https://scan.ibax.network),full implementation (account, transaction,ecossytem, contract, block, honor node), 
depends on [go-ibax](https://github.com/IBAX-io/go-ibax)

## Usage
``` go
go get github.com/IBAX-io/go-explorer
```

### Install

```shell
make all
```

### Set Config
You can modify your configuration in conf/config.yml and connect it to the full node database, do not connect it to the master node network, which may cause unknown security issues.


### Start 🚀

```shell
make startup
```


## License
Use of this work is governed by an MIT License.

You may find a license copy in project root.