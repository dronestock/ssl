package core

const (
	DomainTypeCdn DomainType = iota + 1
	DomainTypeGateway
)

type DomainType uint8
