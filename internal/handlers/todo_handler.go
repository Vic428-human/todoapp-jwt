package handlers

import (
	"net/http"
	"strconv"
	"todo_api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreateTodoRequest struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoRequest struct {
	Title string `json:"title"`
	// 需要把 Completed 定義成 指標型別，才可以和 nil 比較。
	Completed *bool `json:"completed"`
}

/*
	{
	    "title":"buy new book02",
	    "completed": false
	}
*/
func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) { // 閉包
		var input CreateTodoRequest

		//  先驗證從 client 傳來的資料
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// 沒問題後，把資料傳給 repository 層，透過 sql 方式把資料寫入到DB
		todo, err := repository.CreateTodo(pool, input.Title, input.Completed)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		// 資料庫寫入正確後，回傳訊息到 client 端
		c.JSON(http.StatusCreated, todo)
	}
}

func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		todos, err := repository.GetAllTodos(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, todos)
	}
}

func GetTodoByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		// "2" -----> 2, nil
		// 'a' ----> 0,error ("invalid syntax")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID TODO ID"})
			return
		}
		todo, err := repository.GetTodoByID(pool, id)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "todo not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, todo)

	}
}

/*
	{
	    "title":"buy new book02-修改",
	    "completed": true
	}
*/
func UpdateToDoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID TODO ID"})
			return
		}
		var input UpdateTodoRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Go 裡只有指標、slice、map、interface、channel、function 這些類型可以和 nil 比較。
		if input.Title == "" && input.Completed == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title or completed is required"})
			return
		}

		var completed bool
		// 檢查 Completed 是否為 true or false
		if input.Completed != nil && *input.Completed {
			completed = *input.Completed
		}

		todo, err := repository.UpdateTodo(pool, id, input.Title, completed)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "todo not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, todo)

	}
}
