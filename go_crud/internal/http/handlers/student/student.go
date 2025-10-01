package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/prajwal-huggi/backend_go/internal/storage"
	"github.com/prajwal-huggi/backend_go/internal/types"
	"github.com/prajwal-huggi/backend_go/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("Creating the student")

		var student types.Student

		err:= json.NewDecoder(r.Body).Decode(&student)
		//If the json body is empty then in that case the below error will be triggered
		if errors.Is(err, io.EOF){
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			return 
		}

		if err!= nil{
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//Request validation
		if err:= validator.New().Struct(student); err!= nil{
			validateErrs:= err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))

			return
		}

		lastId, err:= storage.CreateStudent(student.Name, student.Email, student.Age,)

		slog.Info("User created successfully", slog.String("userId", fmt.Sprint(lastId)))
		if err!= nil{
			response.WriteJson(w, http.StatusInternalServerError, err)
			return  
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id":lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		id:= r.PathValue("id")
		slog.Info("Getting a student", slog.String("id", id))

		intId, err:= strconv.ParseInt(id,10,64)
		if err!= nil{
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		
		student, err:= storage.GetStudentById(intId)
		if err!= nil{
			slog.Error("Error getting the user", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		
		response.WriteJson(w, http.StatusOK, student)
	}
}