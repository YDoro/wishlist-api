//go:generate mockgen --build_flags=--mod=mod -destination=../../mock/domain/id_mock.go -package=mocks . IDGenerator

package domain

// TODO - move to anoter package?
type IDGenerator interface {
	Generate() (string, error)
}
