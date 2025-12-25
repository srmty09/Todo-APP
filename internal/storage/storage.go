package storage

import "time"


type StorageForUser interface{
	CreateUser(userid int,name string,email string)(int64,error)
}

type StorageForTask interface{
	UserExists(userid int64)(bool,error)
	AddNewTask(userid int64,title string,description string,completed bool,created_at time.Time,updated_at time.Time)(int64,error)
}