package repo

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus uint8

const (
	InvoiceStatusPending = 0
	InvoiceStatusPaid    = 1
	InvoiceStatusRefused = 2
)

type Invoice struct {
	InvoiceId   uuid.UUID
	SenderId    uuid.UUID
	RecipientId uuid.UUID
	Items       []*InvoiceItem
	Status      InvoiceStatus
	Description string
	DueDate     time.Time
	CreatedAt   time.Time
}

type InvoiceItem struct {
	Title       string
	Amount      float64
	ReferenceId string
}

type InvoiceManager struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewInvoiceManager(db *gorm.DB) InvoiceManager {
	return InvoiceManager{
		db: db,
	}
}
