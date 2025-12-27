package types

import "time"


type TaskMetaData struct{
	Title string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Priority string `json:"priority" validate:"required,oneof=low medium high"`
	Completed bool `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct{
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}