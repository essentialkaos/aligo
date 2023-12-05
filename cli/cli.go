package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"go/build"
	"go/types"
	"os"
	"runtime"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/pager"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/usage"
	"github.com/essentialkaos/ek/v12/usage/completion/bash"
	"github.com/essentialkaos/ek/v12/usage/completion/fish"
	"github.com/essentialkaos/ek/v12/usage/completion/zsh"
	"github.com/essentialkaos/ek/v12/usage/man"
	"github.com/essentialkaos/ek/v12/usage/update"

	"github.com/essentialkaos/aligo/v2/cli/support"
	"github.com/essentialkaos/aligo/v2/inspect"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "aligo"
	VER  = "2.1.0"
	DESC = "Utility for viewing and checking Go struct alignment"
)

// Constants with options names
const (
	OPT_ARCH     = "a:arch"
	OPT_STRUCT   = "s:struct"
	OPT_TAGS     = "t:tags"
	OPT_PAGER    = "P:pager"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options map
var optMap = options.Map{
	OPT_ARCH:     {},
	OPT_STRUCT:   {},
	OPT_TAGS:     {Mergeble: true},
	OPT_PAGER:    {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

var colorTagApp, colorTagVer string

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main CLI function
func Run(gitRev string, gomod []byte) {
	runtime.GOMAXPROCS(2)

	preConfigureUI()

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Print(APP, VER, gitRev, gomod)
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) < 2:
		genUsage().Print()
		os.Exit(0)
	}

	err, ok := process(args)

	if err != nil {
		printError(err.Error())
	}

	if !ok {
		os.Exit(1)
	}
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if os.Getenv("NO_COLOR") != "" {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	strutil.EllipsisSuffix = "…"
	fmtutil.SeparatorSymbol = "–"

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{&}{#00ADD8}", "{#5DC9E2}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{&}{#38}", "{#74}"
	default:
		colorTagApp, colorTagVer = "{*}{&}{c}", "{c}"
	}
}

// prepare configures inspector
func prepare() error {
	arch := build.Default.GOARCH

	if options.Has(OPT_ARCH) {
		arch = options.GetS(OPT_ARCH)
	}

	inspect.Sizes = types.SizesFor("gc", arch)

	if inspect.Sizes == nil {
		return fmt.Errorf("Unknown arch %s", arch)
	}

	return nil
}

// process starts source code processing
func process(args options.Arguments) (error, bool) {
	err := prepare()

	if err != nil {
		return err, false
	}

	cmd := args.Get(0).ToLower().String()
	dirs := args.Strings()[1:]
	tags := strings.Split(options.GetS(OPT_TAGS), ",")

	report, err := inspect.ProcessSources(dirs, tags)

	if err != nil {
		return err, false
	}

	if report == nil {
		return nil, true
	}

	if options.GetB(OPT_PAGER) {
		if pager.Setup() == nil {
			defer pager.Complete()
		}
	}

	switch cmd {
	case "view", "v":
		if options.Has(OPT_STRUCT) {
			PrintStruct(report, options.GetS(OPT_STRUCT), false)
		} else {
			PrintFull(report)
		}

	case "check", "c":
		if options.Has(OPT_STRUCT) {
			PrintStruct(report, options.GetS(OPT_STRUCT), true)
		} else if !Check(report) {
			return nil, false
		}

	default:
		return fmt.Errorf("Command %s is unsupported", cmd), false
	}

	return nil, true
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, "aligo"))
	case "fish":
		fmt.Print(fish.Generate(info, "aligo"))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, "aligo"))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(""),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "package…")

	info.AppNameColorTag = colorTagApp

	info.AddCommand("check", "Check package for alignment problems")
	info.AddCommand("view", "Print alignment info for all structs")

	info.AddOption(OPT_ARCH, "Architecture name", "name")
	info.AddOption(OPT_STRUCT, "Print info only about struct with given name", "name")
	info.AddOption(OPT_TAGS, "Build tags {s-}(mergeble){!}", "tag…")
	info.AddOption(OPT_PAGER, "Use pager for long output")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"view .", "Show info about all structs in current package",
	)

	info.AddExample(
		"check .", "Check current package",
	)

	info.AddExample(
		"check ./...", "Check current package and all sub-packages",
	)

	info.AddExample(
		"--tags tag1,tag2,tag3 check ./...", "Check current package and all sub-packages with custom build tags",
	)

	info.AddExample(
		"-s PostMessageParameters view .",
		"Show info about PostMessageParameters struct",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "—",

		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/aligo", update.GitHubChecker},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
