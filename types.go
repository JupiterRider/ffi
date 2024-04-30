package ffi

var (
	TypVoid              = Typ{1, 1, 0, nil}
	TypUint8             = Typ{1, 1, 5, nil}
	TypSint8             = Typ{1, 1, 6, nil}
	TypUint16            = Typ{2, 2, 7, nil}
	TypSint16            = Typ{2, 2, 8, nil}
	TypUint32            = Typ{4, 4, 9, nil}
	TypSint32            = Typ{4, 4, 10, nil}
	TypUint64            = Typ{8, 8, 11, nil}
	TypSint64            = Typ{8, 8, 12, nil}
	TypFloat             = Typ{4, 4, 2, nil}
	TypDouble            = Typ{8, 8, 3, nil}
	TypPointer           = Typ{8, 8, 14, nil}
	TypLongdouble        = Typ{16, 16, 4, nil}
	TypComplexFloat      = Typ{8, 4, 15, &[]*Typ{&TypFloat, nil}[0]}
	TypComplexDouble     = Typ{16, 8, 15, &[]*Typ{&TypDouble, nil}[0]}
	TypComplexLongdouble = Typ{32, 16, 15, &[]*Typ{&TypLongdouble, nil}[0]}
)
