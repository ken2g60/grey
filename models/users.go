package models

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        int     `json:"id" gorm:"type:integer;primaryKey"`
	UserId    string  `json:"user_id" gorm:"type:uuid;not null;index"`
	Email     string  `json:"email" gorm:"uniqueIndex;not null"`
	Password  string  `json:"-" gorm:"not null"`
	Account   Account `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	CreatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UserId = uuid.NewString()
	return nil
}

func CreateUser(ctx context.Context, db *gorm.DB, User *User) (err error) {
	err = db.WithContext(ctx).Create(&User).Error
	if err != nil {
		return err
	}
	return nil
}

func IsEmailExists(ctx context.Context, db *gorm.DB, email string) (*User, error) {
	var user User
	err := db.Debug().WithContext(ctx).Model(&user).Where("email = ?", email).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func PasswordCompare(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func UserProfile(ctx context.Context, db *gorm.DB, user_id string) (*User, error) {
	var user User
	err := db.Model(&User{}).Preload("Account").Where("user_id = ?", user_id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
