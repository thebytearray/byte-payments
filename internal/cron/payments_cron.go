package cron

import (
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/repository"
	"github.com/thebytearray/BytePayments/service"
)

var (
	mu        sync.Mutex
	isRunning bool
)

func safeProcessPendingPayments() {
	mu.Lock()
	if isRunning {
		log.Println("Previous job still running, skipping this run")
		mu.Unlock()
		return
	}
	isRunning = true
	mu.Unlock()

	start := time.Now()
	log.Println("Started processing payments at:", start.Format(time.RFC3339))

	defer func() {
		duration := time.Since(start)
		log.Println("Finished processing in:", duration)

		mu.Lock()
		isRunning = false
		mu.Unlock()
	}()

	paymentService := service.NewPaymentService(repository.NewPaymentRepository(database.DB))
	paymentService.ProcessPendingPayments()
}

func NewPaymentCron() {
	c := cron.New()
	c.AddFunc("@every 30s", func() {
		safeProcessPendingPayments()
	})
	c.Start()
}
