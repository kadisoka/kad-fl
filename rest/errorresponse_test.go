package rest

import "testing"

func TestErrorResponseEmptyError(t *testing.T) {
	var err error = ErrorResponse{}
	assert(t, "<empty error>", err.Error())
}

func TestErrorResponseCodeError(t *testing.T) {
	var err error = ErrorResponse{Code: "MyCODE"}
	assert(t, "code MyCODE", err.Error())
}

func TestErrorResponseDescriptionError(t *testing.T) {
	var err error = ErrorResponse{Description: "This is an error"}
	assert(t, "This is an error", err.Error())
}

func TestErrorResponseFullError(t *testing.T) {
	var err error = ErrorResponse{Code: "MyCODE", Description: "This is an error"}
	assert(t, "[MyCODE] This is an error", err.Error())
}

func TestErrorResponseFieldEmptyError(t *testing.T) {
	var err error = ErrorResponseField{}
	assert(t, "<empty field error>", err.Error())
}

func TestErrorResponseFieldJustFieldError(t *testing.T) {
	var err error = ErrorResponseField{Field: "config.serve_port"}
	assert(t, "field config.serve_port", err.Error())
}

func TestErrorResponseFieldWithCodeError(t *testing.T) {
	var err error = ErrorResponseField{Field: "config.serve_port", Code: "INVALID"}
	assert(t, "config.serve_port: INVALID", err.Error())
}

func TestErrorResponseFieldWithDescError(t *testing.T) {
	var err error = ErrorResponseField{Field: "config.serve_port", Description: "Value must be between 100 and 30000"}
	assert(t, "config.serve_port: Value must be between 100 and 30000", err.Error())
}

func TestErrorResponseFieldFullError(t *testing.T) {
	var err error = ErrorResponseField{Field: "config.serve_port", Code: "INVALID", Description: "Value must be between 100 and 30000"}
	assert(t, "config.serve_port: [INVALID] Value must be between 100 and 30000", err.Error())
}

func TestErrorResponseFieldCodeError(t *testing.T) {
	var err error = ErrorResponseField{Code: "INVALID"}
	assert(t, "code INVALID", err.Error())
}

func TestErrorResponseFieldDescError(t *testing.T) {
	var err error = ErrorResponseField{Description: "Value must be between 100 and 30000"}
	assert(t, "Value must be between 100 and 30000", err.Error())
}

func TestErrorResponseFieldCodeDescError(t *testing.T) {
	var err error = ErrorResponseField{Code: "INVALID", Description: "Value must be between 100 and 30000"}
	assert(t, "[INVALID] Value must be between 100 and 30000", err.Error())
}
