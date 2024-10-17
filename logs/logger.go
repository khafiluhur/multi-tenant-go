package logs

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func SetupLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(&lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})
	return logger
}

func LogWithFields[T any](logger *logrus.Logger, level logrus.Level, context string, data T) {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Struct {
		fields := logrus.Fields{}
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			fieldValue := val.Field(i).Interface()
			fields[field.Name] = fieldValue
		}

		switch level {
		case logrus.TraceLevel:
			logger.WithFields(fields).Trace(context)
		case logrus.DebugLevel:
			logger.WithFields(fields).Debug(context)
		case logrus.InfoLevel:
			logger.WithFields(fields).Info(context)
		case logrus.WarnLevel:
			logger.WithFields(fields).Warn(context)
		case logrus.ErrorLevel:
			logger.WithFields(fields).Error(context)
		case logrus.FatalLevel:
			logger.WithFields(fields).Fatal(context)
		case logrus.PanicLevel:
			logger.WithFields(fields).Panic(context)
		}
	} else {
		logger.WithField("data", data).Error("Provided data is not a struct")
	}
}
