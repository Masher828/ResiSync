package constants

import "fmt"

// consul
const (
	ConsulHost      = "CONSUL_CONFIG_HOST"
	ConsulHttpToken = "CONSUL_HTTP_TOKEN"
	ConsulConfigKey = "CONSUL_CONFIG_KEY"
)

// common
const (
	EnvEnvironment         = "ENV"
	CommonConfigFolderName = "common"
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
	EnabledKey             = "enabled"
)

var (
	ConfigSectionSmtp = "smtp"
	ConfigSmtpKey     = fmt.Sprintf("%s.%s", CommonConfigFolderName, ConfigSectionSmtp)
)

// logging
var (
	ConfigSectionLogging = "logging"
	// ConfigLogToFileKey     = "log_in_file"
	// ConfigPathToLogFileKey = "path_to_log_file"
	// ConfigLogMaxSizeKey    = "max_size_mb"
	// ConfigLogMaxBackupKey  = "max_backups"
	// ConfigLogMaxAgeDaysKey = "max_age_days"
	SendErrorEmailKey     = fmt.Sprintf("%s.%s.%s", CommonConfigFolderName, ConfigSectionLogging, "send_error_email")
	SendErrorEmailToKey   = fmt.Sprintf("%s.%s.%s", CommonConfigFolderName, ConfigSectionLogging, "send_error_email_to")
	SendErrorEmailFromKey = fmt.Sprintf("%s.%s.%s", CommonConfigFolderName, ConfigSectionLogging, "send_error_email_from")
)
