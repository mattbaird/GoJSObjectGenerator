package jsify

import ()

type Foo struct {
	Id   string
	Test int
}

type Bar struct {
	Id         string
	Thing      Baz
	SuperThing int64
}

type Baz struct {
	X int
	Y int
}
