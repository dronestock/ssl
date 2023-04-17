package main

type refresher interface {
	refresh(certificate *certificate)
}
