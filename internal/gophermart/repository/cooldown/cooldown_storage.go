package cooldown

import (
	"sync"
	"time"
)

var (
	cooldownSync    = &sync.Mutex{}
	failTimeUnixSec int64
	cooldownPeriod  int64
)

// NeedAccrualCooldown updates cooldown time.
func NeedAccrualCooldown(requestedCooldownPeriod int64) {
	cooldownSync.Lock()
	defer cooldownSync.Unlock()
	cooldownPeriod = requestedCooldownPeriod
	failTimeUnixSec = time.Now().Unix()
}

// IsAccrualReady returns true if 'Accrual' is ready for requests
// or returns false otherwise.
func IsAccrualReady() bool {
	cooldownSync.Lock()
	defer cooldownSync.Unlock()
	if failTimeUnixSec == 0 {
		return true
	}
	if failTimeUnixSec+cooldownPeriod < time.Now().Unix() {
		failTimeUnixSec = 0

		return true
	}

	return false
}
