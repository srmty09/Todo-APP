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
	"github.com/srmty09/Todo-App/internal/utils/response"
)

//create new user
func New(storage storage.StorageForUser) http.HandlerFunc{
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
		lastId,err:=storage.CreateUser(user.UserId,user.Name,user.Email)
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