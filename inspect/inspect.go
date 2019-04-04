package inspect

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"sort"

	"github.com/kisielk/gotool"

	"golang.org/x/tools/go/loader"

	"github.com/essentialkaos/aligo/report"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Sizes contains info about WordSize and MaxAlign
var Sizes types.Sizes

// ////////////////////////////////////////////////////////////////////////////////// //

var fileSet *token.FileSet

// ////////////////////////////////////////////////////////////////////////////////// //

// ProcessSources starts sources processing
func ProcessSources(dirs []string) (*report.Report, error) {
	importPaths := gotool.ImportPaths(dirs)

	if len(importPaths) == 0 {
		return nil, fmt.Errorf("No import paths found")
	}

	fileSet = token.NewFileSet()

	loaderConfig := &loader.Config{}
	loaderConfig.Fset = fileSet

	for _, importPath := range importPaths {
		loaderConfig.Import(importPath)
	}

	prog, err := loaderConfig.Load()

	if err != nil {
		return nil, nil
	}

	return processProgram(prog)
}

// GetMaxAlign returns MaxAlign
func GetMaxAlign() int64 {
	if Sizes == nil {
		return 8
	}

	return Sizes.(*types.StdSizes).MaxAlign
}

// ////////////////////////////////////////////////////////////////////////////////// //

// processProgram extracs data from all packages
func processProgram(prog *loader.Program) (*report.Report, error) {
	result := &report.Report{}
	packages := prog.InitialPackages()

	for _, pkg := range packages {
		pkgInfo, err := processPackage(pkg)

		if err != nil {
			return nil, err
		}

		result.Packages = append(result.Packages, pkgInfo)
	}

	return result, nil
}

// processPackage extracts package data
func processPackage(pkg *loader.PackageInfo) (*report.Package, error) {
	var strName string
	var strPos token.Position

	result := &report.Package{Path: pkg.Pkg.Path()}

	for _, file := range pkg.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch nt := node.(type) {
			case *ast.GenDecl:
				if nt.Tok == token.TYPE {
					decl := nt.Specs[0].(*ast.TypeSpec)
					strName = decl.Name.Name
					strPos = fileSet.Position(nt.TokPos)
				}
			case *ast.StructType:
				str := getStructInfo(strName, pkg.Types[nt].Type.(*types.Struct), strPos)
				result.Structs = append(result.Structs, str)
			}

			return true
		})
	}

	return result, nil
}

// getStructInfo parses struct info and calculates size
func getStructInfo(name string, str *types.Struct, pos token.Position) *report.Struct {
	result := &report.Struct{
		Name:     name,
		Position: convertPosition(pos),
	}

	numFields := str.NumFields()

	for i := 0; i < numFields; i++ {
		f := str.Field(i)
		size := Sizes.Sizeof(f.Type())

		result.Fields = append(
			result.Fields,
			&report.Field{
				Name: f.Name(),
				Type: f.Type().String(),
				Tag:  str.Tag(i),
				Size: size,
			},
		)
	}

	result.Size = Sizes.Sizeof(str)

	optSize, optFields := getOptimalFields(str, result.Fields)

	if optSize != result.Size {
		result.OptimalSize = optSize
		result.OptimalFields = optFields
	} else {
		result.OptimalSize = result.Size
	}

	return result
}

// getOptimalFields tries to find optimal field order
func getOptimalFields(str *types.Struct, origFields []*report.Field) (int64, []*report.Field) {
	numFields := len(origFields)
	fields := append(origFields[:0:0], origFields...)
	vars := make([]*types.Var, numFields)
	aligns := make([]int64, numFields)
	sizes := make([]int64, numFields)

	for i := 0; i < numFields; i++ {
		fieldVar := str.Field(i)
		vars[i] = fieldVar
		aligns[i] = Sizes.Alignof(fieldVar.Type())
		sizes[i] = fields[i].Size
	}

	sort.Stable(&optimalSorter{vars, fields, aligns, sizes})

	return Sizes.Sizeof(types.NewStruct(vars, nil)), fields
}

// convertPosition converts position between types
func convertPosition(pos token.Position) report.Position {
	return report.Position{
		File: path.Base(pos.Filename),
		Line: pos.Line,
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// optimalSorter is field sorter
type optimalSorter struct {
	Vars   []*types.Var
	Fields []*report.Field
	Aligns []int64
	Sizes  []int64
}

func (s *optimalSorter) Len() int {
	return len(s.Fields)
}

func (s *optimalSorter) Swap(i, j int) {
	s.Vars[i], s.Vars[j] = s.Vars[j], s.Vars[i]
	s.Fields[i], s.Fields[j] = s.Fields[j], s.Fields[i]
	s.Aligns[i], s.Aligns[j] = s.Aligns[j], s.Aligns[i]
	s.Sizes[i], s.Sizes[j] = s.Sizes[j], s.Sizes[i]
}

func (s *optimalSorter) Less(i, j int) bool {
	switch {
	case s.Sizes[i] == 0 && s.Sizes[j] != 0,
		s.Sizes[j] == 0 && s.Sizes[i] != 0:
		return false

	case s.Aligns[i] != s.Aligns[j]:
		return s.Aligns[i] > s.Aligns[j]

	case s.Sizes[i] != s.Sizes[j]:
		return s.Sizes[i] > s.Sizes[j]

	default:
		return false
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //
