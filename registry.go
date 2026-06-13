package udb_plugin_library

var registry []UdbPlugin

// Register adds a plugin to the global registry. Call from an init() function
// so that a blank import of your plugin package triggers registration automatically.
func Register(p UdbPlugin) {
	registry = append(registry, p)
}

// Registered returns all plugins added via Register. Called by udb-core at startup
// after all init() functions have run.
func Registered() []UdbPlugin {
	out := make([]UdbPlugin, len(registry))
	copy(out, registry)
	return out
}
