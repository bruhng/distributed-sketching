package sketches

type Sketch[T any, R any] interface {
	Add(v T)
	Merge(sketch Sketch)
	Query(v T) R
}
