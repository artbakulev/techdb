package validation

import (
	"github.com/artbakulev/techdb/app/models"
	"regexp"
)

var nicknamePattern = regexp.MustCompile("^[A-Za-z0-9_]+$")

func ValidateNickname(nickname string) *models.Error {
	if !nicknamePattern.MatchString(nickname) {
		return models.NewError(400, models.BadRequestError, "invalid nickname")
	}
	return nil
}
