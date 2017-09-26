package ringslice

import (
	"reflect"
	"testing"
)

// type RingSlice interface {
// 	Append(interface{})
// 	Put(int, interface{}) error
// 	Get(int) (interface{}, error)
// 	GetFirst() interface{}
// 	GetLast() interface{}
// 	ToSlice() []interface{}
// }

func TestNew(t *testing.T) {
	t.Run("Generic new", func(t *testing.T) {
		if r, _ := New(10); !reflect.DeepEqual(r, getEmpty10()) {
			t.Error()
		}
	})
	t.Run("Generic new len", func(t *testing.T) {
		if r, _ := New(1234); !reflect.DeepEqual(r, getN(1234)) {
			t.Error()
		}
	})
	t.Run("New negative ints", func(t *testing.T) {
		if _, err := New(-1); err != ErrNegative {
			t.Error()
		}
	})
}

func TestAppend(t *testing.T) {
	r := &genericRingSlice{
		data: []interface{}{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		head: 0,
		tail: 0,
		len:  10,
	}
	t.Run("Generic append", func(t *testing.T) {
		if r.data[0] != 0 {
			t.Error("initialiazing error")
		}
		r.Append(4)
		if r.data[0].(int) != 4 {
			t.Error(r.data)
		}
		r.Append(6)
		if r.data[0].(int) != 4 || r.data[1].(int) != 6 {
			t.Error(r)
		}
	})
	t.Run("Append past head", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			r.Append(10 + i)
		}
		if !reflect.DeepEqual(r.data, []interface{}{18, 19, 10, 11, 12, 13, 14, 15, 16, 17}) {
			t.Error()
		}
		if r.head != 2 {
			t.Error()
		}
		if r.tail != 1 {
			t.Error()
		}
	})
	t.Run("Append past head twice", func(t *testing.T) {
		for i := 1; i < 11; i++ {
			r.Append(10 + i)
		}
		if !reflect.DeepEqual(r.data, []interface{}{19, 20, 11, 12, 13, 14, 15, 16, 17, 18}) {
			t.Error()
		}
		if r.head != 2 {
			t.Error()
		}
		if r.tail != 1 {
			t.Error()
		}
	})
}

func TestPut(t *testing.T) {
	r := &genericRingSlice{
		data: []interface{}{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		head: 0,
		tail: 0,
		len:  10,
	}
	t.Run("Generic put", func(t *testing.T) {
		if r.tail != 0 {
			t.Error()
		}
		r.Put(0, 1)
		r.Put(5, 6)
		if r.data[0].(int) != 1 || r.data[5].(int) != 6 {
			t.Error(r.data)
		}
		if r.tail != 5 {
			t.Error()
		}
	})
	t.Run("Put past len/Neg", func(t *testing.T) {
		r := &genericRingSlice{
			data: []interface{}{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			head: 0,
			tail: 0,
			len:  10,
		}
		err := r.Put(1234, 1234)
		if err != ErrReference {
			t.Error()
		}
		err = r.Put(-1234, 1234)
		if err != ErrNegative {
			t.Error()
		}
	})
}

func TestGet(t *testing.T) {
	r := &genericRingSlice{
		data: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		head: 0,
		tail: 9,
		len:  10,
	}
	looped := &genericRingSlice{
		data: []interface{}{5, 6, 7, 8, 9, 0, 1, 2, 3, 4},
		head: 5,
		tail: 4,
		len:  10,
	}
	t.Run("Generic get", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if val, _ := r.Get(i); val.(int) != i {
				t.Error(i)
			}
		}
	})
	t.Run("Get past len/neg", func(t *testing.T) {
		if _, err := r.Get(999999); err != ErrReference {
			t.Error()
		}
		if _, err := r.Get(-1); err != ErrNegative {
			t.Error()
		}
	})
	t.Run("Get with loop", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			if val, _ := looped.Get(i); val.(int) != i {
				t.Error(i)
			}
		}
	})
}

func TestGetFirst(t *testing.T) {
	r := &genericRingSlice{
		data: []interface{}{999, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		head: 0,
		tail: 9,
		len:  10,
	}
	looped := &genericRingSlice{
		data: []interface{}{5, 6, 7, 8, 9, 1234, 1, 2, 3, 4},
		head: 5,
		tail: 4,
		len:  10,
	}
	new, _ := New(10)
	t.Run("Generic GetFirst", func(t *testing.T) {
		if f, _ := r.GetFirst(); f != 999 {
			t.Error()
		}
	})
	t.Run("GetFirst with loop", func(t *testing.T) {
		if f, _ := looped.GetFirst(); f != 1234 {
			t.Error()
		}
	})
	t.Run("GetFirst empty", func(t *testing.T) {
		if _, err := new.GetFirst(); err != ErrReference {
			t.Error()
		}
	})
}

func TestGetLast(t *testing.T) {
	r := &genericRingSlice{
		data: []interface{}{999, 1, 2, 3, 4, 5, 6, 7, 8, 1234},
		head: 0,
		tail: 9,
		len:  10,
	}
	looped := &genericRingSlice{
		data: []interface{}{5, 6, 7, 8, 999, 1234, 1, 2, 3, 4},
		head: 5,
		tail: 4,
		len:  10,
	}
	new, _ := New(10)
	t.Run("Generic GetLast", func(t *testing.T) {
		if f, _ := r.GetLast(); f != 1234 {
			t.Error()
		}
	})
	t.Run("GetLast with loop", func(t *testing.T) {
		if f, _ := looped.GetLast(); f != 999 {
			t.Error()
		}
	})
	t.Run("GetLast empty", func(t *testing.T) {
		if _, err := new.GetLast(); err != ErrReference {
			t.Error()
		}
	})
}

func TestToSlice(t *testing.T) {
	empty := getEmpty()
	t.Run("ToSlice empty", func(t *testing.T) {
		testSlice(t, empty.ToSlice(), []interface{}{})
	})
	t.Run("ToSlice half full", func(t *testing.T) {
		testSlice(t, getHalf().ToSlice(), []interface{}{0, 1, 2, 3, 4, 5, 0, 0, 0, 0})
	})
	t.Run("ToSlice full", func(t *testing.T) {
		testSlice(t, getFull().ToSlice(), []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	})
	t.Run("ToSlice looped", func(t *testing.T) {
		testSlice(t, getLooped().ToSlice(), []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	})
}

func testSlice(t *testing.T, a []interface{}, b []interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Error(a, "\n\n", b)
	}
}

func getEmpty() RingSlice {
	return &genericRingSlice{
		data: make([]interface{}, 0),
		head: 0,
		tail: 0,
		len:  0,
	}
}

func getEmpty10() RingSlice {
	return &genericRingSlice{
		data: make([]interface{}, 10),
		head: 0,
		tail: 0,
		len:  10,
	}
}

func getHalf() RingSlice {
	return &genericRingSlice{
		data: []interface{}{0, 1, 2, 3, 4, 5, 0, 0, 0, 0},
		head: 0,
		tail: 5,
		len:  10,
	}
}

func getFull() RingSlice {
	return &genericRingSlice{
		data: []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		head: 0,
		tail: 9,
		len:  10,
	}
}

func getLooped() RingSlice {
	return &genericRingSlice{
		data: []interface{}{5, 6, 7, 8, 9, 0, 1, 2, 3, 4},
		head: 5,
		tail: 4,
		len:  10,
	}
}

func getN(n int) RingSlice {
	return &genericRingSlice{
		data: make([]interface{}, n),
		head: 0,
		tail: 0,
		len:  n,
	}
}
