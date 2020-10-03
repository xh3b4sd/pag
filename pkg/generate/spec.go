package generate

type Interface interface {
	Generate() ([]Command, error)
}
