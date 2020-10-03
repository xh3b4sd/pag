package generate

type Interface interface {
	Generate() ([]Context, error)
}
