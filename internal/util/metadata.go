package util

// BuildMeta contains information about the current app build level.
type BuildMeta struct {
	Version string
	Commit  string
	Date    string
}

// NewBuildMeta returns an instance of BuildMeta populated with build metadata.
func NewBuildMeta(version string, commit string, date string) *BuildMeta {
	return &BuildMeta{
		Version: version,
		Commit:  commit,
		Date:    date,
	}
}
