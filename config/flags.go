package config

import "flag"

// These flags represent all the available command line flags
var (
	ConfigFile          = flag.String("config-file", "", "Location of the config file")
	ConfigFileEncrypted = flag.Bool("config-file-encrypted", false, "Force config file to be encrypted even if the file name doesn't match the expected pattern")
	CreateAccount       = flag.Bool("create", false, "If true, attempt to create account")
	DebugFlag           = flag.Bool("debug", false, "Enable debug logging")
	TraceFlag           = flag.Bool("trace", false, "Enable trace logging")
	LogFile             = flag.String("log-file", "", "Log to a file instead of to standard output")
	DebugFunctionCalls  = flag.Bool("debug-function-calls", false, "Enable logging of functions in debug logging")
	AccountFlag         = flag.String("account", "", "The account the CLI should connect to, if more than one is configured")
	MultiFlag           = flag.Bool("multi", false, "If true, will not try to unify the windows, but create separate instances")
	VersionFlag         = flag.Bool("version", false, "Print CoyIM version and exit")
	CPUProfile          = flag.String("cpuprofile", "", "write cpu profile `file`")
)
