package queryable

import (
	"fmt"
	"reflect"
	"time"
)

type re struct {
	msg string
}

func (err *re) Error() string {
	return err.msg
}

func ne(err string) *re {
	return &re{err}
}

func fv(v interface{}) reflect.Value {
	v0 := reflect.ValueOf(v)
	for v0.Kind() == reflect.Ptr {
		v0 = v0.Elem()
	}
	return v0
}

// Creates a new Queryable element.
// You can pass the Raw pointer to unmarshallers.
func New() *Queryable {
	return &Queryable{nil}
	//return &Queryable{make(map[string]interface{})}
}

type Queryable struct {
	Raw interface{}
}

// Runs a query with one or more indexes
func (q *Queryable) QT(indexes ...interface{}) (*Queryable, error) {
	//fmt.Printf("QT( %v ) \n", indexes)
	rval := fv(q.Raw)
	kind := rval.Kind()
	if len(indexes) > 0 {
		switch kind {
		case reflect.Slice, reflect.Array, reflect.String, reflect.Struct, reflect.Map:
			//OK
		default:
			return nil, ne("non queryable element type: " + kind.String() + " " + fmt.Sprintf("%v %v", len(indexes), indexes[0]))
		}
	} else {
		return q, nil
	}
	switch kind {
	case reflect.String:
		if len(indexes) > 1 {
			return nil, ne("a rune is not queryable")
		}
		if reflect.ValueOf(indexes[0]).Kind() != reflect.Int {
			return nil, ne("rune index should be an int")
		}
		return &Queryable{rval.String()[int(reflect.ValueOf(indexes[0]).Int())]}, nil
	case reflect.Map:
		if rval.Type().Key().Kind() != reflect.ValueOf(indexes[0]).Kind() {
			return nil, ne("key kind differs index kind")
		}
		//TODO: see if it needs to use fv() below
		a := rval.MapIndex(reflect.ValueOf(indexes[0])).Interface()
		b := &Queryable{a}
		return b.QT(indexes[1:]...)
	case reflect.Struct:
		if fv(indexes[0]).Kind() != reflect.String {
			return nil, ne("key kind differs index kind")
		}
		a := rval.FieldByName(fv(indexes[0]).String())
		b := &Queryable{a}
		return b.QT(indexes[1:]...)
	case reflect.Slice, reflect.Array:
		//TODO: catch out of range panics
		if fv(indexes[0]).Kind() != reflect.Int {
			return nil, ne("key of slice/array is not int")
		}
		a := rval.Index(int(fv(indexes[0]).Int())).Interface()
		b := &Queryable{a}
		return b.QT(indexes[1:]...)
	}
	//TODO: support chan
	//
	return nil, nil
}

// Runs a query with one or more indexes. Omits errors.
func (q *Queryable) Q(indexes ...interface{}) *Queryable {
	qt, err := q.QT(indexes...)
	if err != nil {
		fmt.Println("FATAL Queryable Q", err)
		return nil
	}
	return qt
}

// Runs a query with one or more indexes. It panics if an error is found.
func (q *Queryable) QMust(indexes ...interface{}) *Queryable {
	qt, err := q.QT(indexes...)
	if err != nil {
		panic("FATAL Queryable Q: " + err.Error())
	}
	return qt
}

// Returns the length of the contained value. It returns 0 if the value is not an Array, Slice, Chan, String or Map.
func (q *Queryable) Len() int {
	//TODO: maybe support struct
	rval := fv(q.Raw)
	kind := rval.Kind()
	switch kind {
	case reflect.Array, reflect.Slice, reflect.Chan, reflect.String, reflect.Map:
		return rval.Len()
	}
	return 0
}

// Loops through all elements of a map, slice or array.
func (q *Queryable) Foreach(f func(key, val *Queryable) bool) error {
	//TODO: maybe support struct
	rval := fv(q.Raw)
	kind := rval.Kind()
	switch kind {
	case reflect.Array, reflect.Slice, reflect.Map:
		// OK
	default:
		return ne("not loopable (Foreach works on arrays, slices and maps)")
	}
	if kind == reflect.Map {
		keys := rval.MapKeys()
		for i := 0; i < len(keys); i++ {
			val := rval.MapIndex(keys[i])
			ok := f(&Queryable{keys[i].Interface()}, &Queryable{val.Interface()})
			if !ok {
				return nil
			}
		}
	}
	// slice
	for i := 0; i < rval.Len(); i++ {
		val := rval.Index(i)
		ok := f(&Queryable{i}, &Queryable{val.Interface()})
		if !ok {
			return nil
		}
	}
	return nil
}

func (q *Queryable) IntT() (int, error) {
	fvv := fv(q.Raw)
	switch fvv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(fvv.Int()), nil
	case reflect.Float64, reflect.Float32:
		return int(fvv.Float()), nil
	}
	return 0, ne("not an int " + fvv.Kind().String())
}

func (q *Queryable) Int64T() (int64, error) {
	fvv := fv(q.Raw)
	switch fvv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fvv.Int(), nil
	case reflect.Float64, reflect.Float32:
		return int64(fvv.Float()), nil
	}
	return 0, ne("not an int64 " + fvv.Kind().String())
}

func (q *Queryable) Float64T() (float64, error) {
	fvv := fv(q.Raw)
	switch fvv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(fvv.Int()), nil
	case reflect.Float64, reflect.Float32:
		return fvv.Float(), nil
	}
	return 0, ne("not a float64 " + fvv.Kind().String())
}

func (q *Queryable) BoolT() (bool, error) {
	fvv := fv(q.Raw)
	switch fvv.Kind() {
	case reflect.Bool:
		return fvv.Bool(), nil
	}
	return false, ne("not a bool " + fvv.Kind().String())
}

func (q *Queryable) StringT() (string, error) {
	fvv := fv(q.Raw)
	switch fvv.Kind() {
	case reflect.String:
		return fvv.String(), nil
	}
	return "", ne("not a string " + fvv.Kind().String())
}

func (q *Queryable) TimeT() (time.Time, error) {
	fvv := fv(q.Raw)
	switch fvv.Kind() {
	case reflect.String:
		t := time.Time{}
		err := t.UnmarshalText([]byte(fvv.String()))
		return t, err
	case reflect.Uint:
		return time.Unix(int64(fvv.Uint()), 0), nil
	case reflect.Int:
		return time.Unix(fvv.Int(), 0), nil
	case reflect.Struct:
		t, ok := q.Raw.(time.Time)
		if ok {
			return t, nil
		}
	}
	return time.Time{}, ne("not a valid time object " + fvv.Kind().String())
}

func (q *Queryable) Int() int {
	v, _ := q.IntT()
	return v
}

func (q *Queryable) Int64() int64 {
	v, _ := q.Int64T()
	return v
}

func (q *Queryable) Float64() float64 {
	v, _ := q.Float64T()
	return v
}

func (q *Queryable) Bool() bool {
	v, _ := q.BoolT()
	return v
}

func (q *Queryable) String() string {
	v, _ := q.StringT()
	return v
}

func (q *Queryable) Time() time.Time {
	v, _ := q.TimeT()
	return v
}
