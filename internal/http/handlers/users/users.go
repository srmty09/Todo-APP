package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/srmty09/Todo-App/internal/storage"
	"github.com/srmty09/Todo-App/internal/types"
	"github.com/srmty09/Todo-App/internal/utils/helpers"
	"github.com/srmty09/Todo-App/internal/utils/response"
)

//create new user
func New(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a user")
		var user types.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if errors.Is(err,io.EOF){
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		if err!=nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return
		}
		validate := validator.New()
		if err := validate.Struct(user);err!=nil{
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}
		lastId,err:=storage.CreateUser(user.Name,user.Email)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return
		}
		slog.Info("user created successfully",slog.String("userId",fmt.Sprint(lastId)))
		response.WriteJson(w,http.StatusCreated,map[string]interface{}{
			"status": "OK",
			"id": lastId,
		})
	}
}

func GetUserInfo(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		intId,err := helpers.ParsePathInt64(r,"id")
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return 
		}
	slog.Info("getting user info for", slog.Int64("userId", intId))
	exist,err := storage.UserExists(intId)
	if err!=nil{
		response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
		return
	}
	if !exist{
		response.WriteJson(w,http.StatusNotFound,response.GeneralError(fmt.Errorf("user with id %d does not exist", intId)))
		return 
	}
		user,err := storage.GetUser(intId)
		if err!=nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return 
		}
		response.WriteJson(w,http.StatusOK,user)
	}
}

func DeleteUserInfo(storage storage.Storage)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		userId,err := helpers.ParsePathInt64(r,"id")
		if err!= nil{
			response.WriteJson(w,http.StatusBadRequest,response.GeneralError(err))
			return  // ← Missing return fixed!
		}
		slog.Info("deleting user with", slog.Int64("userId", userId))
		exist,err := storage.UserExists(userId)
		if err!= nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return 
		}
		if !exist{
			response.WriteJson(w,http.StatusNotFound,response.GeneralError(fmt.Errorf("user with id %d does not exist",userId)))
			return 
		}
		err = storage.DeleteUser(userId)
		if err!= nil{
			response.WriteJson(w,http.StatusInternalServerError,response.GeneralError(err))
			return 
		}
		response.WriteJson(w,http.StatusOK,map[string]interface{}{
			"status":"deleted",
			"userid":userId,
		})
	}
}