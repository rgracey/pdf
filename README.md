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
