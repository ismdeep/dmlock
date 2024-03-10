package dmlock

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/ismdeep/dmlock/pkg/quantumid"
)

// DMLock distribute mysql lock
type DMLock struct {
	ID      string `gorm:"type:varchar(255);not null;primaryKey"` // lock id, format: please generate by `$ openssl rand --hex 16`
	Content string `gorm:"type:varchar(255);not null"`            // content, format: <customer-id>-<last-heart-beat-time>,...
	Quantum string `gorm:"type:varchar(255);not null"`            // quantum id
}

// LockMgmt lock manager
type LockMgmt struct {
	LockID          string
	CustomerID      string
	Bound           int
	MaxKeepDuration time.Duration
	ticker          *time.Ticker
	db              *gorm.DB
}

func NewLockMgmt(db *gorm.DB, lockID string, customerID string, bound int, maxKeepDuration time.Duration) *LockMgmt {
	if err := db.AutoMigrate(&DMLock{}); err != nil {
		fmt.Printf("[ERROR] failed to run AutoMigrate: %v\n", err.Error())
	}

	var cnt int64
	db.Model(&DMLock{}).Where("id = ?", lockID).Count(&cnt)
	if cnt <= 0 {
		db.Where("id = ?", lockID).Create(&DMLock{
			ID:      lockID,
			Content: "",
			Quantum: quantumid.Base58(),
		})
	}

	return &LockMgmt{
		LockID:          lockID,
		CustomerID:      customerID,
		Bound:           bound,
		MaxKeepDuration: maxKeepDuration,
		db:              db,
	}
}

func (receiver *LockMgmt) TryLock() bool {
	// get current lock info
	var cur DMLock
	if err := receiver.db.Where("id = ?", receiver.LockID).First(&cur).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("ERR:", err)
		return false
	}

	// parse content from receiver.Content
	customers := UnmarshalCustomerMap(cur.Content)
	for k, v := range customers {
		if v < time.Now().Add(-receiver.MaxKeepDuration).UnixNano() {
			delete(customers, k)
		}
	}

	if len(customers) >= receiver.Bound {
		return false
	}

	customers[receiver.CustomerID] = time.Now().UnixNano()

	// update
	if receiver.db.Model(&DMLock{}).Where("id = ? AND quantum = ?", receiver.LockID, cur.Quantum).
		Updates(map[string]interface{}{
			"content": MarshalCustomerMap(customers),
			"quantum": quantumid.Base58(),
		}).RowsAffected == 0 {
		return false
	}

	// start an heart beat ticker
	ticker := time.NewTicker(100 * time.Millisecond)
	receiver.ticker = ticker
	go func() {
		for range receiver.ticker.C {
			tick := func() {
				var cur DMLock
				if err := receiver.db.Where("id = ?", receiver.LockID).First(&cur).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return
				}
				customers := UnmarshalCustomerMap(cur.Content)
				if _, ok := customers[receiver.CustomerID]; ok {
					customers[receiver.CustomerID] = time.Now().UnixNano()
				}
				receiver.db.Model(&DMLock{}).Where("id = ? AND quantum = ?", receiver.LockID, cur.Quantum).Updates(map[string]interface{}{
					"content": MarshalCustomerMap(customers),
					"quantum": quantumid.Base58(),
				})
			}
			tick()
		}
	}()

	return true
}

func (receiver *LockMgmt) Lock() {
	for {
		if receiver.TryLock() {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}

func (receiver *LockMgmt) Unlock() {
	receiver.ticker.Stop()
	time.Sleep(130 * time.Millisecond)

	for {
		// get current lock info
		var cur DMLock
		if err := receiver.db.Where("id = ?", receiver.LockID).First(&cur).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		customers := UnmarshalCustomerMap(cur.Content)
		delete(customers, receiver.CustomerID)
		// update
		if receiver.db.Model(&DMLock{}).Where("id = ? AND quantum = ?", receiver.LockID, cur.Quantum).
			Updates(map[string]interface{}{
				"content": MarshalCustomerMap(customers),
				"quantum": quantumid.Base58(),
			}).RowsAffected == 0 {
			continue
		}
		break
	}
}
