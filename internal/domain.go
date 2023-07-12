package internal

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
	default:
		typ = "unknown"
	}

	return &typ
}
