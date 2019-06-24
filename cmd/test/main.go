package main

type inner struct {
}

func (i *inner) Foo() {
	println("foo")
}

type outer interface {
	Foo()
}

type outerImpl struct {
	inner
}

func main() {
	var o outer
	o = &outerImpl{}
	o.Foo()
}
