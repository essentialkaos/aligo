package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"strings"

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fmtutil"
	"github.com/essentialkaos/ek/v13/mathutil"
	"github.com/essentialkaos/ek/v13/terminal"

	"github.com/essentialkaos/aligo/v2/cli/i18n"
	"github.com/essentialkaos/aligo/v2/inspect"
	"github.com/essentialkaos/aligo/v2/report"
)

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
}

// PrintStruct prints info about struct
func PrintStruct(r *report.Report, strName string, optimal bool) {
	if isEmptyReport(r) {
		return
	}

	if strName == "" {
		terminal.Warn(i18n.UI.ERRORS.EMPTY_STRUCT_NAME)
		return
	}

	pkg, str := findStruct(r, strName)

	if pkg == nil && str == nil {
		terminal.Warn(i18n.UI.ERRORS.NO_STRUCT, strName)
		return
	}

	printPackageSeparator(pkg.Path)
	printStructInfo(str, optimal)
}

// Check checks report for problems
func Check(r *report.Report) bool {
	if isEmptyReport(r) {
		return true
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
		fmtc.Println(i18n.UI.INFO.ALL_OPTIMAL.String())
		return true
	}

	return false
}

// ////////////////////////////////////////////////////////////////////////////////// //

// NewRenderer creates new filed info renderer
func NewRenderer(fields []*report.Field, detailed bool) *Renderer {
	r := &Renderer{detailed: detailed}

	var mName, mType, mTag, mComm int

	for _, f := range fields {
		mName = mathutil.Max(mName, len(f.Name))
		mType = mathutil.Max(mType, len(f.Type))
		mTag = mathutil.Max(mTag, len(f.Tag))
		mComm = mathutil.Max(mComm, len(f.Comment))
	}

	if mTag > 0 {
		mTag += 2
	}

	r.format = fmt.Sprintf("    %%-%ds", mName)

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
	var fTag, fComment string

	if f.Tag != "" {
		fTag = "`" + f.Tag + "`"
	}

	if r.hasComments && f.Comment == "" {
		fComment = "-"
	} else {
		fComment = f.Comment
	}

	switch {
	case r.hasTags && r.hasComments:
		fmtc.Printf(r.format, f.Name, f.Type, fTag, fComment)
	case r.hasTags && !r.hasComments:
		fmtc.Printf(r.format, f.Name, f.Type, fTag)
	case !r.hasTags && r.hasComments:
		fmtc.Printf(r.format, f.Name, f.Type, fComment)
	default:
		fmtc.Printf(r.format, f.Name, f.Type)
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
		terminal.Warn(i18n.UI.ERRORS.NO_ANY_STRUCTS)
		return true
	}

	return false
}

// printPackageSeparator prints separator with package name
func printPackageSeparator(path string) {
	if strings.HasPrefix(path, ".") {
		path = "{GOPATH}" + path[1:]
	}

	fmtutil.Separator(true)
	fmtc.Printfn(" ▾ {s}%s{!}", path)
	fmtutil.Separator(true)
	fmtc.NewLine()
}

// printPackageInfo prints package info
func printPackageInfo(pkg *report.Package, onlyProblems bool) {
	printPackageSeparator(pkg.Path)

	for _, str := range pkg.Structs {
		if onlyProblems && isAlignedStruct(str) {
			continue
		}

		printStructInfo(str, onlyProblems)
	}
}

// printStructSizeInfo prints info about struct size
func printStructSizeInfo(str *report.Struct, optimal bool) {
	if optimal {
		fmtc.Printf(
			i18n.UI.INFO.OPTIMIZE_ADVICE.Add("", "\n\n"),
			str.Name, str.Position.File, str.Position.Line, str.Size, str.OptimalSize,
		)
	} else {
		if str.Size != str.OptimalSize {
			fmtc.Printf(
				i18n.UI.INFO.WITH_OPTIMAL.Add("  ", "\n"),
				str.Position.File, str.Position.Line, str.Size, str.OptimalSize,
			)
		} else {
			fmtc.Printf(
				i18n.UI.INFO.ALREADY_OPTIMAL.Add("  ", "\n"),
				str.Position.File, str.Position.Line, str.Size,
			)
		}
	}
}

// printStructInfo prints struct info
func printStructInfo(str *report.Struct, optimal bool) {
	printStructSizeInfo(str, optimal)

	if str.Size == 0 {
		fmtc.Printfn("  type {&}{*}%s{!} struct {s}{ }{!}\n", str.Name)
		return
	}

	fmtc.Printfn("  type {&}{*}%s{!} struct {s}{{!}", str.Name)

	if optimal {
		printAlignedFieldsInfo(str.AlignedFields)
	} else {
		printCurrentFieldsInfo(str.Fields)
	}

	fmtc.Println("  {s}}{!}\n")
}

// printAlignedFieldsInfo prints aligned field data
func printAlignedFieldsInfo(fields []*report.Field) {
	r := NewRenderer(fields, true)

	for _, field := range fields {
		r.PrintField(field)
		fmtc.NewLine()
	}
}

// printCurrentFieldsInfo prints current field data
func printCurrentFieldsInfo(fields []*report.Field) {
	r := NewRenderer(fields, false)

	counter := int64(0)
	maxAlign := inspect.GetMaxAlign()

	for index, field := range fields {
		r.PrintField(field)

		fmt.Print(strings.Repeat("  ", int(counter+1)))

		// Ensure we don't panic when hitting an empty struct{}.
		if field.Size == 0 {
			fmtc.NewLine()
			continue
		}

		for counter%field.Size != 0 {
			fmtc.Printf("{r}□{!} ")
			counter++
		}

		for i := int64(0); i < field.Size; i++ {
			fmtc.Printf("{g}■{!} ")

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
			fmtc.Printf(strings.Repeat("{r}□{!} ", int(maxAlign-counter)))
			counter = 0
		} else if index+1 == len(fields) && counter != 0 {
			fmtc.Printf(strings.Repeat("{g}□{!} ", int(maxAlign-counter)))
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
