package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

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
	fmt.Println("this")
	r.ParseForm()
	taskListFormID := r.PostFormValue("task_list_id")
	title := r.PostFormValue("title")
	description := r.PostFormValue("description")
	// token := r.PostFormValue("token")
	// userID := r.Header.Get("UserID")
	// if the order is supplied then use it otherwise use the len of the TaskLisk childs +1

	validationErrors := map[string]string{"ok": "false"}

	if taskListFormID == "" {
		fmt.Println("task_list_id can't be empty")
		validationErrors["task_list_id"] = "task_list_id can't be empty."

	}
	if title == "" {
		fmt.Println("Title can't be empty")
		validationErrors["title"] = "Title can't be empty."
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
		if err != nil {
			fmt.Println(err)
		}
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
	r.ParseForm()
	title := r.PostFormValue("title")
	description := r.PostFormValue("description")
	token := r.Header.Get("Authentication")
	user, err := utils.IsValidToken(token)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]string{"ok": "false", "token": "token is not valid"})
		return
	}
	fmt.Println(user.ID)
	validationErrors := map[string]string{"ok": "false"}
	if title == "" {
		validationErrors["title"] = "Title can't be empty"
	}
	// if description == "" {
	// 	validationErrors["description"] = "description can't be empty"
	// }
	if len(validationErrors) > 1 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		jData, err := json.Marshal(validationErrors)
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(jData)
	} else {

		TaskListReturn := models.TaskList{Title: title, Description: description, Count: 0, UserID: user.ID}
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
	fmt.Println(r.Method)
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Username : ", r.PostFormValue("username"))
	fmt.Println("Password : ", r.PostFormValue("password"))
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	registerErrors := map[string]string{"ok": "false"}
	if len(username) < 5 {
		fmt.Println("username must be at least 5 characters")
		registerErrors["username"] = "username must be at least 5 characters"
	}
	if len(password) < 8 {
		fmt.Println("password must be at least 8 characters")
		registerErrors["password"] = "password must be at least 8 characters"
	}
	if len(registerErrors) > 1 {
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
	fmt.Println(r.Method)
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Username : ", r.PostFormValue("username"))
	fmt.Println("Password : ", r.PostFormValue("password"))
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if utils.CheckUserExists(username) {
		var user models.User
		models.DB.Where("username = ?", username).First(&user)
		fmt.Println(user.PassowrdHash)
		fmt.Println(password)
		if utils.CheckPassHash(user.PassowrdHash, password) {
			token_hash := utils.GenerateToken()
			const timeLayout = "2006-01-02 15:04:05"
			dt := time.Now()
			expirtyTime := time.Now().Add(time.Hour * 72)
			CreatedAt := dt.Format(timeLayout)
			ExpiryDate := expirtyTime.Format(timeLayout)
			token := models.Token{TokenHash: token_hash, ExpiryDate: ExpiryDate, CreatedAt: CreatedAt, UserID: user.ID}
			models.DB.Debug().Create(&token)
			// if token.ID != 0 {

			// }
			w.Header().Add("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			en := json.NewEncoder(w)
			en.SetIndent("", "    ")
			en.Encode(map[string]string{"ok": "true", "message": "you're logged in!", "token": token.TokenHash})
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
	fmt.Println("here edit task", r.Method)
	r.ParseForm()
	taskListId := r.PostFormValue("task_list_id")
	taskId := r.PostFormValue("task_id")
	title := r.PostFormValue("title")
	description := r.PostFormValue("description")
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
func getTaskItemsHandler(w http.ResponseWriter, r *http.Request) {
	task_list_id := r.FormValue("task_list_id")
	fmt.Println(task_list_id)
	errorsMap := map[string]string{"ok": "false"}
	if task_list_id == "" {
		errorsMap["task_list_id"] = "task_list_id can't be empty"
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(errorsMap)
		return
	}
	var taskItems []models.Task
	models.DB.Debug().Model(&models.Task{}).Where("task_list_id = ?", task_list_id).Find(&taskItems)
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	en := json.NewEncoder(w)
	en.SetIndent("", "    ")
	en.Encode(taskItems)
}
func checkTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskListID := r.PostFormValue("task_list_id")
	taskID := r.PostFormValue("task_id")
	errorsMap := map[string]string{"ok": "false"}
	if taskListID == "" {
		errorsMap["task_list_id"] = "task_list_id can't be empty"
	}
	if taskID == "" {
		errorsMap["task_id"] = "task_id can't be empty"
	}
	if len(errorsMap) > 1 {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(errorsMap)
		return
	}
	var task models.Task
	models.DB.Debug().Model(&models.Task{}).Where("id = ?", taskID).Where("task_list_id = ?", taskListID).Find(&task)
	if task.ID == 0 {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]string{"ok": "false", "task_id": "task is not found"})
		return
	} else {
		print(task.IsComplete)
		if task.IsComplete {
			task.IsComplete = false
		} else {
			task.IsComplete = true
		}

		models.DB.Debug().Model(&task).Updates(&task)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]interface{}{"ok": "true", "todo": task})
	}
}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskListID := r.PostFormValue("task_list_id")
	taskID := r.PostFormValue("task_id")
	errorsMap := map[string]string{"ok": "false"}
	if taskListID == "" {
		errorsMap["task_list_id"] = "task_list_id can't be empty"
	}
	if taskID == "" {
		errorsMap["task_id"] = "task_id can't be empty"
	}
	if len(errorsMap) > 1 {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(errorsMap)
		return
	}
	var task models.Task
	models.DB.Debug().Model(&models.Task{}).Where("id = ?", taskID).Where("task_list_id = ?", taskListID).Find(&task)
	if task.ID == 0 {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]string{"ok": "false", "task_id": "task is not found"})
		return
	} else {
		models.DB.Debug().Model(&task).Delete(&task)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]interface{}{"ok": "true"})
	}
}
func deleteTaskListHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	task_list_id := r.PostFormValue("task_list_id")
	errorsMap := map[string]string{"ok": "false"}
	if task_list_id == "" {
		errorsMap["task_list_id"] = "task_list_id can't be empty"
	}
	if len(errorsMap) > 1 {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnprocessableEntity)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(errorsMap)
		return
	}
	var taskList models.TaskList
	models.DB.Debug().Model(&models.TaskList{}).Where("id = ?", task_list_id).Find(&taskList)
	if taskList.ID == 0 {
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]string{"ok": "false", "task_list_id": "task list is not found"})
		return
	} else {
		models.DB.Debug().Model(&taskList).Delete(&taskList)
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		en := json.NewEncoder(w)
		en.SetIndent("", "    ")
		en.Encode(map[string]interface{}{"ok": "true"})
	}
}
func reactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	http.ServeFile(w, r, "./public/index.html")

}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", reactHandler).Name("home")
	r.HandleFunc("/tasklists/", middlewares.Chain(taskListHandler, middlewares.Logging())).Name("task-list")
	r.HandleFunc("/tasks/", middlewares.Chain(getTaskItemsHandler, middlewares.Logging())).Name("get-tasks")
	r.HandleFunc("/tasklist/{id}", taskListGetHandler).Name("get-task-list").Methods("GET")
	r.HandleFunc("/create-task/", createTaskHandler).Name("create-task").Methods("POST")
	r.HandleFunc("/create-tasklist/", createTaskListHandler).Name("create-task-list").Methods("POST")
	r.HandleFunc("/task/{id}", taskHandler).Name("task").Methods("GET")
	r.HandleFunc("/edit-tasklist/", editTaskListHandler).Name("edit-task-list").Methods("POST")
	r.HandleFunc("/edit-task/", editTaskHandler).Name("edit-task").Methods("POST")
	r.HandleFunc("/check-task/", checkTaskHandler).Name("check-task").Methods("POST")
	r.HandleFunc("/delete-task/", deleteTaskHandler).Name("delete-task").Methods("POST")
	r.HandleFunc("/delete-tasklist/", deleteTaskListHandler).Name("delete-task").Methods("POST")
	r.HandleFunc("/users/register/", registerHandler).Name("register").Methods("POST")
	r.HandleFunc("/users/login/", loginHandler).Name("login").Methods("POST")

	fmt.Println("Running Server on 0.0.0.0:80.")
	fmt.Println("Server is now running.")
	r.Use(middlewares.TokenAuthenticationMiddleware)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authentication", "Content-type"},
	})
	handler := c.Handler(r)

	http.ListenAndServe(":80", handler)
}
