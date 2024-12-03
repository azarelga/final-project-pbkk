package services

import (
    "snipetty.com/main/repositories"
)

type UserService struct {
    repo *repositories.UserRepository
}