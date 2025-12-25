package response

import (
	"encoding/json"
	"net/http"

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



