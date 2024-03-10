# DMLock

DMLock is a Distributed MySQL-based Locking SDK.

Using MySQL as the backend for distributed locking adds reliability and scalability, as MySQL is a robust and widely-used relational database management system capable of handling large amounts of data and concurrent connections.

# Quick Start

## Prerequisites

- [Go](https://go.dev): any one of the three lastest major [release](https://go.dev/doc/devel/release) (we test it with these).

## Getting DMLock

With Go module support, simple add the following import

```
import "github.com/ismdeep/dmlock"
```

to your code, and then `go [build|run|test]` will automatically fetch the nessessary dependencies.

Otherwise, run the following Go command to install the `dmlock` packages:

```
$ go get -u github.com/ismdeep/dmlock
```

## Running DMLock

First you need to import DMLock package for using DMLock, one simplest example like the follow `example.go`:

```
package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ismdeep/dmlock"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:db123456@tcp(127.0.0.1:6392)/db_dmlock"))
	if err != nil {
		panic(err)
	}

	lckID := "lck_my_id"
	customerID := "cus_this_customer_id"
	lck := dmlock.NewLockMgmt(db, lckID, customerID, 1, 10*time.Second)
	lck.Lock()
	fmt.Println(lck.CustomerID)
	time.Sleep(5 * time.Second)
	lck.Unlock()
	fmt.Println(lck.CustomerID, "is gone.")
}
```

And use the Go command to run the demo:

```
$ go run example.go
```

