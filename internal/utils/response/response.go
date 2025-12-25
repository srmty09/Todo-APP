package response

import (
	"encoding/json"
	"net/http"
	"github.com/go-playground/validator/v10"
	"strings"
	"fmt"

)



type Response struct{
	Status string `json:"status"`
	Error string `json:"error"`
}

type Request struct{
	Id int `json:"id"`
	Msg string `json:"msg"`
}

const (
	StatusOK = "OK"
	StatusError = "error"
)



// for writing json files
func WriteJson(w http.ResponseWriter, status int, data interface{}) error{
	w.Header().Set("Content-Type","application/json")

	// first we set the status!!!
	w.WriteHeader(status)

	//then we sent the body or the data!
	return json.NewEncoder(w).Encode(data)
}




func ReadRequest(w http.ResponseWriter, r *http.Request) {
	var req Request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid json",
		})
		return
	}
	if req.Msg == "" {
		WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "msg field is required",
		})
		return
	}

	WriteJson(w, http.StatusOK, req)
}

func GeneralError(err error) Response{
	return Response{
		Status: StatusError,
		Error: err.Error(),
	}
}


func ValidationError(errs validator.ValidationErrors) Response{
	var errMsgs []string

	for _,err := range errs{
		switch err.ActualTag(){
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field",err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid",err.Field()))
		}
	}
	return Response{
		Status: StatusError,
		Error: strings.Join(errMsgs,","),
	}
}

