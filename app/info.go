package app

// This should align with https://schema.org/SoftwareApplication
type Info struct {
	Name string

	BuildInfo BuildInfo
}

// HeaderString returns a line which contains build information about
// the app that can be used to provide such info in logs.
func (info Info) HeaderString() string {
	return info.Name +
		" revision " + info.BuildInfo.RevisionID +
		" built at " + info.BuildInfo.Timestamp
}
