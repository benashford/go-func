package funcs

import (
	"testing"
)

func double(a int) int {
	return a * 2
}

func TestMapInt(t *testing.T) {
	a := []int{1, 2, 3}
	b := Map(a, double).(chan int)
	x := <- b
	y := <- b
	z := <- b
	if x != 2 || y != 4 || z != 6 {
		t.Errorf("does not equal [2, 4, 6]")
	}
}

func TestPMapInt(t *testing.T) {
	a := []int{1, 2, 3}
	b := PMap(a, double).(chan int)
	x := <- b
	y := <- b
	z := <- b
	if x != 2 || y != 4 || z != 6 {
		t.Errorf("does not equal [2, 4, 6]")
	}
}

func stringLength(s string) int {
	return len(s)
}

func TestMapString(t *testing.T) {
	a := []string{"a", "bee", "sea", "dee", "eff", "gee", "aitch"}
	b := Maps(a, stringLength).([]int)
	if len(b) != 7 {
		t.Error("Expecting seven values")
	}
	if b[0] != 1 || b[1] != 3 || b[6] != 5 {
		t.Errorf("%v does not equal [1, 3, 3, 3, 3, 3, 5]", b)
	}
}

func TestPMapString(t *testing.T) {
	a := []string{"a", "bee", "sea", "dee", "eff", "gee", "aitch"}
	b := PMaps(a, stringLength).([]int)
	if len(b) != 7 {
		t.Error("Expecting seven values")
	}
	if b[0] != 1 || b[1] != 3 || b[6] != 5 {
		t.Errorf("%v does not equal [1, 3, 3, 3, 3, 3, 5]", b)
	}
}

type TestStruct struct {
	a int
	b int
}

func doToTest(ts *TestStruct) (*TestStruct) {
	ts.b = ts.a + 1
	return ts
}

func TestMapStruct(t *testing.T) {
	ts1 := &TestStruct{100, 0}
	ts2 := &TestStruct{200, 0}
	a := []*TestStruct{ts1, ts2}
	b := Maps(a, doToTest).([]*TestStruct)
	if b[0].b != 101 || b[1].b != 201 {
		t.Errorf("%v does not equal the expected result")
	}
}

func TestMapChan(t *testing.T) {
	in := make(chan int)
	out := Map(in, double).(chan int)

	in <- 1
	result := <- out

	if result != 2 {
		t.Errorf("%v does not equal 2", result)
	}
}

func isEven(a int) bool {
	return a % 2 == 0
}

func TestFilters(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := Filters(a, isEven).([]int)
	if len(b) != 2 {
		t.Error("Expecting two results")
	}
	if b[0] != 2 || b[1] != 4 {
		t.Error("Unexpected result:", b)
	}
}
