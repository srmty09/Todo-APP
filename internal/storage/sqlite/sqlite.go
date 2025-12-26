package sqlite

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/srmty09/Todo-App/internal/config"
	"github.com/srmty09/Todo-App/internal/types"
)


type Sqlite struct{
	Db *sql.DB
}

func NewUserTb(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.Storage_path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT
	)`)

	if err != nil {
		return nil, err 
	}
	return &Sqlite{
		Db: db,
	}, nil
}

func NewTODOTb(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.Storage_path)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todo(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	title TEXT,
	description TEXT,
	completed INTEGER DEFAULT 0,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		return nil, err 
	}
	return &Sqlite{
		Db: db,
	}, nil
}



func (s *Sqlite) CreateUser(name string,email string)(int64, error){
	stmt,err:= s.Db.Prepare(
		"INSERT INTO user (name,email) VALUES(?,?)")
	if err!=nil{
		return 0,err 
	}
	defer stmt.Close()

	res,err := stmt.Exec(name,email)
	if err!=nil{
		return 0,err 
	}
	
	id,err:=res.LastInsertId()
	
	if err!=nil{
		return 0,err 
	}

	return id,nil
}

func (s *Sqlite) UserExists(userid int64)(bool,error){
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM user WHERE id = ?)"
	err := s.Db.QueryRow(query, userid).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *Sqlite) AddNewTask(userid int64,title string,description string,completed bool,created_at time.Time,updated_at time.Time)(int64, error){
	var completedInt int
	if completed {
		completedInt = 1
	} else {
		completedInt = 0
	}
	
	stmt,err:= s.Db.Prepare(
		"INSERT INTO todo (user_id,title,description,completed,created_at,updated_at) VALUES(?,?,?,?,?,?)")
	if err!=nil{
		return 0,err 
	}
	defer stmt.Close()

	res,err := stmt.Exec(userid,title,description,completedInt,created_at,updated_at)
	if err!=nil{
		return 0,err 
	}
	
	id,err:=res.LastInsertId()
	
	if err!=nil{
		return 0,err 
	}

	return id,nil
}


func (s *Sqlite)GetTaskForId(userid int64) ([]types.TaskMetaData,error){
	stmt,err:= s.Db.Prepare("SELECT title,description,completed,created_at,updated_at FROM todo WHERE user_id = ?")
	if err!= nil{
		return []types.TaskMetaData{},err
	}
	defer stmt.Close()
	rows,err := stmt.Query(userid)
	if err!= nil{
		return []types.TaskMetaData{},err
	}
	defer rows.Close()
	var tasks []types.TaskMetaData

	for rows.Next(){
		var task types.TaskMetaData
		var completedInt int
		err:= rows.Scan(&task.Title,&task.Description,&completedInt,&task.CreatedAt,&task.UpdatedAt)
		if err!= nil{
			return nil,err
		}
		task.Completed = completedInt == 1
		tasks = append(tasks, task)
	}
	return tasks,nil
}


func (s *Sqlite) MarkComplete(userid int64, taskid int64) error {
	stmt, err := s.Db.Prepare("UPDATE todo SET completed = 1, updated_at = ? WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	result, err := stmt.Exec(time.Now(), taskid, userid)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d does not belong to user with id %d or does not exist", taskid, userid)
	}
	
	return nil
}
