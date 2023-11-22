package modcheck

// VersionKey is the key for a version of a module.
type VersionKey struct {
	System  string `json:"system"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Version is a version of a module.
type Version struct {
	VersionKey VersionKey `json:"versionKey"`
	IsDefault  bool       `json:"isDefault"`
}
