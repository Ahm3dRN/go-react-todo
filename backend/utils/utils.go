package utils

import (
	"fmt"
	"log"

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
	// if task.ID == 0 {
	// 	return false
	// }
	// return true
	return task.ID != 0
}

func CheckIfTaskListExists(taskListID string) bool {
	var tasklist models.TaskList
	models.DB.First(&tasklist, taskListID)
	return tasklist.ID != 0
}
