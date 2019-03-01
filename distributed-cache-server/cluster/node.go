package cluster

type Node interface {
	ShouldProcess(key string) (string, bool)
	Members() []string
	Addr() string
}
