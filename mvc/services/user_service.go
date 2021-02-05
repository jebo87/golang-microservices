package services

import (
	"github.com/jebo87/golang-microservices/mvc/domain"
	"github.com/jebo87/golang-microservices/mvc/utils"
)

func GetUser(userId int64) (*domain.User, *utils.ApplicationError) {
	return domain.GetUser(userId)
}