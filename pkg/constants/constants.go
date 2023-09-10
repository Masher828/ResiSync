package pkg_constants

import (
	"fmt"
	"time"
)

// consul
const (
	ConsulHost      = "CONSUL_HOST"
	ConsulHttpToken = "CONSUL_HTTP_TOKEN"
	ConsulConfigKey = "CONSUL_CONFIG_KEY"
)

// common
var (
	EnvFilePathKey         = "RESI_SYNC_ENV"
	AwsEncryptionKey       = "AWS_ENCRYPTION_KEY"
	EnvEnvironment         = "ENV"
	CommonConfigFolderName = "common"
	EnvironmentProduction  = "production"
	EnvironmentDevelopment = "development"
	EnabledKey             = "enabled"
)

// mail
var (
	ConfigSmtpKey = fmt.Sprintf("%s.%s", CommonConfigFolderName, "smtp")
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
	ConfigTracingLevel          = fmt.Sprintf("%s.%s", ConfigSectionTracing, "level")
	ConfigJaegarUrlCollectorKey = fmt.Sprintf("%s.%s", ConfigSectionTracing, "collector_url")
	TracingLevelInfo            = "info"
	TracingLevelDebug           = "debug"
	TracingLevelCritical        = "critical"
)

// aws
var (
	ConfigSectionAWS = fmt.Sprintf("%s.%s", CommonConfigFolderName, "aws")
	AWSS3Bucket      = fmt.Sprintf("%s.%s", ConfigSectionAWS, "s3_bucket")
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
	AccessTokenToUserFormatKey = "acessToken:%s"
	UserToAccessTokenKey       = "user:accessToken:%d"
	RequestUserContextKey      = "user_context"
	SessionExpiryTime          = time.Hour * 24
	RequestContextKey          = "resi_sync_request_context"
)

// security
var (
	EncryptionSection = fmt.Sprintf("%s.%s", CommonConfigFolderName, "encryption")
	EncryptionKey     = fmt.Sprintf("%s.%s", EncryptionSection, "app_encryption_key")
)
