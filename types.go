package ffi

var (
	TypeVoid              = Type{1, 1, Void, nil}
	TypeUint8             = Type{1, 1, Uint8, nil}
	TypeSint8             = Type{1, 1, Sint8, nil}
	TypeUint16            = Type{2, 2, Uint16, nil}
	TypeSint16            = Type{2, 2, Sint16, nil}
	TypeUint32            = Type{4, 4, Uint32, nil}
	TypeSint32            = Type{4, 4, Sint32, nil}
	TypeUint64            = Type{8, 8, Uint64, nil}
	TypeSint64            = Type{8, 8, Sint64, nil}
	TypeFloat             = Type{4, 4, Float, nil}
	TypeDouble            = Type{8, 8, Double, nil}
	TypePointer           = Type{8, 8, Pointer, nil}
	TypeLongdouble        = Type{16, 16, Longdouble, nil}
	TypeComplexFloat      = Type{8, 4, Complex, &[]*Type{&TypeFloat, nil}[0]}
	TypeComplexDouble     = Type{16, 8, Complex, &[]*Type{&TypeDouble, nil}[0]}
	TypeComplexLongdouble = Type{32, 16, Complex, &[]*Type{&TypeLongdouble, nil}[0]}
)
