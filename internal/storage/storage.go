package storage


type Storage interface{
	CreateUser(userid int,name string,email string)(int64,error)
}