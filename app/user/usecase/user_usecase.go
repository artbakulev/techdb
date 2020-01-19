package usecase

import (
	"github.com/artbakulev/techdb/app/models"
	"github.com/artbakulev/techdb/app/user"
)

type userUsecase struct {
	userRepo user.Repository
}

func NewUserUsecase(userRepo user.Repository) user.Usecase {
	return &userUsecase{userRepo: userRepo}
}

func (u userUsecase) GetUserByEmailOrByNickname(data string, isEmail bool) (models.User, *models.Error) {
	if isEmail {
		return u.userRepo.GetByEmail(data)
	}
	return u.userRepo.GetByNickname(data)
}

func (u userUsecase) CreateUser(userNew models.User) (models.User, *models.Error) {
	return u.userRepo.Create(userNew)
}

func (u userUsecase) UpdateUser(userUpdate models.User) (models.User, *models.Error) {
	return u.userRepo.Update(userUpdate)
}
