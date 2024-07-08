package lingo

import (
	"sync"

	"github.com/vanh01/lingo/definition"
)

// FirstOrNil returns the first element of a sequence (with condition if any),
// or a nil value if no element is found
//
// predicate can be nil. If predicate is not empty or nil, we will use the first predicate
func (e Enumerable[T]) FirstOrNil(predicate ...Predicate[T]) T {
	var t T
	first := true
	for value := range e.getIter() {
		if !definition.IsEmptyOrNil(predicate) {
			if !predicate[0](value) {
				continue
			}
		}
		if first {
			first = false
			t = value
		}
	}
	return t
}

// FirstOrDefault returns the first element of a sequence (with condition if any),
// or a default value if no element is found
//
// predicate can be nil. If predicate is not empty or nil, we will use the first predicate
func (e Enumerable[T]) FirstOrDefault(defaultValue T, predicate ...Predicate[T]) T {
	var t T = defaultValue
	first := true
	for value := range e.getIter() {
		if !definition.IsEmptyOrNil(predicate) {
			if !predicate[0](value) {
				continue
			}
		}
		if first {
			first = false
			t = value
		}
	}
	return t
}

// LastOrNil returns the last element of a sequence (with condition if any),
// or a nil value if no element is found
//
// predicate can be nil. If predicate is not empty or nil, we will use the first predicate
func (e Enumerable[T]) LastOrNil(predicate ...Predicate[T]) T {
	var t T
	for value := range e.getIter() {
		if !definition.IsEmptyOrNil(predicate) {
			if !predicate[0](value) {
				continue
			}
		}
		t = value
	}
	return t
}

// LastOrDefault returns the last element of a sequence (with condition if any),
// or a default value if no element is found
//
// predicate can be nil. If predicate is not empty or nil, we will use the first predicate
func (e Enumerable[T]) LastOrDefault(defaultValue T, predicate ...Predicate[T]) T {
	var t T = defaultValue
	for value := range e.getIter() {
		if !definition.IsEmptyOrNil(predicate) {
			if !predicate[0](value) {
				continue
			}
		}
		t = value
	}
	return t
}

// ElementAtOrNil returns the element at a specified index in a sequence or a default value if the index is out of range.
func (e Enumerable[T]) ElementAtOrNil(index int64) T {
	var t T
	var i int64 = 0
	for value := range e.getIter() {
		if i <= index {
			t = value
			i++
		}
	}
	if i != index+1 {
		var tt T
		return tt
	}
	return t
}

// ElementAtOrDefault returns the element at a specified index in a sequence or a default value if the index is out of range.
func (e Enumerable[T]) ElementAtOrDefault(index int64, defaultValue T) T {
	var t T = defaultValue
	var i int64 = 0
	for value := range e.getIter() {
		if i <= index {
			t = value
			i++
		}
	}
	if i != index+1 {
		return defaultValue
	}
	return t
}

// ParallelEnumerable

// FirstOrNil returns the first element of a parallel sequence (with condition if any),
// or a nil value if no element is found
//
// predicate can be nil. If predicate is not empty or nil, we will use the first predicate
func (p ParallelEnumerable[T]) FirstOrNil(predicate ...Predicate[T]) T {
	if p.ordered {
		return p.AsEnumerable().FirstOrNil(predicate...)
	}

	var t T
	first := make(chan bool, 1)
	first <- true
	var wg sync.WaitGroup
	for value := range p.getIter() {
		wg.Add(1)
		temp := value
		go func() {
			tempFirst := <-first
			defer func() {
				first <- tempFirst
				wg.Done()
			}()
			if !definition.IsEmptyOrNil(predicate) {
				if !predicate[0](temp.val) {
					return
				}
			}
			if tempFirst {
				tempFirst = false
				t = temp.val
			}
		}()
	}
	wg.Wait()

	return t
}

