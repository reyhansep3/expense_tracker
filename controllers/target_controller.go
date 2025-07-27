package controllers

import (
	"exp_tracker/middleware"
	"exp_tracker/models"
	"exp_tracker/repository"
	"exp_tracker/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TargetController struct {
	Repo repository.TargetRepository
}

func (c *TargetController) AddTarget(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var target models.Target

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	dst := fmt.Sprintf("./uploads/%s", fileHeader.Filename)
	if err := ctx.SaveUploadedFile(fileHeader, dst); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}
	target.File = fileHeader.Filename

	target.Title = ctx.PostForm("title")
	target.PaymentMethod = ctx.PostForm("payment_method")
	target.Description = ctx.PostForm("description")
	amount, err := strconv.ParseInt(ctx.PostForm("amount"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	prevTotal, err := c.Repo.GetTotalAmountByUser(int(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total amount"})
		return
	}

	target.TotalAmount = prevTotal + amount

	startDateStr := ctx.PostForm("start_date")
	endDateStr := ctx.PostForm("end_date")

	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)
	target.StartDate = startDate
	target.EndDate = endDate
	target.Amount = int64(amount)

	now := time.Now()

	target.ID = utils.GenerateId()
	target.UserId = userID
	target.ModifiedAt = now
	target.ModifiedBy = userID

	if err := c.Repo.Create(&target); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"target": target,
	})
}

func (c *TargetController) GetAllTarget(ctx *gin.Context) {
	userId := ctx.Param("userId")
	fmt.Println("DEBUG userId:", userId)

	target, err := c.Repo.GetTargetExpense(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"target": target,
	})
}

func (c *TargetController) UpdateTargetExpense(ctx *gin.Context) {
	userId, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to convert id in to int",
		})
	}

	var target models.Target

	fileHeader, err := ctx.FormFile("file")
	if err == nil {
		dst := fmt.Sprintf("./uploads/%s", fileHeader.Filename)
		if err := ctx.SaveUploadedFile(fileHeader, dst); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save file",
			})
			return
		}
		target.File = fileHeader.Filename
	}
	target.ID = int64(id)
	target.UserId = userId
	target.Title = ctx.PostForm("title")
	target.PaymentMethod = ctx.PostForm("payment_method")
	target.Description = ctx.PostForm("description")

	amountStr := ctx.PostForm("amount")
	amount, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}
	target.Amount = amount

	startDateStr := ctx.PostForm("start_date")
	endDateStr := ctx.PostForm("end_date")
	startDate, _ := time.Parse("2006-01-02", startDateStr)
	endDate, _ := time.Parse("2006-01-02", endDateStr)
	target.StartDate = startDate
	target.EndDate = endDate
	target.CreateBy = userId

	target.ModifiedAt = time.Now()
	target.ModifiedBy = userId

	prevTotal, err := c.Repo.GetTotalAmountByUser(int(userId))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total amount"})
		return
	}

	target.TotalAmount = prevTotal + amount

	updatedTarget, err := c.Repo.UpdateTarget(id, target)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"updated": updatedTarget,
	})
}

func (c *TargetController) DeleteTarget(ctx *gin.Context) {
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

	if err := c.Repo.DeleteTargetByID(int64(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete expenses by id",
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"success": "data has been deleted",
		})
	}

}
