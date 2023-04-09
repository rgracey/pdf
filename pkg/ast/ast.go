package ast

type Type int

const (
	ROOT Type = iota
	BOOLEAN
	FLOAT
	INTEGER
	NAME
	DICT
	DICT_ENTRY
	STRING
	FUNCTION
	ARRAY
	STREAM
	XREFS
	INDIRECT_OBJECT
	OBJECT_REF
)

type PdfNode interface {
	Type() Type                // Get the type of the node
	Value() interface{}        // Get the value of the node
	SetValue(interface{})      // Set the value of the node
	Children() []PdfNode       // Get the children of the node
	AddChild(PdfNode)          // Add a child to the node
	RemoveChild(int)           // Remove a child from the node
	ReplaceChild(int, PdfNode) // Replace a child with another node
	Clone() PdfNode            // Clone the node
}

type pdfNode struct {
	nodeType Type
	value    interface{}
	children []PdfNode
}

func (n *pdfNode) Type() Type {
	return n.nodeType
}

func (n *pdfNode) Value() interface{} {
	return n.value
}

func (n *pdfNode) SetValue(value interface{}) {
	switch n.nodeType {
	case BOOLEAN:
		if _, ok := value.(bool); !ok {
			panic("Value is not a boolean")
		}
	case FLOAT:
		if _, ok := value.(float64); !ok {
			panic("Value is not a float")
		}
	case INTEGER:
		if _, ok := value.(int); !ok {
			panic("Value is not an integer")
		}
	case STRING, NAME:
		if _, ok := value.(string); !ok {
			panic("Value is not a string")
		}
	}

	n.value = value
}

func (n *pdfNode) Children() []PdfNode {
	return n.children
}

func (n *pdfNode) AddChild(child PdfNode) {
	n.children = append(n.children, child)
}

func (n *pdfNode) RemoveChild(index int) {
	n.children = append(n.children[:index], n.children[index+1:]...)
}

func (n *pdfNode) ReplaceChild(index int, child PdfNode) {
	n.children[index] = child
}

func (n *pdfNode) Clone() PdfNode {
	clone := &pdfNode{
		nodeType: n.nodeType,
		value:    n.value,
	}

	for _, child := range n.children {
		clone.children = append(clone.children, child.Clone())
	}

	return clone
}

type RootNode struct {
	*pdfNode
}

func NewRootNode() *RootNode {
	return &RootNode{
		&pdfNode{
			nodeType: ROOT,
		},
	}
}

type BooleanNode struct {
	*pdfNode
}

func NewBooleanNode(value bool) *BooleanNode {
	return &BooleanNode{
		&pdfNode{
			nodeType: BOOLEAN,
			value:    value,
		},
	}
}

type FloatNode struct {
	*pdfNode
}

func NewFloatNode(value float64) *FloatNode {
	return &FloatNode{
		&pdfNode{
			nodeType: FLOAT,
			value:    value,
		},
	}
}

type IntegerNode struct {
	*pdfNode
}

func NewIntegerNode(value int64) *IntegerNode {
	return &IntegerNode{
		&pdfNode{
			nodeType: INTEGER,
			value:    value,
		},
	}
}

type NameNode struct {
	*pdfNode
}

func NewNameNode(value string) *NameNode {
	return &NameNode{
		&pdfNode{
			nodeType: NAME,
			value:    value,
		},
	}
}

type StringNode struct {
	*pdfNode
}

func NewStringNode(value string) *StringNode {
	return &StringNode{
		&pdfNode{
			nodeType: STRING,
			value:    value,
		},
	}
}

type FunctionNode struct {
	*pdfNode
}

func NewFunctionNode() *FunctionNode {
	return &FunctionNode{
		&pdfNode{
			nodeType: FUNCTION,
		},
	}
}

type ArrayNode struct {
	*pdfNode
}

func NewArrayNode() *ArrayNode {
	return &ArrayNode{
		&pdfNode{
			nodeType: ARRAY,
		},
	}
}

type StreamNode struct {
	*pdfNode
}

func NewStreamNode(value string) *StreamNode {
	return &StreamNode{
		&pdfNode{
			nodeType: STREAM,
			value:    value,
		},
	}
}

type XRefsNode struct {
	*pdfNode
}

func NewXRefsNode() *XRefsNode {
	return &XRefsNode{
		&pdfNode{
			nodeType: XREFS,
		},
	}
}

type IndirectObjectNode struct {
	*pdfNode
	id  int64
	gen int64
}

func NewIndirectObjectNode(id int64, gen int64) *IndirectObjectNode {
	return &IndirectObjectNode{
		&pdfNode{
			nodeType: INDIRECT_OBJECT,
		},
		id,
		gen,
	}
}

func (n *IndirectObjectNode) Id() int64 {
	return n.id
}

func (n *IndirectObjectNode) Gen() int64 {
	return n.gen
}

type ObjectRefNode struct {
	*pdfNode
	id  int64
	gen int64
}

func NewObjectRefNode(id int64, gen int64) *ObjectRefNode {
	return &ObjectRefNode{
		&pdfNode{
			nodeType: OBJECT_REF,
		},
		id,
		gen,
	}
}

func (n *ObjectRefNode) Id() int64 {
	return n.id
}

func (n *ObjectRefNode) Gen() int64 {
	return n.gen
}

type DictNode struct {
	*pdfNode
	entries map[string]PdfNode
}

func NewDictNode() *DictNode {
	return &DictNode{
		&pdfNode{
			nodeType: DICT,
		},
		make(map[string]PdfNode),
	}
}

func (n *DictNode) AddChild(child PdfNode) {
	if len(n.children)%2 == 0 {
		// If the child is not a name, then it's possible that we're dealing
		// with an object reference as the eventual value, and it'll be popped
		// off our children (becase the parsing of objects and references are
		// done retroactively after seeing the keyword [obj or R]). So just add
		// it to our children and return.
		if child.Type() != NAME {
			n.children = append(n.children, child)
			return
		}

		n.entries[child.Value().(string)] = nil
		n.children = append(n.children, child)
		return
	}

	n.children = append(n.children, child)
	n.entries[n.children[len(n.children)-2].Value().(string)] = child
}

func (n *DictNode) Get(key string) PdfNode {
	return n.entries[key]
}
