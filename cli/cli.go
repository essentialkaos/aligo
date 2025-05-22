package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
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

	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fmtutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/pager"
	"github.com/essentialkaos/ek/v13/strutil"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/apps"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"

	"github.com/essentialkaos/aligo/v2/cli/i18n"
	"github.com/essentialkaos/aligo/v2/inspect"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// App info
const (
	APP = "aligo"
	VER = "2.2.5"
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

	if !errs.IsEmpty() {
		terminal.Error(i18n.UI.ERRORS.OPTION_PARSING.Add("", ":"))
		terminal.Error(errs.Error("- "))
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
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			WithApps(apps.Golang()).
			Print()
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) < 2:
		genUsage().Print()
		os.Exit(0)
	}

	err, ok := process(args)

	if err != nil {
		terminal.Error(err)
	}

	if !ok {
		os.Exit(1)
	}
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.IsTrueColorSupported():
		colorTagApp, colorTagVer = "{*}{&}{#00ADD8}", "{#5DC9E2}"
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{&}{#38}", "{#74}"
	default:
		colorTagApp, colorTagVer = "{*}{&}{c}", "{c}"
	}

	i18n.SetLanguage()
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	strutil.EllipsisSuffix = "…"
	fmtutil.SeparatorSymbol = "–"
}

// prepare configures inspector
func prepare() error {
	arch := build.Default.GOARCH

	if options.Has(OPT_ARCH) {
		arch = options.GetS(OPT_ARCH)
	}

	inspect.Sizes = types.SizesFor("gc", arch)

	if inspect.Sizes == nil {
		return fmt.Errorf(i18n.UI.ERRORS.UNKNOWN_ARCH.String(), arch)
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
		return fmt.Errorf(i18n.UI.ERRORS.UNSUPPORTED_COMMAND.String(), cmd), false
	}

	return nil, true
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
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", i18n.UI.USAGE.ARGUMENTS.String())

	info.AppNameColorTag = colorTagApp

	info.AddCommand("check", i18n.UI.USAGE.COMMANDS.CHECK.String())
	info.AddCommand("view", i18n.UI.USAGE.COMMANDS.VIEW.String())

	info.AddOption(OPT_ARCH, i18n.UI.USAGE.OPTIONS.ARCH.String(), i18n.UI.USAGE.OPTIONS.ARCH_VAL.String())
	info.AddOption(OPT_STRUCT, i18n.UI.USAGE.OPTIONS.STRUCT.String(), i18n.UI.USAGE.OPTIONS.STRUCT_VAL.String())
	info.AddOption(OPT_TAGS, i18n.UI.USAGE.OPTIONS.TAGS.String(), i18n.UI.USAGE.OPTIONS.TAGS_VAL.String())
	info.AddOption(OPT_PAGER, i18n.UI.USAGE.OPTIONS.PAGER.String())
	info.AddOption(OPT_NO_COLOR, i18n.UI.USAGE.OPTIONS.NO_COLOR.String())
	info.AddOption(OPT_HELP, i18n.UI.USAGE.OPTIONS.HELP.String())
	info.AddOption(OPT_VER, i18n.UI.USAGE.OPTIONS.VER.String())

	info.AddExample(
		"view .",
		i18n.UI.USAGE.EXAMPLES.EXAMPLE_1.String(),
	)

	info.AddExample(
		"check .",
		i18n.UI.USAGE.EXAMPLES.EXAMPLE_2.String(),
	)

	info.AddExample(
		"check ./...",
		i18n.UI.USAGE.EXAMPLES.EXAMPLE_3.String(),
	)

	info.AddExample(
		"--tags tag1,tag2,tag3 check ./...",
		i18n.UI.USAGE.EXAMPLES.EXAMPLE_4.String(),
	)

	info.AddExample(
		"-s PostMessageParameters view .",
		i18n.UI.USAGE.EXAMPLES.EXAMPLE_5.String(),
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    i18n.UI.USAGE.DESC.String(),
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		License:       "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
		UpdateChecker: usage.UpdateChecker{"essentialkaos/aligo", update.GitHubChecker},
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
