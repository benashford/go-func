package funcs

import (
	"reflect"
)

const (
	defaultCapacity = 10
)

func SliceToChan(dataSlice interface{}) (ch chan interface{}) {
	ch = make(chan interface{})

	go func() {
		dataSliceValue := reflect.ValueOf(dataSlice)
		dataSliceLen   := dataSliceValue.Len()
		for i := 0; i < dataSliceLen; i++ {
			dsv := dataSliceValue.Index(i)
			ch <- dsv.Interface()
		}
		close(ch)
	}()

	return
}

func ChanToSlice(ch chan interface{}) interface{} {
	first := <- ch
	resultType := reflect.SliceOf(reflect.TypeOf(first))
	result := reflect.MakeSlice(resultType, 1, defaultCapacity)
	result.Index(0).Set(reflect.ValueOf(first))
	for val := range ch {
		result = reflect.Append(result, reflect.ValueOf(val))
	}
	return result.Interface()
}

func call(f interface{}, data interface{}) interface{} {
	fVal := reflect.ValueOf(f)
	dataVal := reflect.ValueOf(data)
	params := make([]reflect.Value, 1)
	params[0] = dataVal

	results := fVal.Call(params)

	return results[0].Interface();
}

func MapChan(data chan interface{}, f interface{}) (result chan interface{}) {
	result = make(chan interface{})

	go func() {
		for d := range data {
			result <- call(f, d)
		}
		close(result)
	}()

	return
}

func Map(dataSlice interface{}, mapFunc interface{}) (resultSlice interface{}) {
	out := MapChan(SliceToChan(dataSlice), mapFunc)
	resultSlice = ChanToSlice(out)
	return
}
