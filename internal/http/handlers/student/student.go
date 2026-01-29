package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mahimtalukder/go-student-api/internal/storage"
	"github.com/mahimtalukder/go-student-api/internal/types"
	"github.com/mahimtalukder/go-student-api/internal/utils/response"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		//Decode the body
		decodeError := json.NewDecoder(r.Body).Decode(&student)
		//Error Handel
		if errors.Is(decodeError, io.EOF) {
			err := response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("empty body")))
			if err != nil {
				return
			}
			return
		}
		if decodeError != nil {
			err := response.WriteJson(w, http.StatusBadRequest, response.GeneralError(decodeError))
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		//Request Validation

		if err := validator.New().Struct(&student); err != nil {
			var validateError validator.ValidationErrors
			errors.As(err, &validateError)
			err := response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateError))
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			return
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			err := response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			if err != nil {
				return
			}
		}

		slog.Info("Student created. ID: ", lastId)

		err = response.WriteJson(w, http.StatusOK, map[string]int64{"id": lastId})
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
}
