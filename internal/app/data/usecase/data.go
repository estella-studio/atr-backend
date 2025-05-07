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
	Retrieve(retrieve dto.Retrieve) (dto.ResponseRetrieve, error)
	List(userID uuid.UUID) (*[]dto.ResponseList, error)
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

func (d *DataUseCase) Retrieve(retrieve dto.Retrieve) (dto.ResponseRetrieve, error) {
	data := entity.Data{
		ID:     retrieve.ID,
		UserID: retrieve.UserID,
	}

	err := d.dataRepo.Retrieve(&data)
	if err != nil {
		return dto.ResponseRetrieve{}, err
	}

	return data.ParseToDTOResponseRetrieve(), nil
}

func (d *DataUseCase) List(userID uuid.UUID) (*[]dto.ResponseList, error) {
	data := new([]entity.Data)

	err := d.dataRepo.List(data)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ResponseList, len(*data))

	for i, data := range *data {
		res[i] = data.ParseToDTOResponseList()
	}

	return &res, nil
}
