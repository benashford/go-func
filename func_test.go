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

func add(a, b int) int {
	return a + b
}

func TestReduce(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := Reduce(a, add).(int)
	if b != 10 {
		t.Error("Expecting 10, got:", b)
	}
}

func appendSlice(l []int, v int) []int {
	return append(l, v)
}

func TestReduce2(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := Filter(a, isEven)
	c := Reduce(b, appendSlice).([]int)
	if len(c) != 2 {
		t.Error("Should have length 2, has: ", len(c))
	}
	if c[0] != 2 && c[1] != 4 {
		t.Error("Results not [2, 4]:", c)
	}
}

func TestCombine(t *testing.T) {
	a := []int{1, 2, 3, 4, 5, 6}
	b := PMap(a, func(x int) int { return x * 3})
	c := Filter(b, isEven)
	d := Reduce(c, add).(int)

	if d != 36 {
		t.Error("Expecting 36 was:", d)
	}
}

func TestGroupBy(t *testing.T) {
	a := []string{"a", "bee", "sea", "aitch", "doubleyou", "zed"}
	b := GroupBy(a, stringLength).(map[int][]string)
	if len(b) != 4 {
		t.Error("Should have length 4, has: ", len(b), " b: ", b)
	}
	if len(b[3]) != 3 {
		t.Error("Should have three items of length 3, has: ", b[3])
	}
	if b[9][0] != "doubleyou" {
		t.Error("Expecting 'doubleyou' in position 9: ", b)
	}
}

type IndexableStruct struct {
	id int
	data string
}

func TestIndexBy(t *testing.T) {
	metadata := map[string]string{"ABC": "A12", "XYZ": "911", "BBB": "12X"}
	a := []*IndexableStruct{&IndexableStruct{1, "ABC"}, &IndexableStruct{2, "BBB"}, &IndexableStruct{3, "XYZ"}}
	b := IndexBy(a, func(s *IndexableStruct) string { return metadata[s.data] }).(map[string]*IndexableStruct)
	if len(b) != 3 {
		t.Error("Expecting on three items, actually: ", b)
	}
	if b["911"].data != "XYZ" {
		t.Error("Not expected: ", b)
	}
	if b["12X"].id != 2 {
		t.Error("Not expected: ", b)
	}
}

func fib(idx int) int {
	if idx == 0 {
		return 1
	}	else if idx == 1 {
		return 1
	} else {
		return fib(idx - 2) + fib(idx - 1)
	}
}

const testDataSize = 20

func sourceData() (source []int) {
	source = make([]int, testDataSize)
	for i := 0; i < testDataSize; i++ {
		source[i] = 20
	}
	return
}

func BenchmarkNone(b *testing.B) {
	source := sourceData()
	for i := 0; i < b.N; i++ {
		dest := make([]int, testDataSize)
		for idx, val := range source {
			dest[idx] = fib(val)
		}
	}
}

func BenchmarkMap(b *testing.B) {
	source := sourceData()
	for i := 0; i < b.N; i++ {
		fibs := Maps(source, fib).([]int)
		b.Log(len(fibs))
	}
}

func BenchmarkPMap(b *testing.B) {
	source := sourceData()
	for i := 0; i < b.N; i++ {
		fibs := PMaps(source, fib).([]int)
		b.Log(len(fibs))
	}
}
