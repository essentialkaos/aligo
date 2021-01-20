package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fmtutil"
	"pkg.re/essentialkaos/ek.v12/mathutil"
	"pkg.re/essentialkaos/ek.v12/strutil"

	"github.com/essentialkaos/aligo/inspect"
	"github.com/essentialkaos/aligo/report"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// MAX_TYPE_SIZE maximum type size
const MAX_TYPE_SIZE = 32

// ////////////////////////////////////////////////////////////////////////////////// //

type Renderer struct {
	format      string
	hasTags     bool
	hasComments bool
	detailed    bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

// PrintFull prints all report info
func PrintFull(r *report.Report) {
	if isEmptyReport(r) {
		return
	}

	for _, pkg := range r.Packages {
		if !pkg.IsEmpty() {
			printPackageInfo(pkg, false)
		}
	}

	fmtutil.Separator(true)
	fmtc.NewLine()
}

// PrintStruct prints info about struct
func PrintStruct(r *report.Report, strName string, optimal bool) {
	if isEmptyReport(r) {
		return
	}

	if strName == "" {
		printWarn("You should define struct name")
		return
	}

	pkg, str := findStruct(r, strName)

	if pkg == nil && str == nil {
		printWarn("Can't find struct with name \"%s\"", strName)
	}

	printPackageSeparator(pkg.Path)
	printStructInfo(str, optimal)

	fmtutil.Separator(true)
	fmtc.NewLine()
}

// Check checks report for problems
func Check(r *report.Report) bool {
	if isEmptyReport(r) {
		return false
	}

	var hasProblems bool

	for _, pkg := range r.Packages {
		if pkg.IsEmpty() || !isPackageHasProblems(pkg) {
			continue
		}

		hasProblems = true

		printPackageInfo(pkg, true)
	}

	if !hasProblems {
		fmtc.Println("{g}All structs are well aligned{!}")
	} else {
		fmtutil.Separator(true)
		fmtc.NewLine()
	}

	return hasProblems
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewRenderer creates new filed info renderer
func NewRenderer(fields []*report.Field, detailed bool) *Renderer {
	r := &Renderer{detailed: detailed}

	var mName, mType, mTag, mComm int

	for _, field := range fields {
		mName = mathutil.Max(mName, len(field.Name))
		mType = mathutil.Max(mType, len(strutil.Ellipsis(field.Type, MAX_TYPE_SIZE)))
		mTag = mathutil.Max(mTag, len(field.Tag))
		mComm = mathutil.Max(mComm, len(field.Comment))
	}

	if mTag > 0 {
		mTag += 2
	}

	r.format = fmt.Sprintf("  %%-%ds", mName)

	if mTag > 0 || mComm > 0 || !detailed {
		r.format += fmt.Sprintf(" {*}%%-%ds{!}", mType)
	} else {
		r.format += fmt.Sprintf(" {*}%%s{!}")
	}

	if !detailed {
		return r
	}

	if mTag > 0 {
		r.hasTags = true

		if mComm > 0 {
			r.format += fmt.Sprintf(" {y}%%-%ds{!}", mTag)
		} else {
			r.format += fmt.Sprintf(" {y}%%s{!}")
		}
	}

	if mComm > 0 {
		r.hasComments = true
		r.format += fmt.Sprintf(" {s-}// %%s{!}")
	}

	return r
}

// PrintField prints field info
func (r *Renderer) PrintField(f *report.Field) {
	var fTag, fComment, fType string

	if f.Tag != "" {
		fTag = "`" + f.Tag + "`"
	}

	if r.hasComments && f.Comment == "" {
		fComment = "-"
	} else {
		fComment = f.Comment
	}

	fType = strutil.Ellipsis(f.Type, MAX_TYPE_SIZE)

	switch {
	case r.hasTags && r.hasComments:
		fmtc.Printf(r.format, f.Name, fType, fTag, fComment)
	case r.hasTags && !r.hasComments:
		fmtc.Printf(r.format, f.Name, fType, fTag)
	case !r.hasTags && r.hasComments:
		fmtc.Printf(r.format, f.Name, fType, fComment)
	default:
		fmtc.Printf(r.format, f.Name, fType)
	}
}

// PrintPlaceholder prints placeholder
func (r *Renderer) PrintPlaceholder() {
	switch {
	case r.hasTags && r.hasComments:
		fmtc.Printf(r.format, "", "", "", "")
	case r.hasTags || r.hasComments:
		fmtc.Printf(r.format, "", "", "")
	default:
		fmtc.Printf(r.format, "", "")
	}

	fmt.Printf("  ")
}

// ////////////////////////////////////////////////////////////////////////////////// //

// isEmptyReport returns true if report is empty
func isEmptyReport(r *report.Report) bool {
	if r.IsEmpty() {
		printWarn("Given package doesn't have any structs")
		return true
	}

	return false
}

// printPackageSeparator prints separator with package name
func printPackageSeparator(path string) {
	if strings.HasPrefix(path, ".") {
		path = "{GOPATH}" + path[1:]
	}

	fmtutil.Separator(false, path)
}

// printPackageInfo prints package info
func printPackageInfo(pkg *report.Package, onlyProblems bool) {
	printPackageSeparator(pkg.Path)

	for _, str := range pkg.Structs {
		if onlyProblems && isAlignedStruct(str) {
			continue
		}

		printStructInfo(str, onlyProblems == true)
	}
}

// printStructSizeInfo prints info about struct size
func printStructSizeInfo(str *report.Struct, optimal bool) {
	if optimal {
		fmtc.Printf(
			"Struct {*}%s{!} {s-}(%s:%d){!} fields order can be optimized (%d → %d)\n\n",
			str.Name, str.Position.File, str.Position.Line, str.Size, str.OptimalSize,
		)
	} else {
		if str.Size != str.OptimalSize {
			fmtc.Printf(
				"{s-}// %s:%d | Size: %d (Optimal: %d){!}\n",
				str.Position.File, str.Position.Line, str.Size, str.OptimalSize,
			)
		} else {
			fmtc.Printf(
				"{s-}// %s:%d | Size: %d{!}\n",
				str.Position.File, str.Position.Line, str.Size,
			)
		}
	}
}

// printStructInfo prints struct info
func printStructInfo(str *report.Struct, optimal bool) {
	printStructSizeInfo(str, optimal)

	fmtc.Printf("type {*}%s{!} struct {s}{{!}\n", str.Name)

	if optimal {
		printAlignedFieldsInfo(str.AlignedFields)
	} else {
		printCurrentFieldsInfo(str.Fields)
	}

	fmtc.Println("{s}}{!}\n")
}

func printAlignedFieldsInfo(fields []*report.Field) {
	r := NewRenderer(fields, true)

	for _, field := range fields {
		r.PrintField(field)
		fmtc.NewLine()
	}
}

func printCurrentFieldsInfo(fields []*report.Field) {
	r := NewRenderer(fields, false)

	counter := int64(0)
	maxAlign := inspect.GetMaxAlign()

	for index, field := range fields {
		r.PrintField(field)

		fmt.Printf(strings.Repeat("  ", int(counter+1)))

		for i := int64(0); i < field.Size; i++ {
			fmtc.Printf("{g}■ {!}")

			counter++

			if counter == maxAlign {
				if i+1 != field.Size {
					fmtc.NewLine()
					r.PrintPlaceholder()
				}
				counter = 0
			}
		}

		if index+1 < len(fields) && counter != 0 && fields[index+1].Size > maxAlign-counter {
			fmtc.Printf(strings.Repeat("{r}□ {!}", int(maxAlign-counter)))
			counter = 0
		} else if index+1 == len(fields) && counter != 0 {
			fmtc.Printf(strings.Repeat("{g}□ {!}", int(maxAlign-counter)))
		}

		fmtc.NewLine()
	}
}

// findStruct finds struct with given name
func findStruct(r *report.Report, name string) (*report.Package, *report.Struct) {
	for _, pkg := range r.Packages {
		for _, str := range pkg.Structs {
			if str.Name == name {
				return pkg, str
			}
		}
	}

	return nil, nil
}

// isPackageHasProblems returns true if package has structs with
// unaligned fields
func isPackageHasProblems(pkg *report.Package) bool {
	for _, str := range pkg.Structs {
		if !isAlignedStruct(str) {
			return true
		}
	}

	return false
}

// isAlignedStruct returns false if struct has unaligned fields
func isAlignedStruct(str *report.Struct) bool {
	return str.Size == str.OptimalSize || str.Ignore
}
