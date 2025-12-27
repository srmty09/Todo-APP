package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/srmty09/Todo-App/internal/storage"
	"github.com/srmty09/Todo-App/internal/types"
	"github.com/srmty09/Todo-App/internal/utils/helpers"
	"github.com/srmty09/Todo-App/internal/utils/response"
)


func Add(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helpers.ParsePathInt64(r, "id")
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		} 
		slog.Info("Adding new task to user",slog.Int64("userId", userId))
		var task types.TaskMetaData
		err = json.NewDecoder(r.Body).Decode(&task)
		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err != nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		validate := validator.New()
		if err := validate.Struct(task);err!=nil{
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError(validateErrs))
			return 
		}
		// Check if user exists
		exists,err := storage.UserExists(userId)
		if err != nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		if !exists{
			response.WriteJson(w,http.StatusNotFound,response.GeneralError(fmt.Errorf("user with id %d does not exist",userId)))
			return
		}
		lastId,err := storage.AddNewTask(userId,task.Title,task.Description,task.Priority,task.Completed,time.Now(),time.Now())
		if err != nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		slog.Info("task added successfully",slog.String("taskId",fmt.Sprint(lastId)))
		response.WriteJson(w,http.StatusCreated,map[string]interface{}{
			"status": "OK",
			"id": lastId,
		})
	}
}


func GetTodo(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helpers.ParsePathInt64(r, "id")
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}
		
		status := r.URL.Query().Get("status")
		keyword := r.URL.Query().Get("search")
		
		slog.Info("Getting tasks for user", slog.Int64("userId", userId), slog.String("status", status))
		
		exist,err:= storage.UserExists(userId)
		if err!= nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return 
		}
		if !exist{
			response.WriteJson(w,http.StatusNotFound,response.GeneralError(fmt.Errorf("user with id %d does not exist",userId)))
			return 
		}

	var tasks []types.TaskMetaData
	
	if keyword != "" && status != "" {
		tasks, err = storage.GetTaskWithFilters(userId, "%"+keyword+"%", status)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		if tasks == nil {
			response.WriteJson(w, http.StatusOK, map[string]interface{}{
				"message": fmt.Sprintf("No %s tasks found matching '%s'", status, keyword),
				"tasks":   []types.TaskMetaData{},
			})
			return
		}
	} else if keyword != "" {
		tasks, err = storage.GetTaskWithTitle(userId, "%"+keyword+"%")
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		if tasks == nil {
			response.WriteJson(w, http.StatusOK, map[string]interface{}{
				"message": "No task found",
				"tasks":   []types.TaskMetaData{},
			})
			return
		}
	} else {
		switch status {
		case "completed":
			tasks, err = storage.GetCompletedTask(userId)
			if err != nil {
				response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
				return
			}
			if tasks == nil {
				response.WriteJson(w, http.StatusOK, map[string]interface{}{
					"message": "No completed tasks found",
					"tasks":   []types.TaskMetaData{},
				})
				return
			}
		case "incomplete", "incompleted":
			tasks, err = storage.GetIncompletedTask(userId)
			if err != nil {
				response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
				return
			}
			if tasks == nil {
				response.WriteJson(w, http.StatusOK, map[string]interface{}{
					"message": "All tasks completed!",
					"tasks":   []types.TaskMetaData{},
				})
				return
			}
		default:
			tasks, err = storage.GetTaskForId(userId)
			if err != nil {
				response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
				return
			}
		}
	}

		response.WriteJson(w,http.StatusOK,tasks)
	}
}



func CompletedTask(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helpers.ParsePathInt64(r, "id")
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}	
		taskId, err := helpers.ParsePathInt64(r, "task_id")
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}
		slog.Info("Marking task as complete", slog.Int64("userId", userId), slog.Int64("taskId", taskId))
		err = storage.MarkComplete(userId, taskId)
		if err!= nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return 
		}
		response.WriteJson(w,http.StatusOK,map[string]interface{}{
			"status": "Completed",
		})
	}
}


func IncompletedTask(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helpers.ParsePathInt64(r, "id")
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}	
		taskId, err := helpers.ParsePathInt64(r, "task_id")
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}
		slog.Info("Marking task as incomplete", slog.Int64("userId", userId), slog.Int64("taskId", taskId))
		err = storage.MarkIncomplete(userId, taskId)
		if err!= nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return 
		}
		response.WriteJson(w,http.StatusOK,map[string]interface{}{
			"status": "Incompleted",
		})
	}
}


func GetSingleTask(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helpers.ParsePathInt64(r, "id")
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}
		taskId, err := helpers.ParsePathInt64(r, "task_id")
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}
		slog.Info("Getting single task", slog.Int64("userId", userId), slog.Int64("taskId", taskId))
		task,err := storage.GetSingleTask(userId, taskId)
		if err!= nil{
			response.WriteJson(w,http.StatusNotFound,response.GeneralError(err))
			return 
		}
		response.WriteJson(w,http.StatusOK,task)
	}
}

func DeleteTask(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := helpers.ParsePathInt64(r, "id")
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}	
		taskId, err := helpers.ParsePathInt64(r, "task_id")
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}
		slog.Info("Deleting task", slog.Int64("userId", userId), slog.Int64("taskId", taskId))
		err = storage.DeletingTask(userId, taskId)
		if err!= nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return 
		}
		response.WriteJson(w,http.StatusOK,map[string]interface{}{
			"status": "Deleted",
		})
	}
}


func EditTask(storage storage.Storage)http.HandlerFunc{
	return  func(w http.ResponseWriter, r *http.Request) {
		userId, err := helpers.ParsePathInt64(r, "id")
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		taskId, err := helpers.ParsePathInt64(r, "task_id")
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}
		slog.Info("Editing task", slog.Int64("userId", userId), slog.Int64("taskId", taskId))
		
		// Get existing task first
		existingTask, err := storage.GetSingleTask(userId, taskId)
		if err!=nil{
			response.WriteJson(w,http.StatusNotFound,response.GeneralError(err))
			return
		}
		
		// Decode the update request
		var updateRequest types.TaskMetaData
		err = json.NewDecoder(r.Body).Decode(&updateRequest)
		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(fmt.Errorf("empty body")))
			return
		}
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		
		// Update only provided fields (partial update)
		if updateRequest.Title != "" {
			existingTask.Title = updateRequest.Title
		}
		if updateRequest.Description != "" {
			existingTask.Description = updateRequest.Description
		}
		if updateRequest.Priority != "" {
			existingTask.Priority = updateRequest.Priority
		}
		
		// Validate the final task
		validate := validator.New()
		if err:= validate.Struct(existingTask);err!=nil{
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w,http.StatusBadRequest,response.ValidationError(validateErrs))
			return 
		}
		
		err = storage.EditTask(userId,taskId,existingTask.Title,existingTask.Description,existingTask.Priority)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		slog.Info("task edited successfully", slog.Int64("userId", userId), slog.Int64("taskId", taskId))
		response.WriteJson(w,http.StatusOK,map[string]interface{}{
			"status": "Updated",
			"message": "Task updated successfully",
		})
	}
}