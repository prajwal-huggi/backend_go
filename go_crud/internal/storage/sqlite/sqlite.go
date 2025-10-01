package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prajwal-huggi/backend_go/internal/config"
	"github.com/prajwal-huggi/backend_go/internal/types"
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

func (s *Sqlite) GetStudentById(id int64) (types.Student, error){
	stmt, err:= s.Db.Prepare("SELECT id, name, email, age FROM students WHERE id=?")
	if err!= nil{
		return types.Student{}, fmt.Errorf("prepare statement failed: %w", err)
	}
	defer stmt.Close()

	var student types.Student

	err= stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err!= nil{
		if err== sql.ErrNoRows{
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}

		return types.Student{}, fmt.Errorf("query error: %w", err)
	}

	return student, nil
}