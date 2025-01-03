package registry

type Integration interface {
	Authenticate() error
	PullImage(image string) error
}
