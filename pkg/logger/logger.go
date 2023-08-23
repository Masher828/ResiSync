package logger

import (
	"ResiSync/pkg/constants"
	"ResiSync/pkg/models"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func GetAppStartupLogger() *zap.Logger {
	zapConfig := zap.NewDevelopmentConfig()

	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	appStartupLogger, err := zapConfig.Build()

	if err != nil {
		log.Panic("Error while creating app startup logger ", err)
	}

	return appStartupLogger
}

func GetBasicLogger() *zap.Logger {
	return logger
}

func InitializeLoggerWithHook(hookFn func(entry zapcore.Entry) error) {

	initializeLogger()

	logger = logger.WithOptions(zap.Hooks(hookFn))
}

func initializeLogger() {
	keyPath := []string{constants.CommonConfigFolderName, constants.ConfigSectionLogging}

	key := strings.Join(keyPath, ".")

	var conf models.LoggerConfig

	err := viper.UnmarshalKey(key, &conf)
	if err != nil {
		log.Panic("Error while creating app startup logger ", err)
	}

	validate := validator.New()

	err = validate.Struct(conf)
	if err != nil {
		log.Panic("Error while creating app startup logger ", err)
	}

	level := getLevel(conf.Level)

	encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	consoleEncoderConfig := encoderConfig

	consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(log.Writer()), zap.NewAtomicLevelAt(level))

	// keyPath = append(keyPath, constants.ConfigLogToFileKey)
	//to be added file logger using lumberjack logger
	logger = zap.New(zapcore.NewTee(consoleCore), zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))

}

func getLevel(l string) zapcore.Level {
	level, err := zapcore.ParseLevel(l)

	if err != nil {
		return zapcore.InfoLevel
	}
	return level
}
