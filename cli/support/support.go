package support

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/hash"
	"github.com/essentialkaos/ek/v12/strutil"

	"github.com/essentialkaos/depsy"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Print prints verbose info about application, system, dependencies and
// important environment
func Print(app, ver, gitRev string, gomod []byte) {
	fmtutil.SeparatorTitleColorTag = "{s-}"
	fmtutil.SeparatorFullscreen = false
	fmtutil.SeparatorColorTag = "{s-}"
	fmtutil.SeparatorSize = 80

	showApplicationInfo(app, ver, gitRev)
	showOSInfo()
	showEnvInfo()
	showDepsInfo(gomod)

	fmtutil.Separator(false)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showApplicationInfo shows verbose information about application
func showApplicationInfo(app, ver, gitRev string) {
	fmtutil.Separator(false, "APPLICATION INFO")

	printInfo(7, "Name", app)
	printInfo(7, "Version", ver)

	printInfo(7, "Go", fmtc.Sprintf(
		"%s {s}(%s/%s){!}",
		strings.TrimLeft(runtime.Version(), "go"),
		runtime.GOOS, runtime.GOARCH,
	))

	if gitRev == "" {
		gitRev = extractGitRevFromBuildInfo()
	}

	if gitRev != "" {
		if !fmtc.DisableColors && fmtc.IsTrueColorSupported() {
			printInfo(7, "Git SHA", gitRev+getHashColorBullet(gitRev))
		} else {
			printInfo(7, "Git SHA", gitRev)
		}
	}

	bin, _ := os.Executable()
	binSHA := hash.FileHash(bin)

	if binSHA != "" {
		binSHA = strutil.Head(binSHA, 7)
		if !fmtc.DisableColors && fmtc.IsTrueColorSupported() {
			printInfo(7, "Bin SHA", binSHA+getHashColorBullet(binSHA))
		} else {
			printInfo(7, "Bin SHA", binSHA)
		}
	}
}

// showEnvInfo shows info about environment
func showEnvInfo() {
	fmtutil.Separator(false, "ENVIRONMENT")

	cmd := exec.Command("go", "version")
	out, err := cmd.Output()

	if err != nil {
		printInfo(2, "Go", "")
		return
	}

	goVer := string(out)
	goVer = strutil.ReadField(goVer, 2, false, ' ')
	goVer = strutil.Exclude(goVer, "go")

	printInfo(2, "Go", goVer)
}

// showDepsInfo shows information about all dependencies
func showDepsInfo(gomod []byte) {
	deps := depsy.Extract(gomod, false)

	if len(deps) == 0 {
		return
	}

	fmtutil.Separator(false, "DEPENDENCIES")

	for _, dep := range deps {
		if dep.Extra == "" {
			fmtc.Printf(" {s}%8s{!}  %s\n", dep.Version, dep.Path)
		} else {
			fmtc.Printf(" {s}%8s{!}  %s {s-}(%s){!}\n", dep.Version, dep.Path, dep.Extra)
		}
	}
}

// extractGitRevFromBuildInfo extracts git SHA from embedded build info
func extractGitRevFromBuildInfo() string {
	info, ok := debug.ReadBuildInfo()

	if !ok {
		return ""
	}

	for _, s := range info.Settings {
		if s.Key == "vcs.revision" && len(s.Value) > 7 {
			return s.Value[:7]
		}
	}

	return ""
}

// getHashColorBullet return bullet with color from hash
func getHashColorBullet(v string) string {
	if len(v) > 6 {
		v = strutil.Head(v, 6)
	}

	return fmtc.Sprintf(" {#" + strutil.Head(v, 6) + "}● {!}")
}

// printInfo formats and prints info record
func printInfo(size int, name, value string) {
	name += ":"
	size++

	if value == "" {
		fm := fmt.Sprintf("  {*}%%-%ds{!}  {s-}—{!}\n", size)
		fmtc.Printf(fm, name)
	} else {
		fm := fmt.Sprintf("  {*}%%-%ds{!}  %%s\n", size)
		fmtc.Printf(fm, name, value)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //
