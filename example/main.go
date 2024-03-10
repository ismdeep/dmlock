package main

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ismdeep/dmlock"
	"github.com/ismdeep/dmlock/pkg/quantumid"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:db123456@tcp(127.0.0.1:6392)/db_dmlock"))
	if err != nil {
		panic(err)
	}

	lckID := "lck_" + quantumid.UUIDTidy()
	//lckID := "lck_my_id"
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
			}()

			customerID := "cus_" + quantumid.UUIDTidy()
			lck := dmlock.NewLockMgmt(db, lckID, customerID, 1, 10*time.Second)
			lck.Lock()
			fmt.Println(lck.CustomerID)
			time.Sleep(5 * time.Second)
			lck.Unlock()
			fmt.Println(lck.CustomerID, "is gone.")
		}(i)
	}

	wg.Wait()
}
