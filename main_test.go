package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task/config"
	"task/models"
	"task/routes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, &models.Task{})
	router := gin.Default()
	routes.SetupRoutes(router)
	return router
}

func TestSignUp(t *testing.T) {
	router := setupRouter()
	user := models.User{
		Username: "testuser",
		Password: "password123",
	}
	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "Failed to create user")
}

func TestLogin(t *testing.T) {
	router := setupRouter()
	config.DB.Create(&models.User{
		Username: "testuser",
		Password: "$2a$10$somethinghashedhere",
	})

	credentials := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	credentialsJSON, _ := json.Marshal(credentials)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(credentialsJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Login successful")
}

func TestCreateTask(t *testing.T) {
	router := setupRouter()

	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Priority:    "High",
		Deadline:    "2024-12-31",
		Status:      "Pending",
		Category:    "Testing",
		UserID:      1,
	}
	taskJSON, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "Failed to create task")
}

func TestGetTasks(t *testing.T) {
	router := setupRouter()

	config.DB.Create(&models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Priority:    "High",
		Deadline:    "2024-12-31",
		Status:      "Pending",
		Category:    "Testing",
		UserID:      1,
	})

	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Task")
}

func TestUpdateTask(t *testing.T) {
	router := setupRouter()

	task := models.Task{
		Title:       "Old Task",
		Description: "This is an old task",
		Priority:    "Low",
		Deadline:    "2024-12-31",
		Status:      "Pending",
		Category:    "Old",
		UserID:      1,
	}
	config.DB.Create(&task)

	updatedTask := map[string]string{
		"title":       "Updated Task",
		"description": "This task has been updated",
		"priority":    "High",
		"status":      "Completed",
	}
	updatedTaskJSON, _ := json.Marshal(updatedTask)

	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBuffer(updatedTaskJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Contains(t, w.Body.String(), "Task not found")
}

func TestDeleteTask(t *testing.T) {
	router := setupRouter()

	task := models.Task{
		Title:       "Task to Delete",
		Description: "This task will be deleted",
		Priority:    "Medium",
		Deadline:    "2024-12-31",
		Status:      "Pending",
		Category:    "Delete",
		UserID:      1,
	}
	config.DB.Create(&task)

	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Task deleted successfully")
}
