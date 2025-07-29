package cron

import (
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/repository"
	"github.com/thebytearray/BytePayments/service"
)

func NewPaymentCron() {
	c := cron.New()
	c.AddFunc("@every 1m", func() {
		log.Println("Running payment checker at : ", time.Now().Format(time.RFC3339))
		paymentService := service.NewPaymentService(repository.NewPaymentRepository(database.DB))
		paymentService.ProcessPendingPayments()
	})
	c.Start()
}
