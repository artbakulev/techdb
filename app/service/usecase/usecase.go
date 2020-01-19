package usecase

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/service"
)

type serviceUsecase struct {
	serviceRepo service.Repository
}

func NewServiceUsecase(serviceRepo service.Repository) service.Usecase {
	return &serviceUsecase{serviceRepo: serviceRepo}
}

func (s serviceUsecase) ClearDB() *models.Error {
	return s.ClearDB()
}

func (s serviceUsecase) GetDBStatus() (models.Status, *models.Error) {
	return s.GetDBStatus()
}
