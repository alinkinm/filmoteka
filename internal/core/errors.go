package core

import (
	"fmt"
	"net/http"
)

type Info struct {
	Msg        string
	StatusCode int
}

type MyError struct {
	Inf  Info
	Type string
}

func (err *MyError) Error() string {
	return fmt.Sprintf("message:'%s'", err.Inf.Msg)
}

func NewErrActorAlreadyExists() *MyError {
	return &MyError{Type: "ErrActorAlreadyExists", Inf: Info{Msg: "actor already exists in database", StatusCode: http.StatusBadRequest}}
}

func NewErrActorDoesNotExist() *MyError {
	return &MyError{Type: "ErrActorDoesNotExist", Inf: Info{Msg: "actor with this id does not exists", StatusCode: http.StatusBadRequest}}
}

func NewErrMovieAlreadyExists() *MyError {
	return &MyError{Type: "ErrActorAlreadyExists", Inf: Info{Msg: "movie already exists in database", StatusCode: http.StatusBadRequest}}
}

func NewErrMovieDoesNotExist() *MyError {
	return &MyError{Type: "ErrMovieDoesNotExist", Inf: Info{Msg: "movie with this id does not exists", StatusCode: http.StatusBadRequest}}
}
