# Queryable
A helper package for querying structured data (json, xml) without having to create structs for them.

Usage:
```Go
package main

import(
	"fmt"
	"github.com/gabstv/queryable"
	"encoding/json"
)

func main() {
	data := []byte(`{"foo":1,"bar":"something","foobar":{"apples":[1,45,67],"author":"gabs"}}`)
	v := queryable.New()
	json.Unmarshal(data, &v.Raw)

	foo := v.Q("foo").Int()
	fmt.Printf("foo: %v\n", foo)

	apple2 := v.Q("foobar","apples", 1).Int() // or v.Q("foobar").Q("apples").Q(1).Int()
	fmt.Printf("apple2: %v\n", apple2)

	author := v.Q("foobar", "author").String()
	fmt.Printf("author: %v\n", author)

	// you can also loop through maps or slices
	v.Q("foobar", "apples").Foreach(func(key, val *queryable.Queryable) bool {
		fmt.Printf("key: %v val: %v\n", key.Int(), val.Int())
		return true // return false to stop the loop
	})
}
```