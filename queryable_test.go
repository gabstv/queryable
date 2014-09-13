package queryable

import (
	"encoding/json"
	"log"
	"testing"
)

func TestQueryable(t *testing.T) {

	q0 := New()

	jsn := []byte(`{"foo":1,"bar":"nope","pie":[3,18,34]}`)
	err := json.Unmarshal(jsn, &q0.Raw)
	if err != nil {
		t.Fatal(err)
	}
	v00, v01 := q0.QT("foo")
	if v00 != nil && v01 != nil {
		t.Fatal("ERROR!!!!")
	}
	log.Println(q0.Q("pie", 0).Raw)
	v000, err := q0.Q("pie", 0).IntT()
	if err != nil {
		t.Fatal("v000", err)
	}
	v001, _ := q0.Q("pie").Q(0).IntT()
	if v000 == v001 && v000 == 3 && v001 == 3 {
		log.Println("Awesome!")
	} else {
		t.Fatal("SHOULD BE EQUAL", v000, v001)
	}
	err = q0.Q("pie").Foreach(func(k, v *Queryable) bool {
		log.Println(k.Int(), v.Int())
		return true
	})
	if err != nil {
		t.Fatal(err)
	}
}
