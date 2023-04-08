package document

type Pdf struct {
	version string
	objects []Object
	xrefs   []Xref
	trailer Dict
}

func (p *Pdf) Version() string {
	return p.version
}

func (p *Pdf) SetVersion(version string) {
	p.version = version
}

func (p *Pdf) Objects() []Object {
	return p.objects
}

func (p *Pdf) Object(id int) Object {
	for _, obj := range p.objects {
		if obj.Ref.Id == id {
			return obj
		}
	}

	return Object{}
}

func (p *Pdf) AddObject(obj Object) {
	p.objects = append(p.objects, obj)
}

func (p *Pdf) Xrefs() []Xref {
	return p.xrefs
}

func (p *Pdf) Xref(id int64) Xref {
	for _, xref := range p.xrefs {
		if xref.Offset == id {
			return xref
		}
	}

	return Xref{}
}

func (p *Pdf) AddXref(xref Xref) {
	p.xrefs = append(p.xrefs, xref)
}

func (p *Pdf) Trailer() Dict {
	return p.trailer
}

func (p *Pdf) SetTrailer(trailer Dict) {
	p.trailer = trailer
}
