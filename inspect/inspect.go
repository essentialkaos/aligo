package inspect

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"go/ast"
	"go/token"
	"go/types"
	"path"
	"reflect"
	"slices"
	"sort"
	"strings"

	"github.com/essentialkaos/ek/v13/sliceutil"

	"github.com/kisielk/gotool"

	"golang.org/x/tools/go/packages"

	"github.com/essentialkaos/aligo/v2/i18n"
	"github.com/essentialkaos/aligo/v2/report"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const IGNORE_FLAG = "aligo:ignore"

// ////////////////////////////////////////////////////////////////////////////////// //

// Sizes contains info about WordSize and MaxAlign
var Sizes types.Sizes

// ////////////////////////////////////////////////////////////////////////////////// //

type structInfo struct {
	Name     string
	Type     *types.Struct
	AST      *ast.StructType
	Pos      token.Position
	Mappings map[string]string
	Skip     bool
}

// ////////////////////////////////////////////////////////////////////////////////// //

var fileSet *token.FileSet

// ////////////////////////////////////////////////////////////////////////////////// //

// ProcessSources starts sources processing
func ProcessSources(dirs, tags, excludes []string) (*report.Report, error) {
	importPaths := sliceutil.Filter(gotool.ImportPaths(dirs), func(importPath string, _ int) bool {
		return !slices.ContainsFunc(excludes, func(exclude string) bool {
			return strings.Contains(importPath, exclude)
		})
	})

	if len(importPaths) == 0 {
		return nil, i18n.UI.ERRORS.NO_IMPORT_PATHS.Error()
	}

	fileSet = token.NewFileSet()

	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Fset:  fileSet,
		Tests: false,
	}, importPaths...)

	if err != nil {
		return nil, err
	}

	return processPackages(pkgs)
}

// GetMaxAlign returns MaxAlign
func GetMaxAlign() int64 {
	if Sizes == nil {
		return 8
	}

	t, ok := Sizes.(*types.StdSizes)

	if ok {
		return t.MaxAlign
	}

	// Get MaxAlign from private struct like *types.gcSizes
	ptr := reflect.ValueOf(Sizes)

	if ptr.IsValid() {
		f := reflect.Indirect(ptr).FieldByName("MaxAlign")

		if f.IsValid() {
			return f.Int()
		}
	}

	return 8
}

// ////////////////////////////////////////////////////////////////////////////////// //

// processPackages checks given packages and returns report for them
func processPackages(pkgs []*packages.Package) (*report.Report, error) {
	result := &report.Report{}

	for _, pkg := range pkgs {
		pkgInfo, err := processPackage(pkg)

		if err != nil {
			return nil, err
		}

		result.Packages = append(result.Packages, pkgInfo)
	}

	return result, nil
}

// processPackage checks given package and returns report for it
func processPackage(pkg *packages.Package) (*report.Package, error) {
	var strName string
	var strPos token.Position
	var strIgnore bool

	result := &report.Package{Path: pkg.ID}
	mappings := map[string]string{pkg.ID + ".": ""}

	for _, file := range pkg.Syntax {
		commentMap := ast.NewCommentMap(fileSet, file, file.Comments)

		ast.Inspect(file, func(node ast.Node) bool {
			switch nt := node.(type) {
			case *ast.GenDecl:
				if nt.Tok == token.TYPE {
					decl := nt.Specs[0].(*ast.TypeSpec)
					strName = decl.Name.Name
					strPos = fileSet.Position(nt.TokPos)
					strIgnore = checkIgnoreFlag(commentMap.Filter(nt))
				}

			case *ast.ImportSpec:
				ntPath := strings.Trim(nt.Path.Value, `"`)

				if strings.Contains(ntPath, "/") {
					if nt.Name == nil {
						mappings[ntPath] = formatPackageName(ntPath)
					} else {
						mappings[ntPath] = nt.Name.Name
					}
				}

			case *ast.StructType:
				if strName == "" {
					return true // ignore unnamed structs defined in methods
				}

				info := &structInfo{
					Name:     strName,
					Type:     pkg.TypesInfo.Types[nt].Type.(*types.Struct),
					AST:      nt,
					Pos:      strPos,
					Mappings: mappings,
					Skip:     strIgnore,
				}

				structReport := getStructReport(info)

				if structReport != nil {
					result.Structs = append(result.Structs, structReport)
				}

				strName = ""
			}

			return true
		})
	}

	return result, nil
}

// getStructInfo parses struct info and calculates size
func getStructReport(info *structInfo) *report.Struct {
	result := &report.Struct{
		Name:     info.Name,
		Position: convertPosition(info.Pos),
		Ignore:   info.Skip,
	}

	numFields := info.Type.NumFields()

	// Recover from panic of checking size of non-generic types
	defer func() { recover() }()

	for i := range numFields {
		f := info.Type.Field(i)
		fs := findFieldInfo(info.AST.Fields.List, i, f.Name())
		utyp := f.Type().Underlying()
		size := Sizes.Sizeof(utyp)
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

// findFieldInfo tries to find field info in fields slice
func findFieldInfo(list []*ast.Field, index int, name string) *ast.Field {
	for _, field := range list {
		for _, ident := range field.Names {
			if ident.Name == name {
				return field
			}
		}
	}

	return list[index]
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
			if v == "." {
				k, v = k+".", "" // Format local type name
			}

			typ = strings.ReplaceAll(typ, k, v)
		}
	}

	if strings.ContainsRune(typ, '/') {
		return path.Base(typ)
	}

	return typ
}

// formatPackageName formats package name
func formatPackageName(p string) string {
	p = path.Base(p)
	p = strings.ReplaceAll(p, "go-", "")
	p = strings.ReplaceAll(p, "go.", "")

	if strings.Contains(p, ".") {
		p = p[:strings.Index(p, ".")]
	}

	return p
}

// checkIgnoreFlag checks struct comments for ignore flag
func checkIgnoreFlag(cm ast.CommentMap) bool {
	if cm == nil || len(cm.Comments()) == 0 {
		return false
	}

	for _, cg := range cm.Comments() {
		for _, c := range cg.List {
			if strings.Contains(strings.ToLower(c.Text), IGNORE_FLAG) {
				return true
			}
		}
	}

	return false
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
