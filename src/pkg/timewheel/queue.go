package timewheel

type Queue[T any] interface {
	Enqueue(t T) bool
}
