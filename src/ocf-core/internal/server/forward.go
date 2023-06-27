package server

import (
	"fmt"
	"net/http"
)

func ErrorHandler(res http.ResponseWriter, req *http.Request, err error) {
	res.Write([]byte(fmt.Sprintf("ERROR: %s", err.Error())))
}
