package report

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

// Report contains aligning info about packages
type Report struct {
	Packages []*Package `json:"packages"`
}

// Package contains info about all structs in package
type Package struct {
	Path    string    `json:"path"`
	Structs []*Struct `json:"structs"`
}

// Struct contains info about fields aligning
type Struct struct {
	Name          string   `json:"name"`
	Position      Position `json:"position"`
	Fields        []*Field `json:"fields"`
	OptimalFields []*Field `json:"optimal_fields"` // nil if Size == OptimalSize
	Size          int64    `json:"size"`
	OptimalSize   int64    `json:"optimal_size"`
}

// Field contains info about field
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Tag  string `json:"tag"`
	Size int64  `json:"size"`
}

// Position contains info about struct position
type Position struct {
	File string `json:"file"`
	Line int    `json:"line"`
}

// ////////////////////////////////////////////////////////////////////////////////// //

// IsEmpty returns true if report is empty
func (r *Report) IsEmpty() bool {
	if r == nil || len(r.Packages) == 0 {
		return true
	}

	for _, pkg := range r.Packages {
		if !pkg.IsEmpty() {
			return false
		}
	}

	return true
}

// IsEmpty returns true if package is empty
func (p *Package) IsEmpty() bool {
	if p == nil {
		return true
	}

	return len(p.Structs) == 0
}
