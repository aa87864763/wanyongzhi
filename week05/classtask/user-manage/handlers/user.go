package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"user-manage/models"

	"github.com/gin-gonic/gin"
)

var users []models.User

func loadUsers() error {
	data, err := ioutil.ReadFile("user.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &users)
}

func saveUsers() error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	return ioutil.WriteFile("user.json", data, 0644)
}

func GetUsers(c *gin.Context) {
	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func CreateUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查email是否已存在
	for _, user := range users {
		if user.Email == newUser.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}
	}

	users = append(users, newUser)
	if err := saveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func UpdateUser(c *gin.Context) {
	var updatedUser models.User
	if err := c.ShouldBindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查找并更新用户
	found := false
	for i, user := range users {
		if user.Email == updatedUser.Email {
			users[i] = updatedUser
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := saveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func DeleteUser(c *gin.Context) {
	email := c.Param("email")

	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 查找并删除用户
	found := false
	for i, user := range users {
		if user.Email == email {
			users = append(users[:i], users[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := saveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func SearchUser(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email parameter is required"})
		return
	}

	if err := loadUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	for _, user := range users {
		if user.Email == email {
			c.JSON(http.StatusOK, user)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
}
