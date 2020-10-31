package generate

type Interface interface {
	Commands() ([]Command, error)
	Files() ([]File, error)
}
