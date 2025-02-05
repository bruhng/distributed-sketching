package sketches

type Sketch[T any, R any, X any] interface {
	Add(v T)
	Merge(sketch X)
	Query(v T) R
}
