package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2019 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"go/build"
	"go/types"
	"os"
	"runtime"

	"pkg.re/essentialkaos/ek.v11/fmtc"
	"pkg.re/essentialkaos/ek.v11/fmtutil"
	"pkg.re/essentialkaos/ek.v11/options"
	"pkg.re/essentialkaos/ek.v11/strutil"
	"pkg.re/essentialkaos/ek.v11/usage"
	"pkg.re/essentialkaos/ek.v11/usage/update"

	"github.com/essentialkaos/aligo/inspect"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "aligo"
	VER  = "1.2.0"
	DESC = "Utility for viewing and checking Golang struct alignment"
)

// Constants with options names
const (
	OPT_ARCH     = "a:arch"
	OPT_STRUCT   = "s:struct"
	OPT_DETAILED = "d:detailed"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options map
var optMap = options.Map{
	OPT_ARCH:     {},
	OPT_STRUCT:   {},
	OPT_DETAILED: {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VER:      {Type: options.BOOL, Alias: "ver"},
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Init is main CLI func
func Init() {
	runtime.GOMAXPROCS(2)

	args, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError("Options parsing errors:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}

	configureUI()

	if options.GetB(OPT_HELP) || len(args) < 2 {
		showUsage()
		return
	}

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	prepare()
	process(args)
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	strutil.EllipsisSuffix = "…"
	fmtutil.SeparatorTitleColorTag = "{*}"
}

// prepare configure inspector
func prepare() {
	arch := build.Default.GOARCH

	if options.Has(OPT_ARCH) {
		arch = options.GetS(arch)
	}

	inspect.Sizes = types.SizesFor("gc", arch)

	if inspect.Sizes == nil {
		printErrorAndExit("Uknown arch %s", arch)
	}
}

// process starts processing
func process(args []string) {
	cmd := args[0]
	dirs := args[1:]

	report, err := inspect.ProcessSources(dirs)

	if err != nil {
		printErrorAndExit(err.Error())
	}

	if report == nil && err == nil {
		return
	}

	switch cmd {
	case "view", "v":
		if options.Has(OPT_STRUCT) {
			PrintStruct(report, options.GetS(OPT_STRUCT), options.GetB(OPT_DETAILED), false)
		} else {
			PrintFull(report)
		}
	case "check", "c":
		if options.Has(OPT_STRUCT) {
			PrintStruct(report, options.GetS(OPT_STRUCT), options.GetB(OPT_DETAILED), true)
		} else {
			Check(report, options.GetB(OPT_DETAILED))
		}
	default:
		printErrorAndExit("Command %s is unsupported", cmd)
	}
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage print usage info
func showUsage() {
	info := usage.NewInfo("", "package…")

	info.AddCommand("check", "Check package for alignment problems")
	info.AddCommand("view", "Print alignment info for all structs")

	info.AddOption(OPT_ARCH, "Architecture name", "name")
	info.AddOption(OPT_STRUCT, "Print info only about struct with given name", "name")
	info.AddOption(OPT_DETAILED, "Print detailed alignment info {s-}(useful with check command){!}")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"view .", "Show info about all structs in current package",
	)

	info.AddExample(
		"check .", "Check current package for alignment problems",
	)

	info.AddExample(
		"-s PostMessageParameters view .",
		"Show info about PostMessageParameters struct",
	)

	info.Render()
}

// showAbout print info about version
func showAbout() {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "Essential Kaos",
		License:       "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/aligo", update.GitHubChecker},
	}

	about.Render()
}
