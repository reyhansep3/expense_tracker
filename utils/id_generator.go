package utils

import (
	"strconv"
	"sync"
	"time"
)

var (
	mu      sync.Mutex
	counter int64
)

func GenerateId() int64 {
	mu.Lock()
	defer mu.Unlock()
	currentTime := time.Now()
	counter++
	formatedTime := currentTime.Format("20060102150405")

	timePart, err := strconv.ParseInt(formatedTime, 10, 64)
	if err != nil {
		timePart = time.Now().Unix()
	}

	return timePart*1000 + counter
}
