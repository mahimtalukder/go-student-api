package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mahimtalukder/go-student-api/internal/types"
	"github.com/mahimtalukder/go-student-api/internal/utils/response"
)

func New() http.HandlerFunc {
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

		err := response.WriteJson(w, http.StatusOK, student)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}
}
