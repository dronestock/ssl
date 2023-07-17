package core

type Domain struct {
	Id   string
	Name string
	Type DomainType
}

func (d *Domain) TencentType() *string {
	typ := ""
	switch d.Type {
	case DomainTypeCdn:
		typ = "cdn"
	case DomainTypeGateway:
		typ = "apigateway"
	default:
		typ = "unknown"
	}

	return &typ
}
