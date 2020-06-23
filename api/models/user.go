package models

import (
	"errors"
	"html"
	"os"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.com/norfabagas/auth/api/utils/crypto"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	PublicID  string    `gorm:"size:255;not null;unique" json:"public_id"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func EscapeAndTrimString(input string) string {
	return html.EscapeString(strings.TrimSpace(input))
}

func (user *User) BeforeSave() error {
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (user *User) Prepare() {
	user.ID = 0
	user.Name = EscapeAndTrimString(user.Name)
	user.Email = EscapeAndTrimString(user.Email)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
}

func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email Format")
		}
		return nil
	case "register":
		if user.Name == "" {
			return errors.New("Required Name")
		}
		if user.Email == "" {
			return errors.New("Required Email")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid Email Format")
		}
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if passwordLength := len([]rune(user.Password)); passwordLength < 6 {
			return errors.New("Password minimum is 6 characters")
		}
		return nil
	case "update":
		if user.Name == "" {
			return errors.New("Required Name")
		}
		return nil
	default:
		return errors.New("Undefined Action")
	}
}

func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error

	user.Name, err = crypto.Encrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		return &User{}, err
	}

	user.PublicID = crypto.MD5Hash(user.Email + user.CreatedAt.String())

	err = db.Debug().Create(&user).Error
	if err != nil {
		return &User{}, err
	}

	user.Name, _ = crypto.Decrypt(user.Name, os.Getenv("APP_KEY"))
	return user, nil
}

func (user *User) FindUserByID(db *gorm.DB, id uint32) (*User, error) {
	err := db.Debug().Model(&User{}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return user, nil
}

func (user *User) UpdateUser(db *gorm.DB, id uint32) (*User, error) {
	encryptedName, err := crypto.Encrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		return &User{}, err
	}

	db = db.Debug().Model(&User{}).Where("id = ?", id).UpdateColumns(
		map[string]interface{}{
			"name":       encryptedName,
			"updated_at": user.UpdatedAt,
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return &User{}, err
	}

	user.Name, err = crypto.Decrypt(user.Name, os.Getenv("APP_KEY"))
	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) DeleteUser(db *gorm.DB, id uint32) (string, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", id).Take(&user).Delete(&user)
	if db.Error != nil {
		return "", db.Error
	}
	return user.Email, nil
}
