package tasks

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/srmty09/Todo-App/internal/storage"
	"github.com/srmty09/Todo-App/internal/types"
	"github.com/srmty09/Todo-App/internal/utils/response"
)


func Add(storage storage.StorageForTask) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		intId,err := strconv.ParseInt(id,10,64)
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		} 
		slog.Info("Adding new task to user",slog.String("userId",id))
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
		exists,err := storage.UserExists(intId)
		if err != nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		if !exists{
			response.WriteJson(w,http.StatusNotFound,response.GeneralError(fmt.Errorf("user with id %d does not exist",intId)))
			return
		}
		lastId,err := storage.AddNewTask(intId,task.Title,task.Description,task.Completed,time.Now(),time.Now())
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


