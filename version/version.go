package version

// the build version of safeline-utils
var (
	Version = VersionDefault
	Build   = BuildDefault
	Repo    = RepoDefault
)

// if not use make to build project, use the following default variables for version command
const (
	// VersionDefault is the default version of NEU_IPGW
	VersionDefault = "unknown"
	// BuildDefault is the default build time of NEU_IPGW
	BuildDefault = "1970-01-01T00:00:00+0000"
	RepoDefault  = "DoraTiger/safeline-utils"
)