// FirstOrDefault returns the first element of a parallel sequence (with condition if any),
// or a default value if no element is found
//
// predicate can be nil. If predicate is not empty or nil, we will use the first predicate
func (p ParallelEnumerable[T]) FirstOrDefault(defaultValue T, predicate ...Predicate[T]) T {
	if p.ordered {
		return p.AsEnumerable().FirstOrDefault(defaultValue, predicate...)
	}

	var t T = defaultValue
	first := make(chan bool, 1)
	first <- true
	var wg sync.WaitGroup
	for value := range p.getIter() {
		wg.Add(1)
		temp := value
		go func() {
			tempFirst := <-first
			defer func() {
				first <- tempFirst
				wg.Done()
			}()
			if !definition.IsEmptyOrNil(predicate) {
				if !predicate[0](temp.val) {
					return
				}
			}
			if tempFirst {
				tempFirst = false
				t = temp.val
			}
		}()
	}
	wg.Wait()

	return t
}

// LastOrNil returns the last element of a parallel sequence (with condition if any),
// or a nil value if no element is found
//
// predicate can be nil. If predicate is not empty or nil, we will use the first predicate
func (p ParallelEnumerable[T]) LastOrNil(predicate ...Predicate[T]) T {
	if p.ordered {
		return p.AsEnumerable().LastOrNil(predicate...)
	}

	tc := make(chan T)
	go func() {
		defer close(tc)
		var wg sync.WaitGroup
		for value := range p.getIter() {
			wg.Add(1)
			temp := value
			go func() {
				defer wg.Done()
				if !definition.IsEmptyOrNil(predicate) {
					if !predicate[0](temp.val) {
						return
					}
				}
				tc <- temp.val
			}()
		}
		wg.Wait()
	}()
	var t T

	for value := range tc {
		t = value
	}

	return t
}

// LastOrDefault returns the last element of a parallel sequence (with condition if any),
// or a default value if no element is found
//
// predicate can be nil. If predicate is not empty or nil, we will use the first predicate
func (p ParallelEnumerable[T]) LastOrDefault(defaultValue T, predicate ...Predicate[T]) T {
	if p.ordered {
		return p.AsEnumerable().LastOrDefault(defaultValue, predicate...)
	}

	tc := make(chan T)
	go func() {
		defer close(tc)
		var wg sync.WaitGroup
		for value := range p.getIter() {
			wg.Add(1)
			temp := value
			go func() {
				defer wg.Done()
				if !definition.IsEmptyOrNil(predicate) {
					if !predicate[0](temp.val) {
						return
					}
				}
				tc <- temp.val
			}()
		}
		wg.Wait()
	}()
	var t T = defaultValue

	for value := range tc {
		t = value
	}

	return t
}

// ElementAtOrNil returns the element at a specified index in a parallel sequence or a default value if the index is out of range.
func (p ParallelEnumerable[T]) ElementAtOrNil(index int64) T {
	if p.ordered {
		return p.AsEnumerable().ElementAtOrNil(index)
	}

	var t T
	i := make(chan int64, 1)
	i <- 0
	for value := range p.getIter() {
		temp := value
		go func() {
			tempi := <-i
			if tempi <= index {
				t = temp.val
				tempi++
			}
			i <- tempi
		}()
	}
	if <-i != index+1 {
		var tt T
		return tt
	}
	return t
}

// ElementAtOrDefault returns the element at a specified index in a parallel sequence or a default value if the index is out of range.
func (p ParallelEnumerable[T]) ElementAtOrDefault(index int64, defaultValue T) T {
	if p.ordered {
		return p.AsEnumerable().ElementAtOrDefault(index, defaultValue)
	}

	var t T = defaultValue
	i := make(chan int64, 1)
	i <- 0
	for value := range p.getIter() {
		temp := value
		go func() {
			tempi := <-i
			if tempi <= index {
				t = temp.val
				tempi++
			}
			i <- tempi
		}()
	}
	if <-i != index+1 {
		return defaultValue
	}
	return t
}
