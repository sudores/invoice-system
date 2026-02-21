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

	Status uint8

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

func (im GormInvoiceRepo) Create(ctx context.Context, u *Invoice) (*Invoice, error) {
	im.log.Trace().Msg("Creating invoice")
	tx := im.db.WithContext(ctx).Create(u)
	return u, tx.Error
}

func (im GormInvoiceRepo) Update(ctx context.Context, id uuid.UUID, u *Invoice) (*Invoice, error) {
	im.log.Trace().Str("invoice_id", id.String()).Msg("Updating invoice")
	tx := im.db.WithContext(ctx).Model(&Invoice{Id: id}).Updates(u)
	return u, tx.Error
}

func (im GormInvoiceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	im.log.Trace().Str("invoice_id", id.String()).Msg("Deleting invoice")
	tx := im.db.WithContext(ctx).Delete(&Invoice{}, id)
	return tx.Error
}

func (im GormInvoiceRepo) Get(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	im.log.Trace().Str("invoice_id", id.String()).Msg("Getting invoice")
	u := &Invoice{}
	tx := im.db.WithContext(ctx).First(&u, id)
	return u, tx.Error

}

func (im GormInvoiceRepo) GetFromUser(ctx context.Context, uuid uuid.UUID) (*Invoice, error) {
	panic("Not implemented")
}
func (im GormInvoiceRepo) GetToUser(ctx context.Context, uuid uuid.UUID) (*Invoice, error) {
	panic("Not implemented")
}
