package reporting

type TypedStruct[abcd any] struct { // want `Type parameter name contains 'abcd', consider rephrasing to something else.`
	value abcd
}

type TypedStructWithInterface[t interface{ abcd() }] struct { // want `Method name contains 'abcd', consider rephrasing to something else.`
	value t
}

func TypedFunc[abcd any]() { // want `Type parameter name contains 'abcd', consider rephrasing to something else.`
}
