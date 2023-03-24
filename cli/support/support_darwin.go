package support

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/system"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// showOSInfo shows verbose information about system
func showOSInfo() {
	systemInfo, err := system.GetSystemInfo()

	if err != nil {
		return
	}

	osInfo, err := system.GetOSInfo()

	if err == nil {
		fmtutil.Separator(false, "OS INFO")

		printInfo(12, "Name", osInfo.Name)
		printInfo(12, "Version", osInfo.VersionID)
		printInfo(12, "Build", osInfo.Build)
	}

	fmtutil.Separator(false, "SYSTEM INFO")

	printInfo(7, "Name", systemInfo.OS)
	printInfo(7, "Arch", systemInfo.Arch)
	printInfo(7, "Kernel", systemInfo.Kernel)
}
