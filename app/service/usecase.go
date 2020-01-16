package service

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	ClearDB() *models.Error
	GetDBStatus() (models.Status, *models.Error)
}
