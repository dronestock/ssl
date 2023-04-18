package main

type statusCoder interface {
	code() int
	message() string
}
