package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type User struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email        string    `gorm:"not null;unique"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"softDelete:milli"`
}

type GormUserRepo struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewUserManager(db *gorm.DB, log *zerolog.Logger) (*GormUserRepo, error) {
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	err := db.AutoMigrate(
		&User{},
	)
	if err != nil {
		return nil, err
	}
	return &GormUserRepo{
		db:  db,
		log: log.With().Str("component", "pkg.repo.user").Logger(),
	}, nil
}

func (im GormUserRepo) CreateUser(ctx context.Context, u *User) (*User, error) {
	im.log.Trace().Msg("Creating user")
	tx := im.db.WithContext(ctx).Create(u)
	return u, tx.Error
}

func (im GormUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	im.log.Trace().Str("user_id", id.String()).Msg("Deleting user")
	tx := im.db.WithContext(ctx).Delete(&User{}, id)
	return tx.Error
}

func (im GormUserRepo) UpdateUser(ctx context.Context, id uuid.UUID, u *User) (*User, error) {
	im.log.Trace().Str("user_id", id.String()).Msg("Updating user")
	tx := im.db.WithContext(ctx).Model(&User{Id: id}).Updates(u)
	return u, tx.Error
}

func (im GormUserRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	im.log.Trace().Str("user_email", email).Msg("Getting user")
	u := &User{}
	tx := im.db.WithContext(ctx).Where("email = ?", email).First(&u)
	return u, tx.Error
}

func (im GormUserRepo) GetUserById(ctx context.Context, id uuid.UUID) (*User, error) {
	im.log.Trace().Str("user_id", id.String()).Msg("Getting user")
	u := &User{}
	tx := im.db.WithContext(ctx).First(&u, id)
	return u, tx.Error
}
