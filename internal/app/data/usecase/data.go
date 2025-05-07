package usecase

import (
	"github.com/estella-studio/leon-backend/internal/app/data/repository"
	"github.com/estella-studio/leon-backend/internal/domain/dto"
	"github.com/estella-studio/leon-backend/internal/domain/entity"
	"github.com/estella-studio/leon-backend/internal/infra/jwt"
	"github.com/google/uuid"
)

type DataUseCaseItf interface {
	Add(add dto.Add) (dto.ResponseAdd, error)
}

type DataUseCase struct {
	dataRepo repository.DataMySQLItf
	jwt      jwt.JWTItf
}

func NewDataUseCase(dataRepo repository.DataMySQLItf, jwt *jwt.JWT) DataUseCaseItf {
	return &DataUseCase{
		dataRepo: dataRepo,
		jwt:      jwt,
	}
}

func (d *DataUseCase) Add(add dto.Add) (dto.ResponseAdd, error) {
	data := entity.Data{
		ID:     uuid.New(),
		UserID: add.UserID,
		Data:   add.Data,
	}

	err := d.dataRepo.Add(&data)
	if err != nil {
		return dto.ResponseAdd{}, err
	}

	return data.ParseToDTOResponseAdd(), nil
}
