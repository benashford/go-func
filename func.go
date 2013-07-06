package funcs

import (
	"fmt"
	"reflect"
	"runtime"
)

const (
	defaultCapacity = 10
)

var (
	cpus = runtime.NumCPU()
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

func mapChan(dataChan interface{}, f interface{}) (resultChan interface{}) {
	fType := reflect.TypeOf(f)
	fRetType := fType.Out(0)
	resultValue := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, fRetType), 0)
	resultChan = resultValue.Interface()

	go func() {
		fVal := reflect.ValueOf(f)
		chanValue := reflect.ValueOf(dataChan)
		value, ok := chanValue.Recv()
		for ok {
			results := fVal.Call([]reflect.Value{value})
			resultValue.Send(results[0])
			value, ok = chanValue.Recv()
		}
		resultValue.Close()
	}()

	return
}

func Maps(data interface{}, mapFunc interface{}) (resultSlice interface{}) {
	out := Map(data, mapFunc)
	resultSlice = ChanToSlice(out)
	return
}

func Map(data interface{}, mapFunc interface{}) (resultChan interface{}) {
	dataType := reflect.TypeOf(data)
	switch dataType.Kind() {
	case reflect.Chan:
		return mapChan(data, mapFunc)
	case reflect.Slice:
		return mapChan(SliceToChan(data), mapFunc)
	default:
		panic(fmt.Sprintf("Unrecognised data: %s", data))
	}
}

func pMapChanInt(inChan reflect.Value, f interface{}, outChan reflect.Value) {
	fVal := reflect.ValueOf(f)
	val, ok := inChan.Recv()
	for ok {
		results := fVal.Call([]reflect.Value{val})
		outChan.Send(results[0])
		val, ok = inChan.Recv()
	}
	outChan.Close()
}

func pMapFeedInChans(dataChan interface{}, inChans []reflect.Value) {
	dataValue := reflect.ValueOf(dataChan)
	numChans := len(inChans)
	idx := 0
	val, ok := dataValue.Recv()
	for ok {
		inChans[idx % numChans].Send(val)
		idx++
		val, ok = dataValue.Recv()
	}
	for _, inChan := range inChans {
		inChan.Close()
	}
}

func pMapDrainOutChans(outChans []reflect.Value, resultChan reflect.Value) {
	numChans := len(outChans)
	idx := 0
	val, ok := outChans[idx].Recv()
	for ok {
		resultChan.Send(val)
		idx++
		val, ok = outChans[idx % numChans].Recv()
	}
	resultChan.Close()
}

func pMapChan(dataChan interface{}, f interface{}) (resultChan interface{}) {
	fType := reflect.TypeOf(f)
	fRetType := fType.Out(0)
	dataChanType := reflect.TypeOf(dataChan)
	dataChanElemType := dataChanType.Elem()
	inChans := make([]reflect.Value, cpus)
	outChans := make([]reflect.Value, cpus)
	for i := 0; i < cpus; i++ {
		inChans[i] = reflect.MakeChan(reflect.ChanOf(reflect.BothDir, dataChanElemType), cpus)
		outChans[i] = reflect.MakeChan(reflect.ChanOf(reflect.BothDir, fRetType), cpus)

		go pMapChanInt(inChans[i], f, outChans[i])
	}
	go pMapFeedInChans(dataChan, inChans)

	resultChanValue := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, fRetType), cpus)
	resultChan = resultChanValue.Interface()

	go pMapDrainOutChans(outChans, resultChanValue)

	return
}

func PMaps(data interface{}, mapFunc interface{}) (resultSlice interface{}) {
	out := PMap(data, mapFunc)
	resultSlice = ChanToSlice(out)
	return
}

func PMap(data interface{}, mapFunc interface{}) (resultChan interface{}) {
	dataType := reflect.TypeOf(data)
	switch dataType.Kind() {
	case reflect.Chan:
		return pMapChan(data, mapFunc)
	case reflect.Slice:
		return pMapChan(SliceToChan(data), mapFunc)
	default:
		panic(fmt.Sprintf("Unexpected data: %v", data))
	}
}

func filterChan(dataChan interface{}, f interface{}) (resultChan interface{}) {
	chanType := reflect.TypeOf(dataChan)
	elemType := chanType.Elem()
	resultValue := reflect.MakeChan(reflect.ChanOf(reflect.BothDir, elemType), 0)
	resultChan = resultValue.Interface()

	go func() {
		fVal := reflect.ValueOf(f)
		chanValue := reflect.ValueOf(dataChan)
		value, ok := chanValue.Recv()
		for ok {
			results := fVal.Call([]reflect.Value{value})
			pass := results[0].Interface().(bool)
			if pass {
				resultValue.Send(value)
			}
			value, ok = chanValue.Recv()
		}
		resultValue.Close()
	}()

	return
}

func Filters(data interface{}, f interface{}) (resultSlice interface{}) {
	out := Filter(data, f)
	resultSlice = ChanToSlice(out)
	return
}

func Filter(data interface{}, f interface{}) (resultChan interface{}) {
	dataType := reflect.TypeOf(data)
	switch dataType.Kind() {
	case reflect.Chan:
		resultChan = filterChan(data, f)
	case reflect.Slice:
		resultChan = filterChan(SliceToChan(data), f)
	default:
		panic(fmt.Sprintf("Unexpected data: %v", data))
	}
	return
}

func reduceChan(dataChan interface{}, f interface{}) (result interface{}) {
	fType := reflect.TypeOf(f)
	fRetType := fType.Out(0)
	val := reflect.Zero(fRetType)

	fVal := reflect.ValueOf(f)

	chanValue := reflect.ValueOf(dataChan)
	value, ok := chanValue.Recv()
	for ok {
		results := fVal.Call([]reflect.Value{val, value})
		val = results[0]
		value, ok = chanValue.Recv()
	}
	result = val.Interface()
	return
}

func Reduce(data interface{}, f interface{}) (result interface{}) {
	dataType := reflect.TypeOf(data)
	switch dataType.Kind() {
	case reflect.Chan:
		result = reduceChan(data, f)
	case reflect.Slice:
		result = reduceChan(SliceToChan(data), f)
	default:
		panic(fmt.Sprintf("Unexpected data: %v", data))
	}
	return
}
