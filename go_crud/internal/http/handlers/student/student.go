package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/prajwal-huggi/backend_go/internal/types"
	"github.com/prajwal-huggi/backend_go/internal/utils/response"
)

func New() http.HandlerFunc{
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

		response.WriteJson(w, http.StatusCreated, map[string]string{"success":"OK"})
	}
}
