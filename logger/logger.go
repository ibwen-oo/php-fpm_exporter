package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(logLevel, logPath string) (err error) {
	writer := getLogWriter(logPath)

	encoder := getEncoder()

	// 指定日志级别
	var l = new(zapcore.Level)
	// 从viper读取日志级别配置
	err = l.UnmarshalText([]byte(logLevel))
	if err != nil {
		return
	}
	var core zapcore.Core

	core = zapcore.NewCore(encoder, writer, l)

	// 输出文件名和行号
	lg := zap.New(core, zap.AddCaller())
	// 替换zap包中全局的logger实例,后续在其它包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(lg)
	return
}

// getEncoder 获取Encoder配置
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 获取日志输出配置
func getLogWriter(filename string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,
		MaxBackups: 7,
	}
	return zapcore.AddSync(lumberJackLogger)
}