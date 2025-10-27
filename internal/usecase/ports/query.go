package ports

type Query[T any] interface {
	Get(id string) T
}
