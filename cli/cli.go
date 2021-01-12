package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2020 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"go/build"
	"go/types"
	"os"
	"runtime"

	"pkg.re/essentialkaos/ek.v12/fmtc"
	"pkg.re/essentialkaos/ek.v12/fmtutil"
	"pkg.re/essentialkaos/ek.v12/options"
	"pkg.re/essentialkaos/ek.v12/strutil"
	"pkg.re/essentialkaos/ek.v12/usage"
	"pkg.re/essentialkaos/ek.v12/usage/completion/bash"
	"pkg.re/essentialkaos/ek.v12/usage/completion/fish"
	"pkg.re/essentialkaos/ek.v12/usage/completion/zsh"
	"pkg.re/essentialkaos/ek.v12/usage/man"
	"pkg.re/essentialkaos/ek.v12/usage/update"

	"github.com/essentialkaos/aligo/inspect"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP  = "aligo"
	VER  = "1.4.0"
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

	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
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

	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
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

	if options.Has(OPT_COMPLETION) {
		os.Exit(genCompletion())
	}

	if options.Has(OPT_GENERATE_MAN) {
		genMan()
		os.Exit(0)
	}

	configureUI()

	if options.GetB(OPT_VER) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) || len(args) < 2 {
		showUsage()
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
			if Check(report, options.GetB(OPT_DETAILED)) {
				os.Exit(1)
			}
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

// showUsage prints usage info
func showUsage() {
	genUsage().Render()
}

// showAbout prints info about version
func showAbout() {
	genAbout().Render()
}

// genCompletion generates completion for different shells
func genCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Printf(bash.Generate(info, "aligo"))
	case "fish":
		fmt.Printf(fish.Generate(info, "aligo"))
	case "zsh":
		fmt.Printf(zsh.Generate(info, optMap, "aligo"))
	default:
		return 1
	}

	return 0
}

// genMan generates man page
func genMan() {
	fmt.Println(
		man.Generate(
			genUsage(),
			genAbout(),
		),
	)
}

// genUsage generates usage info
func genUsage() *usage.Info {
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

	return info
}

// genAbout generates info about version
func genAbout() *usage.About {
	about := &usage.About{
		App:           APP,
		Version:       VER,
		Desc:          DESC,
		Year:          2009,
		Owner:         "ESSENTIAL KAOS",
		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/aligo", update.GitHubChecker},
	}

	return about
}
