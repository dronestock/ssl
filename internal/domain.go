package internal

type Domain struct {
	Id   string
	Name string
	Type DomainType
}

func (d *Domain) Tencent() (domain string) {
	switch d.Type {
	case DomainTypeCdn:

	}

	return
}
