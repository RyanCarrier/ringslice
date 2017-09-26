package ringslice

import "errors"

//ErrReference is and error thrown when trying to reference outside the size
// of the RingSlice
var ErrReference = errors.New("index out side of slice length")

//ErrNegative is an error thrown when trying to create a ring or reference an
// index with a negative
var ErrNegative = errors.New("index/length cannot be negative")

//RingSlice is an implementation of looped array, can be used as an indexed
// linked list
type RingSlice interface {
	Append(interface{})
	Put(int, interface{}) error
	Get(int) (interface{}, error)
	GetFirst() (interface{}, error)
	GetLast() (interface{}, error)
	ToSlice() []interface{}
}

type genericRingSlice struct {
	data        []interface{}
	head        int
	tail        int
	len         int
	initialized bool
}

//New creates a new RingSlice of size n
func New(n int) (RingSlice, error) {
	if n <= 0 {
		return nil, ErrNegative
	}
	return &genericRingSlice{make([]interface{}, n), 0, 0, n, false}, nil
}

func (g *genericRingSlice) Append(val interface{}) {
	if g.tail+1 >= g.len {
		g.tail = 0
	} else if g.initialized {
		g.tail++
	}
	if g.tail == g.head && g.initialized {
		if g.head+1 >= g.len {
			g.head = 0
		} else {
			g.head++
		}
	} else if !g.initialized {
		g.initialized = true
	}
	g.data[g.tail] = val
}
func (g *genericRingSlice) Put(i int, val interface{}) error {
	if i >= g.len {
		return ErrReference
	}
	if i < 0 {
		return ErrNegative
	}
	abs := g.absoluteIndex(i)
	if abs > g.head && abs > g.tail {
		//don't need to worry about the other way, as in all regular circumstances;
		// head will start at 0, if head is not at 0 then more then len items
		// have been inserted (tail has moved)
		g.tail = abs
	}
	g.data[abs] = val
	return nil
}
func (g *genericRingSlice) Get(i int) (interface{}, error) {
	if i < 0 {
		return nil, ErrNegative
	}
	if i >= g.len || g.data[g.absoluteIndex(i)] == nil {
		return nil, ErrReference
	}
	return g.data[g.absoluteIndex(i)], nil
}
func (g *genericRingSlice) GetFirst() (interface{}, error) {
	if g.data[g.head] == nil {
		return nil, ErrReference
	}
	return g.data[g.head], nil
}
func (g *genericRingSlice) GetLast() (interface{}, error) {
	if g.data[g.tail] == nil {
		return nil, ErrReference
	}
	return g.data[g.tail], nil
}
func (g *genericRingSlice) ToSlice() []interface{} {
	if g.tail < g.head {
		return append(g.data[g.head:], g.data[:g.tail+1]...)
	}
	return g.data
}

func (g *genericRingSlice) absoluteIndex(i int) int {
	return (g.head + i) % g.len
}
