package main

type TestInterfaceOne interface {
	Func1()
	Func2()
}

type TestInterfaceTwo interface {
	Func2()
	Func3()
}

type TestType2 struct{}

func (t *TestType2) Func1() {

}

func (t *TestType2) Func2() {

}

func (t *TestType2) Func3() {

}

func main() {
	ty := TestType2{}
	i := TestInterfaceOne(&ty)
	print(i)
}
