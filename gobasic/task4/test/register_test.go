package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/model"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/repository"
	"github.com/wangxg-sj/web3-learning/gobasic/task4/router"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := "root:1qaz@WSX@tcp(127.0.0.1:3306)/task4?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Ensure tables exist
	repository.AutoMigrate(db)

	return db
}

func TestRegister(t *testing.T) {
	db := setupTestDB()
	r := router.InitRouter(db)

	// Cleanup: Delete test user if exists
	testEmail := "test_register@example.com"
	db.Where("email = ?", testEmail).Delete(&model.User{})

	// Prepare request
	reqBody := map[string]string{
		"email":    testEmail,
		"password": "password123",
		"name":     "Test User",
	}
	jsonBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(jsonBytes))
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assertions
	// The util.SuccessResponse always returns 200 OK, with the actual status code in the body
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	// Response code is float64 because of json unmarshal default
	assert.Equal(t, float64(http.StatusCreated), response["code"])

	// Verify database
	var user model.User
	result := db.Where("email = ?", testEmail).First(&user)
	assert.NoError(t, result.Error)
	assert.Equal(t, "Test User", user.Name)

	// Update: Cleanup after test
	db.Delete(&user)
}
