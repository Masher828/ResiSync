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

// mail
var (
	ConfigSectionSmtp = "smtp"
	ConfigSmtpKey     = fmt.Sprintf("%s.%s", CommonConfigFolderName, ConfigSectionSmtp)
)

// logging
var (
	ConfigSectionLogging  = fmt.Sprintf("%s.%s", CommonConfigFolderName, "logging")
	SendErrorEmailKey     = fmt.Sprintf("%s.%s", ConfigSectionLogging, "send_error_email")
	SendErrorEmailToKey   = fmt.Sprintf("%s.%s", ConfigSectionLogging, "send_error_email_to")
	SendErrorEmailFromKey = fmt.Sprintf("%s.%s", ConfigSectionLogging, "send_error_email_from")
)

// otel
var (
	ConfigSectionOtel           = fmt.Sprintf("%s.%s", CommonConfigFolderName, "otel")
	ConfigSectionTracing        = fmt.Sprintf("%s.%s", ConfigSectionOtel, "tracing")
	ConfigTracingEnabled        = fmt.Sprintf("%s.%s", ConfigSectionTracing, "enabled")
	ConfigJaegarUrlCollectorKey = fmt.Sprintf("%s.%s", ConfigSectionTracing, "collector_url")
)

// database
var (
	ConfigSectionDatabase = fmt.Sprintf("%s.%s", CommonConfigFolderName, "database")
	ConfigSectionPostgres = fmt.Sprintf("%s.%s", ConfigSectionDatabase, "postgres")
	ConfigSectionRedis    = fmt.Sprintf("%s.%s", ConfigSectionDatabase, "redis")
)

// api
const (
	RequestAuthenticatedKey    = "is_authenticated"
	AccessTokenToUserFormatKey = "acessToken: %s"
	RequestUserContextKey      = "user_context"
	RequestContextKey          = "resi_sync_request_context"
)

// security
var (
	EncryptionSection = fmt.Sprintf("%s.%s", CommonConfigFolderName, "encryption")
	EncryptionKey     = fmt.Sprintf("%s.%s", EncryptionSection, "app_encryption_key")
)
