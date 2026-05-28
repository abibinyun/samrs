package repository

import (
	userrepo "samrs-backend/internal/repository/user"

	"gorm.io/gorm"
)

type UserFilter = userrepo.UserFilter

type UserRepository = userrepo.UserRepository

func NewUserRepository(db *gorm.DB) UserRepository {
	return userrepo.NewUserRepository(db)
}
