package r3meters

import (
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type FlowRateMeter struct {
	counter uint64
	last    uint64
}

func (m *FlowRateMeter) Start(filePath string) error {
	logFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	logger := log.New(logFile, "", log.LstdFlags)
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			delta := atomic.LoadUint64(&m.counter) - atomic.LoadUint64(&m.last)
			logger.Println(delta)
			atomic.StoreUint64(&m.last, atomic.LoadUint64(&m.counter))
		}
	}()
	return nil
}

func (m *FlowRateMeter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		atomic.AddUint64(&m.counter, 1)
	})
}
