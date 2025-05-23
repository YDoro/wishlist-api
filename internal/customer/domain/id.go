package domain

// TODO - move to anoter package?
type IDGenerator interface {
	Generate() (string, error)
}
