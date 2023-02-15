package main

type TestInterface interface {
	DoSomeWork()
}

type TestType struct {
	a int
}

func (t *TestType) DoSomeWork() {
}

func main() {
	t := &TestType{}
	i := TestInterface(t)
	print(i)
}
