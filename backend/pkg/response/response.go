package response

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// SuccessBody wraps a successful JSON response.
type SuccessBody struct {
	Data interface{} `json:"data"`
}

// ErrorBody wraps an error JSON response.
type ErrorBody struct {
	Error string      `json:"error"`
	Meta  interface{} `json:"meta,omitempty"`
}

// WriteJSON encodes v as JSON and writes it with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError writes a standard JSON error response.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ErrorBody{Error: message})
}

// WriteValidationError formats validator.ValidationErrors into a readable response.
func WriteValidationError(w http.ResponseWriter, err error) {
	var ve validator.ValidationErrors
	if ok := isValidationError(err, &ve); !ok {
		WriteError(w, http.StatusBadRequest, "validation failed")
		return
	}

	fields := make(map[string]string, len(ve))
	for _, fe := range ve {
		fields[fe.Field()] = fieldErrMsg(fe)
	}

	WriteJSON(w, http.StatusBadRequest, ErrorBody{
		Error: "validation failed",
		Meta:  fields,
	})
}

// DecodeJSON decodes the request body into v.
func DecodeJSON(r *http.Request, v interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

// ── helpers ──────────────────────────────────────────────────────────────────

func isValidationError(err error, ve *validator.ValidationErrors) bool {
	switch e := err.(type) {
	case validator.ValidationErrors:
		*ve = e
		return true
	}
	return false
}

func fieldErrMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "value is too short (min " + fe.Param() + ")"
	case "max":
		return "value is too long (max " + fe.Param() + ")"
	default:
		return "invalid value"
	}
}
