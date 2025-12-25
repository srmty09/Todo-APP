package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/srmty09/Todo-App/internal/config"
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
	userid INTEGER,
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



func (s *Sqlite) CreateUser(userid int,name string,email string)(int64, error){
	stmt,err:= s.Db.Prepare(
		"INSERT INTO user (userid,name,email) VALUES(?,?,?)")
	if err!=nil{
		return 0,err 
	}
	defer stmt.Close()

	res,err := stmt.Exec(userid,name,email)
	if err!=nil{
		return 0,err 
	}
	
	id,err:=res.LastInsertId()
	
	if err!=nil{
		return 0,err 
	}

	return id,nil
}

