package support

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fmtutil"
	"github.com/essentialkaos/ek/v12/system"
	"github.com/essentialkaos/ek/v12/system/container"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// showOSInfo shows verbose information about system
func showOSInfo() {
	osInfo, err := system.GetOSInfo()

	if err == nil {
		fmtutil.Separator(false, "OS INFO")

		printInfo(12, "Name", osInfo.ColoredName())
		printInfo(12, "Pretty Name", osInfo.ColoredPrettyName())
		printInfo(12, "Version", osInfo.Version)
		printInfo(12, "ID", osInfo.ID)
		printInfo(12, "ID Like", osInfo.IDLike)
		printInfo(12, "Version ID", osInfo.VersionID)
		printInfo(12, "Version Code", osInfo.VersionCodename)
		printInfo(12, "Platform ID", osInfo.PlatformID)
		printInfo(12, "CPE", osInfo.CPEName)
	}

	systemInfo, err := system.GetSystemInfo()

	if err != nil {
		return
	} else if osInfo == nil {
		fmtutil.Separator(false, "SYSTEM INFO")
		printInfo(12, "Name", systemInfo.OS)
	}

	printInfo(12, "Arch", systemInfo.Arch)
	printInfo(12, "Kernel", systemInfo.Kernel)

	containerEngine := "No"

	switch container.GetEngine() {
	case container.DOCKER:
		containerEngine = "Yes (Docker)"
	case container.PODMAN:
		containerEngine = "Yes (Podman)"
	case container.LXC:
		containerEngine = "Yes (LXC)"
	}

	fmtc.NewLine()

	printInfo(12, "Container", containerEngine)
}
