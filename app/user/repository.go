package user

import "github.com/artbakulev/techdb/app/models"

type Repository interface {
	GetByNickname(nickname string) (models.User, *models.Error)
	GetByEmail(email string) (models.User, *models.Error)
	Create(userNew models.User) (models.User, *models.Error)
	Update(userUpdate models.User) (models.User, *models.Error)
}
