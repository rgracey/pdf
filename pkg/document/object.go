package document

type Type int

const (
	UNKNOWN Type = iota
	BOOLEAN
	NUMBER_FLOAT
	NUMBER_INTEGER
	DICT
	STRING
	FUNCTION
	NAME
	ARRAY
	STREAM
	XREFS
	INDIRECT_OBJECT
	OBJECT_REF
)

type ObjectRef struct {
	Id         int
	Generation int
}

type Object struct {
	Type Type
	Ref  struct {
		Id         int
		Generation int
	}
	Header interface{}
	Data   interface{}
}
