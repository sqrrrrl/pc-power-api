package repo

import (
	"github.com/go-errors/errors"
	"github.com/pc-power-api/src/exceptions"
	"github.com/pc-power-api/src/infra/entity"
	"gorm.io/gorm"
)

var DeviceNotFoundError = exceptions.NewObjectNotFound("device not found")

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) Create(device *entity.Device) *errors.Error {
	err := r.db.Create(device).Error
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func (r *DeviceRepository) Update(device *entity.Device) *errors.Error {
	err := r.db.Save(device).Error
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func (r *DeviceRepository) Delete(device *entity.Device) *errors.Error {
	err := r.db.Delete(device).Error
	if err != nil {
		return errors.New(err)
	}
	return nil
}

func (r *DeviceRepository) GetById(id string) (*entity.Device, *errors.Error) {
	var device entity.Device
	err := r.db.First(&device, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New(DeviceNotFoundError)
		}
		return nil, errors.New(err)
	}
	return &device, nil
}
