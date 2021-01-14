package inspect

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path"
	"sort"
	"strings"

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

type structInfo struct {
	Name     string
	Type     *types.Struct
	AST      *ast.StructType
	Pos      token.Position
	Mappings map[string]string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// ProcessSources starts sources processing
func ProcessSources(dirs []string) (*report.Report, error) {
	importPaths := gotool.ImportPaths(dirs)

	if len(importPaths) == 0 {
		return nil, fmt.Errorf("No import paths found")
	}

	fileSet = token.NewFileSet()

	loaderConfig := &loader.Config{ParserMode: parser.ParseComments}
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

	sort.Sort(pkgSlice(result.Packages))

	return result, nil
}

// processPackage extracts package data
func processPackage(pkg *loader.PackageInfo) (*report.Package, error) {
	var strName string
	var strPos token.Position

	result := &report.Package{Path: pkg.Pkg.Path()}

	mappings := map[string]string{
		pkg.Pkg.Path() + ".": "",
	}

	for _, file := range pkg.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch nt := node.(type) {
			case *ast.GenDecl:
				if nt.Tok == token.TYPE {
					decl := nt.Specs[0].(*ast.TypeSpec)
					strName = decl.Name.Name
					strPos = fileSet.Position(nt.TokPos)
				}

			case *ast.ImportSpec:
				ntPath := strings.Trim(nt.Path.Value, "\"")

				if strings.Contains(ntPath, "/") {
					if nt.Name == nil {
						mappings[ntPath] = formatPackageName(ntPath)
					} else {
						mappings[ntPath] = nt.Name.Name
					}
				}

			case *ast.StructType:
				info := &structInfo{
					Name:     strName,
					Type:     pkg.Types[nt].Type.(*types.Struct),
					AST:      nt,
					Pos:      strPos,
					Mappings: mappings,
				}

				result.Structs = append(result.Structs, getStructInfo(info))
			}

			return true
		})
	}

	return result, nil
}

// getStructInfo parses struct info and calculates size
func getStructInfo(info *structInfo) *report.Struct {
	result := &report.Struct{
		Name:     info.Name,
		Position: convertPosition(info.Pos),
	}

	numFields := info.Type.NumFields()

	for i := 0; i < numFields; i++ {
		f := info.Type.Field(i)
		fs := info.AST.Fields.List[i]
		size := Sizes.Sizeof(f.Type())
		comm := strings.Trim(fs.Comment.Text(), "\n\r")
		typ := formatValueType(f.Type().String(), info.Mappings)

		result.Fields = append(
			result.Fields,
			&report.Field{
				Name:    f.Name(),
				Type:    typ,
				Tag:     info.Type.Tag(i),
				Comment: comm,
				Size:    size,
			},
		)
	}

	result.Size = Sizes.Sizeof(info.Type)

	alnSize, alnFields := getAlignedFields(info.Type, result.Fields)

	if alnSize != result.Size {
		result.OptimalSize = alnSize
		result.AlignedFields = alnFields
	} else {
		result.OptimalSize = result.Size
	}

	return result
}

// getAlignedFields tries to find optimal field order
func getAlignedFields(str *types.Struct, origFields []*report.Field) (int64, []*report.Field) {
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

// formatValueType formats value type
func formatValueType(typ string, mappings map[string]string) string {
	for k, v := range mappings {
		if strings.Contains(typ, k) {
			return strings.Replace(typ, k, v, -1)
		}
	}

	return typ
}

// formatPackageName formats package name
func formatPackageName(p string) string {
	p = path.Base(p)
	p = strings.Replace(p, "go-", "", -1)
	p = strings.Replace(p, "go.", "", -1)

	if strings.Contains(p, ".") {
		p = p[:strings.Index(p, ".")]
	}

	return p
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

type pkgSlice []*report.Package

func (s pkgSlice) Len() int           { return len(s) }
func (s pkgSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s pkgSlice) Less(i, j int) bool { return s[i].Path < s[j].Path }

// ////////////////////////////////////////////////////////////////////////////////// //
