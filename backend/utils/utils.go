package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Ahm3dRN/go-react-todo/models"
	"golang.org/x/crypto/bcrypt"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Println("Error Found:")
		fmt.Println(err)
	}
}

func CheckUserExists(username string) bool {
	var user models.User
	result := models.DB.Where("username = ?", username).First(&user)
	return result.Error == nil
}
func GetHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
func CheckPassHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println(err)
	return err == nil
}

func CheckIfTaskExists(taskId string) bool {
	var task models.Task
	models.DB.First(&task, taskId)
	return task.ID != 0
}

func CheckIfTaskListExists(taskListID string) bool {
	var tasklist models.TaskList
	models.DB.First(&tasklist, taskListID)
	return tasklist.ID != 0
}
func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
func IsValidToken(token string) (models.Token, error) {
	token = strings.TrimPrefix(token, "\"")
	token = strings.TrimSuffix(token, "\"")
	tokenObj := models.Token{}
	models.DB.Debug().Where("token_hash = ?", token).Find(&tokenObj)
	fmt.Println(tokenObj)
	if tokenObj.ID == 0 {
		return tokenObj, errors.New("token not found")

	}
	return tokenObj, nil
}
