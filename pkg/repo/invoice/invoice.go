package invoice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Invoice struct {
	Id uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`

	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"softDelete:milli"`
	UpdatedAt *time.Time
	DueAt     time.Time

	FromUserId uuid.UUID
	ToUserId   uuid.UUID

	Status InvoiceStatus

	Description string
	Items       []*InvoiceItem
	// TODO: Added tags to here
}

type InvoiceItem struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	InvoiceID   uuid.UUID `gorm:"type:uuid;index"`
	Description string
	Amount      int64
	// TODO: Added tags to here
}

type InvoiceStatus uint8

const (
	InvoiceStatus_PENDING InvoiceStatus = iota
	InvoiceStatus_PAID
	InvoiceStatus_REFUSED
	InvoiceStatus_CANCELED
)

var InvoiceStatusNames = map[InvoiceStatus]string{
	InvoiceStatus_PENDING:  "pending",
	InvoiceStatus_PAID:     "paid",
	InvoiceStatus_REFUSED:  "refused",
	InvoiceStatus_CANCELED: "canceled",
}

type GormInvoiceRepo struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewInvoiceGormRepo(db *gorm.DB, log *zerolog.Logger) (*GormInvoiceRepo, error) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err := db.AutoMigrate(
		&Invoice{},
	)
	if err != nil {
		return nil, err
	}
	return &GormInvoiceRepo{
		db:  db,
		log: log.With().Str("component", "pkg.repo.invoice").Logger(),
	}, nil
}

func (gir GormInvoiceRepo) Create(ctx context.Context, u *Invoice) (*Invoice, error) {
	panic("not implemented") // TODO: Implement
}

func (gir GormInvoiceRepo) Update(ctx context.Context, id uuid.UUID, u *Invoice) (*Invoice, error) {
	panic("not implemented") // TODO: Implement
}

func (gir GormInvoiceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	panic("not implemented") // TODO: Implement
}

func (gir GormInvoiceRepo) Get(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	panic("not implemented") // TODO: Implement
}

func (gir GormInvoiceRepo) GetByUser(ctx context.Context, uuid uuid.UUID) (*Invoice, error) {
	panic("not implemented") // TODO: Implement
}
