package main

import (
	"encoding/json"
	"fmt"
)

type Test struct {
	C string
	D string
}

func main()  {
	Test2()
}

func Test1()  {
	//var test1 *Test

	//if test1 == nil {
	//	test1.A = "xxx"
	//}
}

func Test2()( a *Test)  {
	testStr:=`{
	"a": "a1",
		"b": "b1"
}`
	//fmt.Println( reflect.TypeOf(a) )
	//fmt.Printf("%t", a)

	err:=json.Unmarshal([]byte(testStr), &a)
	if err!=nil {
		fmt.Println(err, err.Error())
	}
	fmt.Printf("%+v", a)
	if a.C == "" {
		a.C= "c1"
	}
	fmt.Printf("%+v", a)
	//fmt.Println(a)
	return a
}