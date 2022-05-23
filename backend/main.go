package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/Ahm3dRN/go-react-todo/middlewares"
	"github.com/Ahm3dRN/go-react-todo/models"
	"github.com/Ahm3dRN/go-react-todo/utils"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func taskListHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.TaskList
	models.DB.Find(&tasks).Where("deleted_at = ?", "null")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	en := json.NewEncoder(w)
	en.SetIndent("", "    ")
	en.Encode(tasks)
}

// CreateTaskHandler Handles POST requests to /create-task/
// currently takes 3 form fields: task_list id, titlte, description
func createTaskHandler(w http.ResponseWriter, r *http.Request) {

	taskListFormID := r.FormValue("task_list_id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	order := r.FormValue("order") // if the order is supplied then use it otherwise use the len of the TaskLisk childs +1
	validationErrors := map[string]string{"ok": "false"}

	if taskListFormID == "" {
		fmt.Println("task_list_id can't be empty")
		validationErrors["task_list_id"] = "task_list_id can't be empty."

	}
	if title == "" {
		fmt.Println("Title can't be empty")
		validationErrors["title"] = "Title can't be empty."
	}
	if description == "" {
		fmt.Println("description can't be empty")
		validationErrors["description"] = "description can't be empty."
	}
	if order == "" {
		fmt.Println("use current order")
	}
	if len(validationErrors) > 1 {

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		jData, err := json.Marshal(validationErrors)
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jData)
		return
	} else {
		taskListFormID := strings.ReplaceAll(taskListFormID, "\"", "")
		Taa, err := strconv.ParseInt(taskListFormID, 10, 32)

		var taskList models.TaskList
		models.DB.First(&taskList, Taa)

		newTask := models.Task{Title: title, Description: description, IsComplete: false, IsDeleted: false, TaskListID: Taa}
		models.DB.Model(&taskList).Association("Tasks").Append(&newTask)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		jData, err := json.Marshal(map[string]string{"ok": "true"})
		if err != nil {
			println("errr")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(jData)
	}
}
func createTaskListHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	description := r.FormValue("description")
	validationErrors := map[string]string{"ok": "false"}
	if title == "" {
		validationErrors["title"] = "Title can't be empty"
	}
	if description == "" {
		validationErrors["description"] = "description can't be empty"
	}
	if len(validationErrors) > 1 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		jData, err := json.Marshal(validationErrors)
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jData)
	} else {
		TaskListReturn := models.TaskList{Title: title, Description: description, Count: 0}
		models.DB.Create(&TaskListReturn)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		jData, err := json.Marshal(map[string]string{"ok": "true"})
		if err != nil {
			println("errr")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(jData)
	}
}
func taskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["id"]
	if utils.CheckIfTaskExists(taskID) {
		var task models.Task
		models.DB.First(&task, taskID)

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(task)
	} else {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"ok": "false", "message": "task not found"})
	}
}
func taskListGetHandler(w http.ResponseWriter, r *http.Request) {
	taskListID := mux.Vars(r)["id"]
	if utils.CheckIfTaskListExists(taskListID) {
		var tasklist models.TaskList
		models.DB.First(&tasklist, taskListID)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(tasklist)
	} else {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"ok": "false", "message": "task list not found"})
	}
}
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new user has been registered")
	username := r.FormValue("username")
	password := r.FormValue("password")
	registerErrors := map[string]string{}
	if len(username) < 5 {
		fmt.Println("username must be at least 5 characters")
		registerErrors["username"] = "username must be at least 5 characters"
	}
	if len(password) < 8 {
		fmt.Println("password must be at least 8 characters")
	}
	if len(registerErrors) > 0 {
		registerErrors["ok"] = "false"
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(registerErrors)
		return
	}
	if utils.CheckUserExists(username) {
		fmt.Println("username already exists")
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]string{"ok": "false", "username": "username already exists"})
		return
	}
	hashed_pass := utils.GetHash([]byte(password))
	new_user := models.User{Username: username, PassowrdHash: hashed_pass}
	models.DB.Create(&new_user)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	en := json.NewEncoder(w)
	en.SetIndent("", "    ")
	en.Encode(map[string]string{"ok": "true"})
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if utils.CheckUserExists(username) {
		var user models.User
		models.DB.Where("username = ?", username).First(&user)
		fmt.Println(user.PassowrdHash)
		fmt.Println(password)
		if utils.CheckPassHash(user.PassowrdHash, password) {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			en := json.NewEncoder(w)
			en.SetIndent("", "    ")
			en.Encode(map[string]string{"ok": "true", "message": "you're logged in!"})
		} else {
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusUnprocessableEntity)
			en := json.NewEncoder(w)
			en.SetIndent("", "    ")
			en.Encode(map[string]string{"ok": "false", "password": "password isn't valid"})
		}
	} else {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]string{"ok": "false", "username": "username doesn't exist"})
	}
}
func editTaskListHandler(w http.ResponseWriter, r *http.Request) {
	taskListId := r.FormValue("task_list_id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	testi := map[string]interface{}{}
	errorMap := map[string]string{"ok": "false"}
	if taskListId == "" {
		errorMap["task_list_id"] = "task list id can't be empty"
	} else {
		obj := models.DB.Debug().Model(models.TaskList{}).Where("id = ?", taskListId)
		err := obj.Error
		if err != nil {
			errorMap["task_list_id"] = "Task list is not found."
		}
	}
	if title != "" {
		testi["title"] = title
	}
	if description != "" {
		testi["description"] = description
	}

	if len(errorMap) > 1 {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(errorMap)
	} else {
		re := models.DB.Debug().Model(&models.TaskList{}).Where("id = ?", taskListId).Updates(&testi)
		if re.Error != nil {
			fmt.Println(re.Error)
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]string{"ok": "true"})
	}
}
func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskListId := r.FormValue("task_list_id")
	taskId := r.FormValue("task_id")
	title := r.FormValue("title")
	description := r.FormValue("description")
	errorMap := map[string]string{"ok": "false"}
	testi := map[string]interface{}{}
	if taskListId == "" {
		errorMap["task_list_id"] = "task list id can't be empty"
	}
	if taskId == "" {
		errorMap["task_id"] = "task id can't be empty"
	}
	if title != "" {
		testi["title"] = title
	}
	if taskId != "" {
		testi["description"] = description
	}
	if len(errorMap) > 1 {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(errorMap)
	} else {
		re := models.DB.Debug().Model(&models.Task{}).Where("id = ?", taskId).Where("task_list_id = ?", taskListId).Updates(&testi)
		if re.Error != nil {
			fmt.Println(re.Error)
		}
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]string{"ok": "true"})
	}
}
func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", homeHandler).Name("home")
	r.HandleFunc("/tasks/", middlewares.Chain(taskListHandler, middlewares.Logging())).Name("task-list")

	r.HandleFunc("/tasklist/{id}", taskListGetHandler).Name("get-task-list").Methods("GET")
	r.HandleFunc("/create-task/", createTaskHandler).Name("create-task").Methods("POST")
	r.HandleFunc("/create-tasklist/", createTaskListHandler).Name("create-task-list").Methods("POST")
	r.HandleFunc("/task/{id}", taskHandler).Name("task").Methods("GET")
	r.HandleFunc("/edit-tasklist/", editTaskListHandler).Name("task").Methods("POST")
	r.HandleFunc("/edit-task/", editTaskHandler).Name("task").Methods("POST")
	r.HandleFunc("/users/register/", registerHandler).Name("register").Methods("GET", "POST")
	r.HandleFunc("/users/login/", loginHandler).Name("login").Methods("GET", "POST")

	fmt.Println("Running Server on 0.0.0.0:80.")
	fmt.Println("Server is now running.")

	http.ListenAndServe(":80", r)
}
