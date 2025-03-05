package version

import (
	"bytes"
	"fmt"
)

// Info Information about the binary version
type Info struct {
	Revision          string `json:"revision,omitempty"`
	Version           string `json:"version,omitempty"`
	VersionPrerelease string `json:"version_prerelease,omitempty"`
	VersionMetadata   string `json:"version_metadata,omitempty"`
	BuildDate         string `json:"build_date,omitempty"`
}

// GetVersion returns the information about the binary version
func GetVersion() *Info {
	ver := Version
	rel := VersionPrerelease
	md := VersionMetadata
	if GitDescribe != "" {
		ver = GitDescribe
	}
	if GitDescribe == "" && rel == "" && VersionPrerelease != "" {
		rel = "dev"
	}

	return &Info{
		Revision:          GitCommit,
		Version:           ver,
		VersionPrerelease: rel,
		VersionMetadata:   md,
		BuildDate:         BuildDate,
	}
}

// VersionNumber return the version number
func (c *Info) VersionNumber() string {
	if Version == "unknown" && VersionPrerelease == "unknown" {
		return "(version unknown)"
	}

	version := c.Version

	if c.VersionPrerelease != "" {
		version = fmt.Sprintf("%s-%s", version, c.VersionPrerelease)
	}

	if c.VersionMetadata != "" {
		version = fmt.Sprintf("%s+%s", version, c.VersionMetadata)
	}

	return version
}

// FullVersionNumber return the version number with metadata and build date
func (c *Info) FullVersionNumber(rev bool) string {
	var versionString bytes.Buffer

	if Version == "unknown" && VersionPrerelease == "unknown" {
		return "Oversee (version unknown)"
	}

	_, _ = fmt.Fprintf(&versionString, "Oversee v%s", c.Version)
	if c.VersionPrerelease != "" {
		_, _ = fmt.Fprintf(&versionString, "-%s", c.VersionPrerelease)
	}

	if c.VersionMetadata != "" {
		_, _ = fmt.Fprintf(&versionString, "+%s", c.VersionMetadata)
	}

	if rev && c.Revision != "" {
		_, _ = fmt.Fprintf(&versionString, " (%s)", c.Revision)
	}

	if c.BuildDate != "" {
		_, _ = fmt.Fprintf(&versionString, ", built %s", c.BuildDate)
	}

	return versionString.String()
}
