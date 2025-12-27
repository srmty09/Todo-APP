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

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.Storage_path)
	if err != nil {
		return nil, err
	}
	
	// Create user table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS user(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		return nil, err 
	}
	
	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}
	
	// Create todo table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS todo(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	description TEXT NOT NULL,
	completed BOOL DEFAULT FALSE,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
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
	stmt, err := s.Db.Prepare("UPDATE todo SET completed = TRUE, updated_at = ? WHERE id = ? AND user_id = ?")
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

func (s *Sqlite) MarkIncomplete(userid int64, taskid int64) error {
	stmt, err := s.Db.Prepare("UPDATE todo SET completed = FALSE, updated_at = ? WHERE id = ? AND user_id = ?")
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


func (s *Sqlite) DeletingTask(userid int64, taskid int64) error {
	stmt, err := s.Db.Prepare("DELETE FROM todo WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	result, err := stmt.Exec(taskid, userid)
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

func (s *Sqlite) GetSingleTask(userid int64, taskid int64) (*types.TaskMetaData, error) {
	stmt, err := s.Db.Prepare("SELECT title, description, completed, created_at, updated_at FROM todo WHERE id = ? AND user_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	
	var task types.TaskMetaData
	var completedInt int
	err = stmt.QueryRow(taskid, userid).Scan(&task.Title, &task.Description, &completedInt, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task with id %d does not belong to user with id %d or does not exist", taskid, userid)
	}
	if err != nil {
		return nil, err
	}
	
	task.Completed = completedInt == 1
	return &task, nil
}

func (s *Sqlite) EditTask(userid int64, taskid int64, title string, description string) error {
	stmt, err := s.Db.Prepare("UPDATE todo SET title = ?, description = ?, updated_at = ? WHERE id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	result, err := stmt.Exec(title, description, time.Now(), taskid, userid)
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


func (s *Sqlite) GetUser(userId int64)(*types.User,error){
	stmt,err := s.Db.Prepare("SELECT name, email FROM user WHERE id = ?")
	if err!= nil{
		return nil,err
	}
	defer stmt.Close()
	var user types.User
	err = stmt.QueryRow(userId).Scan(&user.Name,&user.Email)
	if err!=nil{
		return nil,err
	}
	return &user,nil
}

func (s *Sqlite) DeleteUser(userid int64)(error){
	// This single query deletes the user AND all their tasks automatically
	// because we have FOREIGN KEY (user_id) REFERENCES user(id) ON DELETE CASCADE
	stmt, err := s.Db.Prepare("DELETE FROM user WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	
	result, err := stmt.Exec(userid)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user with id %d does not exist", userid)
	}
	
	return nil
}


// Close closes the database connection
func (s *Sqlite) Close() error {
	return s.Db.Close()
}