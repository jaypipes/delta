package types

type MyInterfaceType interface {
	F()
}

var (
	NilInterface    MyInterfaceType = (*MyStructType)(nil)
	NotNilInterface MyInterfaceType = &MyStructType{}
)
