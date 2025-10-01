package sqlite

import (
	"database/sql"

	"github.com/prajwal-huggi/backend_go/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct{
	Db *sql.DB
}

func New(cfg *config.Config)(*Sqlite, error){
	db, err:= sql.Open("sqlite3", cfg.StoragePath)
	if err != nil{
		return nil, err
	}

	_, err= db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	age INTEGER,
	email TEXT
	)`)	

	if err!= nil{
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

//we are doing the implementation in the below function
func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error){
	stmt, err:= s.Db.Prepare("INSERT INTO students (name, age, email) VALUES (?,?,?)")
	if err!= nil{
		return 0, err
	}

	defer stmt.Close()

	res, err:= stmt.Exec(name, age, email)
	if err!= nil{
		return 0, err
	}

	lastId, err:= res.LastInsertId()
	if err!= nil{
		return 0, nil
	}
	
	return lastId, nil
}
