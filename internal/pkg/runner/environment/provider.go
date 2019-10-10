package environment

type Provider interface {
	GetSecret(name string) (string, error)
	GetValue(name string) (string, error)
}

type Factory interface {
	Create(searchPath []string) Provider
}
