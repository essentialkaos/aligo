package report

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
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
	AlignedFields []*Field `json:"aligned_fields"` // nil if Size == OptimalSize
	Size          int64    `json:"size"`
	OptimalSize   int64    `json:"optimal_size"`
	Ignore        bool     `json:ignore`
}

// Field contains info about field
type Field struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Tag     string `json:"tag"`
	Comment string `json:"comment"`
	Size    int64  `json:"size"`
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
