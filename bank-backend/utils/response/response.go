package response

import (
	"fmt"
	"net/http"
)

// response struct
type Response struct {
	Code    int         `json:"code,omitempty"`
	Status  string      `json:"status,omitempty"`
	Content interface{} `json:"content,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ListResponse struct {
	Body interface{} `json:"body"`
	Page interface{} `json:"pagination"`
}

func NewResponse(code int, content interface{}, msg, err string) *Response {
	res := &Response{
		Code:    code,
		Status:  http.StatusText(code),
		Content: content,
		Message: msg,
		Error:   err,
	}

	return res
}

func Respond(err error, result interface{}, code int) *Response {
	rC, msg := MappingError(err, code)
	resp := &Response{
		Content: result,
		Code:    rC,
		Message: msg,
	}

	return resp
}

func ListRepond(body interface{}, pagination interface{}) *ListResponse {
	resp := &ListResponse{
		Body: body,
		Page: pagination,
	}

	return resp
}

func MappingError(err error, resCode int) (int, string) {
	switch resCode {
	case http.StatusOK:
		return resCode, ""
	case http.StatusPaymentRequired:
		return resCode, "Payment Required"
	case http.StatusBadRequest:
		return resCode, fmt.Sprint("Bad Request: ", err.Error())
	case http.StatusConflict:
		return resCode, "Data Conflict"
	case http.StatusNotFound:
		return resCode, "Data Not Found"
	case http.StatusInsufficientStorage:
		return resCode, "Error Database"
	case http.StatusForbidden:
		return resCode, "Access Denied"
	case http.StatusMethodNotAllowed:
		return resCode, "Action Not Allowed"
	case http.StatusUnauthorized:
		return resCode, "User is not authorized"
	case http.StatusUnprocessableEntity:
		return resCode, "Data can't be processed"
	case http.StatusInternalServerError:
		return resCode, fmt.Sprint("Internal Server Error: ", err.Error())
	default:
		return resCode, fmt.Sprint("Unexpected Error: ", err.Error())
	}
}
