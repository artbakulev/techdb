package user

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	GetUserByNickname(nickname string) (models.User, *models.Error)
	GetUserByEmail(email string) (models.User, *models.Error)
	GetUserByEmailOrByNickname(data string, isEmail bool) (models.User, *models.Error)
	CreateUser(userNew models.UserNew) (models.User, *models.Error)
	UpdateUser(userUpdate models.UserUpdate) (models.User, *models.Error)
}
