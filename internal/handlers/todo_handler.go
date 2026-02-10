package handlers

import (
	"net/http"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoRequest struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateTodoRequest

		//  先驗證從 client 傳來的資料
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 沒問題後，把資料傳給 repo層，透過sql方式把資料寫入到DB
		todo, err := repository.CreateTodo(pool, input.Title, input.Completed)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// 資料庫寫入正確後，回傳訊息到 client 端
		c.JSON(http.StatusCreated, todo)
	}
}
