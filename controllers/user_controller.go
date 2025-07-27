package controllers

import (
	"exp_tracker/models"
	"exp_tracker/repository"
	"exp_tracker/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Repo repository.UserRepository
}

func (c *UserController) CreateUsers(ctx *gin.Context) {
	var newUser models.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	newUser.ID = utils.GenerateId()

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	newUser.Password = hashedPassword
	generatedToken, err := utils.CreateToken(newUser.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	newUser.Token = generatedToken
	exist, err := c.Repo.IsUserExist(newUser.Name, newUser.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Database error",
		})
		return
	}
	if exist {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Username or email already used",
		})
		return
	}
	if err := c.Repo.Create(&newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"users ": newUser,
		})
	}

}
