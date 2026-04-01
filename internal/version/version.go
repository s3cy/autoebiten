package version

import "runtime/debug"

// VersionForCLI returns the CLI version from build info.
func VersionForCLI() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		return info.Main.Version
	} else {
		return "unknown"
	}
}

// VersionForLibrary assumes autoebiten is running as a library andreturns the library version from build info.
func VersionForLibrary() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}
	for _, dep := range info.Deps {
		if dep.Path == "github.com/s3cy/autoebiten" {
			if dep.Replace != nil {
				return dep.Replace.Path + " " + dep.Replace.Version
			}
			return dep.Version
		}
	}
	return "unknown"
}
