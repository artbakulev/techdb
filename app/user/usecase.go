package user

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	GetUserByEmailOrByNickname(data string, isEmail bool) (models.User, *models.Error)
	CreateUser(userNew models.User) (models.User, *models.Error)
	UpdateUser(userUpdate models.User) (models.User, *models.Error)
}
