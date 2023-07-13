package core

type StatusCoder interface {
	Code() int
	Message() string
}
