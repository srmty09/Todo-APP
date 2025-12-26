package storage

import (
	"time"

	"github.com/srmty09/Todo-App/internal/types"
)


type StorageForUser interface{
	CreateUser(name string,email string)(int64,error)
}

type StorageForTask interface{
	UserExists(userid int64)(bool,error)
	AddNewTask(userid int64,title string,description string,completed bool,created_at time.Time,updated_at time.Time)(int64,error)
	GetTaskForId(userid int64)([]types.TaskMetaData,error)
	MarkComplete(userid int64, taskid int64)(error)
}