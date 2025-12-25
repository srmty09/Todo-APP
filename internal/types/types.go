package types

import "time"


type TaskMetaData struct{
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Completed bool `json:"completed"`
	created_at time.Time 
	updated_at time.Time
}

type User struct{
	UserId int `validate:"required"`
	Name string `validate:"required"`
	Email string `validate:"required"`
}