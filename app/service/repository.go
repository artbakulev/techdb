package service

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	Clear() *models.Error
	GetStatus() (models.Status, *models.Error)
}
