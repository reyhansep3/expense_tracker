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

type ExpenseController struct {
	Repo repository.ExpenseRepository
}

func (c *ExpenseController) AddExpenses(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var expense models.Expenses
	if err := ctx.ShouldBindJSON(&expense); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	expense.ID = utils.GenerateId()
	expense.UserId = userID
	expense.CreateAt = now
	expense.CreateBy = userID
	expense.ModifiedAt = now
	expense.ModifiedBy = userID

	if err := c.Repo.Create(&expense); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"expenses": expense})
}

func (c *ExpenseController) GetAllExpense(ctx *gin.Context) {
	usrIdParams := ctx.Param("userId")
	usrId, err := strconv.Atoi(usrIdParams)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
	}
	expense, err := c.Repo.GetUserExpense(usrId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"expsenses": expense,
	})
}

func (c *ExpenseController) GetUserExpenseByDate(ctx *gin.Context) {
	userId := ctx.Param("userId")
	usrId, err := strconv.Atoi(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if userId == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user id is required"})
		return
	}

	var expenseRange models.DateRangeExpenses
	if err := ctx.ShouldBindJSON(&expenseRange); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	if expenseRange.StartDate == "" || expenseRange.EndDate == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}

	startDate, err := time.Parse("2006-01-02", expenseRange.StartDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", expenseRange.EndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	if endDate.Before(startDate) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after or equal to start_date"})
		return
	}

	expenses, err := c.Repo.GetExpenseByDate(usrId, startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expenses", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":     userId,
		"start_date":  expenseRange.StartDate,
		"end_date":    expenseRange.EndDate,
		"total_items": len(expenses),
		"data":        expenses,
	})
}

func (c *ExpenseController) TotalExpenseByUser(ctx *gin.Context) {
	userId, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		t, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
			return
		}
		startDate = &t
	}

	if endDateStr != "" {
		t, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
			return
		}
		endDate = &t
	}

	total, err := c.Repo.GetTotalExpenseByUserWithDateRange(int(userId), startDate, endDate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
	})
}

func (c *ExpenseController) UpdateExpenseByID(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	idParams := ctx.Param("id")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var expense models.Expenses
	if err := ctx.BindJSON(&expense); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expenseTime := expense.ExpenseDate.ToTime()
	expense.ID = int64(id)
	expense.ModifiedBy = userID
	expense.CreateBy = userID
	expense.ModifiedAt = time.Now()

	updatedExpense, err := c.Repo.UpdateUserExpense(id, expense, expenseTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"update": updatedExpense})
}

func (c *ExpenseController) DeleteExpenses(ctx *gin.Context) {
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

	if err := c.Repo.DeleteExpenseByID(int64(id)); err != nil {
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
