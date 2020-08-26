package cli

// OctopusServerSettings :
type OctopusServerSettings struct {
	BaseURL string
	APIKey  string
}

// LocalCacheSettings :
type LocalCacheSettings struct {
	FilePath   string
	TTLMinutes int
}

// Config is a struct
type Config struct {
	OctopusServer OctopusServerSettings
	LocalCache    LocalCacheSettings
}
