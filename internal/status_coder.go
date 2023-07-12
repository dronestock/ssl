package internal

type StatusCoder interface {
	Code() int
	Message() string
}
