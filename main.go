package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("mysql", "root:zhangrun@tcp(47.103.90.59:3306)/gin_api?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("can not connect to database")
	}

	db.AutoMigrate(&todoModel{})
}

type (
	todoModel struct {
		gorm.Model
		Title     string `json:"title"`
		Completed int    `json:"completed"`
	}

	transformedTodo struct {
		ID        uint   `JSON:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
)

func main() {
	router := gin.Default()
	v1 := router.Group("/api/v1/todos")
	{
		v1.GET("/", fetchAllTodo)
		v1.POST("/", createTodo)
		v1.GET("/:id", fetchSingeTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}
	router.Run()
}

func fetchAllTodo(c *gin.Context) {
	var todos []todoModel
	var _todos []transformedTodo

	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"msg":    "No todo found!",
		})
		return
	}

	for _, item := range todos {
		completed := true
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, transformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

func createTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := todoModel{Title: c.PostForm("title"), Completed: completed}

	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "msg": "scucess created", "resourseId": todo.ID})
}

func fetchSingeTodo(c *gin.Context) {
	var todo todoModel
	todoId := c.Param("id")

	db.First(&todo, todoId)
	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "msg": "Not Found"})
		return
	}

	completed := true
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := transformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

func updateTodo(c *gin.Context) {
	var todo todoModel
	itemID := c.Param("id")

	db.First(&todo, itemID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "msg": "no todo found"})
		return
	}

	title := c.PostForm("title")
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	data := make(map[string]interface{})
	data["title"] = title
	data["completed"] = completed
	db.Model(&todo).Updates(data)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "success updated"})
}

func deleteTodo(c *gin.Context) {
	var todo todoModel
	todoId := c.Param("id")

	db.First(&todo, todoId)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "msg": "no todo found"})
		return
	}

	db.Delete(&todo)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "msg": "todo deleted success"})
}
