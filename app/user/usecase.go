package user

import "github.com/artbakulev/techdb/app/models"

type Usecase interface {
	GetUserByEmail(email string) (models.User, *models.Error)
	GetUserByNickname(nickname string) (models.User, *models.Error)
	CreateUser(userNew models.User) (models.User, *models.Error)
	UpdateUser(userUpdate models.User) (models.User, *models.Error)
}
