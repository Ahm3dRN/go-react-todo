package models

import (
	"fmt"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB, err = gorm.Open(sqlite.Open("main.db"), &gorm.Config{})

type TaskList struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	Count       int    `json:"count"`
	UserID      uint   `json:"-"`
	// User        User
	Tasks     []Task
	CreatedAt datatypes.Date `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt datatypes.Date `json:"updated_at" gorm:"autoUpdateTime:true"`
}

type Task struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	// State       string `json:"state"`
	TaskListID int64  `json:"task_list"`
	Order      uint32 `json:"order"`
	IsComplete bool   `json:"is_complete" gorm:"default:false"`
	IsDeleted  bool   `json:"-" gorm:"default:false"`
}

// func (t *Task) AfterUpdate(tx *gorm.DB) (err error) {
// 	var taskList TaskList
// 	tx.Model(&taskList).Where("ID = ?", t.TaskListID).Update("count", )
// }

type User struct {
	gorm.Model
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email" gorm:"uniqueIndex"`
	Username     string `json:"user_name" gorm:"uniqueIndex"`
	PassowrdHash string `json:"-"`
	Role         uint8  `json:"role"`
}

type Token struct {
	ID         uint `gorm:"primarykey"`
	TokenHash  string
	ExpiryDate string
	CreatedAt  string
	UserID     uint
	// User        User
}

func CheckErr(err error) {
	if err != nil {
		fmt.Println("Error Found:")
		fmt.Println(err)
	}
}

func init() {
	CheckErr(err)
	DB.AutoMigrate(&TaskList{})
	DB.AutoMigrate(&Task{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Token{})
}
