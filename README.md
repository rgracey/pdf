# `pdf`
A simple PDF parser/serialiser for Go.

## Getting Started
### Example
#### Parsing
##### Parsing a file
```go
fileName := "sample.pdf"

ast, _ := pdf.ParseFile(fileName)
```

##### Parse a stream
```go
f, err := os.Open(fileName)

if err != nil {
    fmt.Println(err)
    return
}

defer f.Close()

ast, _ := pdf.ParseStream(f)
```

#### Serialising
```go
ast, _ := pdf.ParseFile("sample.pdf")

// ... manipulate the AST

serialised, _ := pdf.Serialise(ast)

f, err := os.Create("serialised.pdf")

if err != nil {
    fmt.Println(err)
    return
}

defer f.Close()

f.Write([]byte(serialised))
```

#### Finding a node
Finding a node by its ID
```go
type Filter func(n ast.PdfNode) bool

// A recursive function to find a node based on some criteria
func findNode(filter Filter, node ast.PdfNode) ast.PdfNode {
	if filter(node) {
		return node
	}

	for _, child := range node.Children() {
		if found := findNode(filter, child); found != nil {
			return found
		}
	}

	return nil
}

// Find an indirect object with ID 27.
node := findNode(func(n ast.PdfNode) bool {
    if n.Type() == ast.INDIRECT_OBJECT {
        if n.(*ast.IndirectObjectNode).Id() == 27 {
            return true
        }
    }

    return false
}, root)


// ... do something with node
```
