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

func (u userUsecase) GetUserByNickname(nickname string) (models.User, *models.Error) {
	return u.userRepo.GetByNickname(nickname)
}

func (u userUsecase) GetUserByEmail(email string) (models.User, *models.Error) {
	return u.userRepo.GetByEmail(email)
}

func (u userUsecase) CreateUser(userNew models.User) (models.User, *models.Error) {
	return u.userRepo.Create(userNew)
}

func (u userUsecase) UpdateUser(userUpdate models.User) (models.User, *models.Error) {
	return u.userRepo.Update(userUpdate)
}
