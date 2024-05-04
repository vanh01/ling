package lingo

import "reflect"

// Append appends a value to the end of the sequence.
func (e Enumerable[T]) Append(t T) Enumerable[T] {
	return Enumerable[T]{
		getIter: func() <-chan T {
			out := make(chan T)

			go func() {
				defer close(out)
				for value := range e.getIter() {
					out <- value
				}
				out <- t
			}()

			return out
		},
	}
}

// AppendRange appends the elements of the specified collection to the end of the sequence.
func (e Enumerable[T]) AppendRange(second Enumerable[T]) Enumerable[T] {
	return Enumerable[T]{
		getIter: func() <-chan T {
			out := make(chan T)

			go func() {
				defer close(out)
				for value := range e.getIter() {
					out <- value
				}
				for value := range second.getIter() {
					out <- value
				}
			}()

			return out
		},
	}
}

// Prepend adds a value to the beginning of the sequence.
func (e Enumerable[T]) Prepend(t T) Enumerable[T] {
	return Enumerable[T]{
		getIter: func() <-chan T {
			out := make(chan T)

			go func() {
				defer close(out)
				out <- t
				for value := range e.getIter() {
					out <- value
				}
			}()

			return out
		},
	}
}

// PrependRange appends the elements of the specified collection to the beginning of the sequence.
func (e Enumerable[T]) PrependRange(second Enumerable[T]) Enumerable[T] {
	return Enumerable[T]{
		getIter: func() <-chan T {
			out := make(chan T)

			go func() {
				defer close(out)
				for value := range second.getIter() {
					out <- value
				}
				for value := range e.getIter() {
					out <- value
				}
			}()

			return out
		},
	}
}

// Clear removes all elements of the sequence.
func (e Enumerable[T]) Clear() Enumerable[T] {
	return Empty[T]()
}

// Insert inserts an element into the sequence at the specified index.
func (e Enumerable[T]) Insert(index int, t T) Enumerable[T] {
	return Enumerable[T]{
		getIter: func() <-chan T {
			out := make(chan T)

			go func() {
				defer close(out)
				i := 0
				for value := range e.getIter() {
					if i == index {
						out <- t
					}
					out <- value
					i++
				}
				if i == index {
					out <- t
				}
			}()

			return out
		},
	}
}

// Remove removes the first occurrence of the given element, if found.
func (e Enumerable[T]) Remove(t T, comparer Comparer[T]) Enumerable[T] {
	return Enumerable[T]{
		getIter: func() <-chan T {
			out := make(chan T)

			go func() {
				defer close(out)
				isFirst := true
				for value := range e.getIter() {
					if isFirst {
						if comparer != nil && comparer(value, t) {
							isFirst = false
							continue
						}
						if comparer == nil && reflect.ValueOf(value).Interface() == reflect.ValueOf(t).Interface() {
							isFirst = false
							continue
						}
					}
					out <- value
				}
			}()

			return out
		},
	}
}

// RemoveAt removes the element at the specified index of the sequence.
func (e Enumerable[T]) RemoveAt(index int) Enumerable[T] {
	return Enumerable[T]{
		getIter: func() <-chan T {
			out := make(chan T)

			go func() {
				defer close(out)
				i := 0
				for value := range e.getIter() {
					if i == index {
						i++
						continue
					}
					out <- value
					i++
				}
			}()

			return out
		},
	}
}

// RemoveRange removes a range of elements from the sequence.
func (e Enumerable[T]) RemoveRange(index int, count int) Enumerable[T] {
	return Enumerable[T]{
		getIter: func() <-chan T {
			out := make(chan T)

			go func() {
				defer close(out)
				i := 0
				for value := range e.getIter() {
					if i >= index && count > 0 {
						count--
						continue
					}
					out <- value
					i++
				}
			}()

			return out
		},
	}
}
