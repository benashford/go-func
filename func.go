package funcs

import (
	"reflect"
)

const (
	defaultCapacity = 10
)

func SliceToChan(dataSlice interface{}) (ch interface{}) {
	sliceType := reflect.TypeOf(dataSlice)
	elementType := sliceType.Elem()
	channel := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, elementType), 0)
	ch = channel.Interface()

	go func() {
		dataSliceValue := reflect.ValueOf(dataSlice)
		dataSliceLen   := dataSliceValue.Len()
		for i := 0; i < dataSliceLen; i++ {
			dsv := dataSliceValue.Index(i)
			channel.Send(dsv)
		}
		channel.Close()
	}()

	return
}

func ChanToSlice(ch interface{}) interface{} {
	chType := reflect.TypeOf(ch)
	resultType := chType.Elem()
	result := reflect.MakeSlice(reflect.SliceOf(resultType), 0, defaultCapacity)
	chValue := reflect.ValueOf(ch)
	value, ok := chValue.Recv();
	for ok {
		result = reflect.Append(result, value)
		value, ok = chValue.Recv();
	}
	return result.Interface()
}

func MapChan(dataChan interface{}, f interface{}) (resultChan interface{}) {
	fType := reflect.TypeOf(f)
	fRetType := fType.Out(0)
	resultValue := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, fRetType), 0)
	resultChan = resultValue.Interface()

	go func() {
		fVal := reflect.ValueOf(f)
		chanValue := reflect.ValueOf(dataChan)
		value, ok := chanValue.Recv();
		for ok {
			results := fVal.Call([]reflect.Value{value})
			resultValue.Send(results[0])
			value, ok = chanValue.Recv();
		}
		resultValue.Close();
	}()

	return
}

func Maps(dataSlice interface{}, mapFunc interface{}) (resultSlice interface{}) {
	out := MapChan(SliceToChan(dataSlice), mapFunc)
	resultSlice = ChanToSlice(out)
	return
}

func Map(dataSlice interface{}, mapFunc interface{}) (resultChan interface{}) {
	return MapChan(SliceToChan(dataSlice), mapFunc)
}
