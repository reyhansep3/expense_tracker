package routers

import (
	"database/sql"
	"exp_tracker/controllers"
	"exp_tracker/middleware"
	"exp_tracker/repository"

	"github.com/gin-gonic/gin"
)

func StartServer(db *sql.DB) {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	userRepo := repository.NewUserRepository(db)
	userController := controllers.UserController{
		Repo: userRepo,
	}

	categoryRepo := repository.NewCategoryRepository(db)
	categoryController := controllers.CategoryController{
		Repo: categoryRepo,
	}

	expenseRepo := repository.CExpenseRepository(db)
	expenseController := controllers.ExpenseController{
		Repo: expenseRepo,
	}

	targetRepo := repository.CTargetRepository(db)
	targetRepository := controllers.TargetController{
		Repo: targetRepo,
	}

	router.POST("/api/user", userController.CreateUsers)

	protected := router.Group("/api")
	protected.Use(middleware.ValidateUser())
	{
		protected.POST("/categories", categoryController.CreateCategory)
		protected.GET("/users/:userId/categories", categoryController.GetAllCategory)
		protected.GET("/users/:userId/categories/:id", categoryController.GetCategoryByID)
		protected.DELETE("/users/:userId/categories/:id", categoryController.DeleteCategoryByID)

		protected.POST("/expenses", expenseController.AddExpenses)
		protected.POST("/expenses/:userId/ByDateRange", expenseController.GetUserExpenseByDate)
		protected.GET("/users/:userId/expenses", expenseController.GetAllExpense)
		protected.GET("/expenses/total", expenseController.TotalExpenseByUser)
		protected.PUT("/expenses/:id", expenseController.UpdateExpenseByID)
		protected.DELETE("/expenses/:id", expenseController.DeleteExpenses)

		protected.POST("/target", targetRepository.AddTarget)
		protected.GET("/target/:userId", targetRepository.GetAllTarget)
		protected.PUT("/target/:id", targetRepository.UpdateTargetExpense)
		protected.DELETE("/target/:id", targetRepository.DeleteTarget)

	}

	var PORT = "8080"
	router.Run(":" + PORT)
}
