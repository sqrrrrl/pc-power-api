package repo

import (
	"github.com/go-errors/errors"
	"github.com/pc-power-api/src/exceptions"
	"github.com/pc-power-api/src/infra/entity"
	"gorm.io/gorm"
)

var UsernameAlreadyExistsError = exceptions.NewObjectAlreadyExist("This username is already taken")
var UserNotFoundError = exceptions.NewObjectNotFound("The user was not found")

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(user *entity.User) *errors.Error {
	err := r.db.Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errors.New(UsernameAlreadyExistsError)
		}
		return errors.New(err)
	}
	return nil
}

func (r *UserRepository) GetById(id string) (*entity.User, *errors.Error) {
	var user entity.User
	err := r.db.Preload("Devices").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(UserNotFoundError)
		}
		return nil, errors.New(err)
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*entity.User, *errors.Error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(UserNotFoundError)
		}
		return nil, errors.New(err)
	}
	return &user, nil
}
