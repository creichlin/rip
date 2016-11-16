package rip

type PathElement interface {
	Template() string
}

type Parameter interface {
	Name() string
	Doc() string
}

type fixedPath string

func (f *fixedPath) Template() string {
	return string(*f)
}

type pathParameter struct {
	name string
	doc  string
}

func (f *pathParameter) Template() string {
	return "{" + f.name + "}"
}

func (f *pathParameter) Name() string {
	return f.name
}

func (f *pathParameter) Doc() string {
	return f.doc
}

type queryParameter struct {
	name string
	doc  string
}

func (f *queryParameter) Name() string {
	return f.name
}

func (f *queryParameter) Doc() string {
	return f.doc
}

type Validation interface{}