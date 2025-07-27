package controllers

import (
	"exp_tracker/middleware"
	"exp_tracker/models"
	"exp_tracker/repository"
	"exp_tracker/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	Repo repository.CategoryRepository
}

func (c *CategoryController) CreateCategory(ctx *gin.Context) {
	var newCategory models.Categories

	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := ctx.ShouldBindJSON(&newCategory); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	now := time.Now().Format("2006-01-02")
	newCategory.ID = utils.GenerateId()
	newCategory.UserId = userID
	newCategory.CreateAt = now
	newCategory.CreateBy = userID
	newCategory.ModifiedAt = now
	newCategory.ModifiedBy = userID

	if err := c.Repo.Create(&newCategory); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"category": newCategory,
	})
}

func (c *CategoryController) GetAllCategory(ctx *gin.Context) {
	usrIdParams := ctx.Param("userId")
	usrId, err := strconv.Atoi(usrIdParams)
	if err != nil {
		panic(err)
	}
	categories, err := c.Repo.GetAllData(usrId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get all data",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

func (c *CategoryController) GetCategoryByID(ctx *gin.Context) {

	_, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}

	idParams := ctx.Param("id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		panic(err)
	}
	usrIdParams := ctx.Param("userId")
	usrId, err := strconv.Atoi(usrIdParams)
	if err != nil {
		panic(err)
	}
	categories, err := c.Repo.GetDataByID(id, usrId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get category data by id",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})

}

func (c *CategoryController) DeleteCategoryByID(ctx *gin.Context) {
	_, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	idParams := ctx.Param("id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		panic(err)
	}
	usrIdParams := ctx.Param("userId")
	usrId, err := strconv.Atoi(usrIdParams)
	if err != nil {
		panic(err)
	}
	if err := c.Repo.DeleteDataByID(id, usrId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete category by id",
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "data has been deleted",
		})
	}

}
