package ffi

var (
	TypeVoid              = Type{1, 1, 0, nil}
	TypeUint8             = Type{1, 1, 5, nil}
	TypeSint8             = Type{1, 1, 6, nil}
	TypeUint16            = Type{2, 2, 7, nil}
	TypeSint16            = Type{2, 2, 8, nil}
	TypeUint32            = Type{4, 4, 9, nil}
	TypeSint32            = Type{4, 4, 10, nil}
	TypeUint64            = Type{8, 8, 11, nil}
	TypeSint64            = Type{8, 8, 12, nil}
	TypeFloat             = Type{4, 4, 2, nil}
	TypeDouble            = Type{8, 8, 3, nil}
	TypePointer           = Type{8, 8, 14, nil}
	TypeLongdouble        = Type{16, 16, 4, nil}
	TypeComplexFloat      = Type{8, 4, 15, &[]*Type{&TypeFloat, nil}[0]}
	TypeComplexDouble     = Type{16, 8, 15, &[]*Type{&TypeDouble, nil}[0]}
	TypeComplexLongdouble = Type{32, 16, 15, &[]*Type{&TypeLongdouble, nil}[0]}
)
