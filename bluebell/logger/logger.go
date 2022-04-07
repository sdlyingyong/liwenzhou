package logger

import (
	"bluebell/settings"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
)

func Init(cfg *settings.AppConfig, mode string) (err error) {
	//准备配置参数
	zapConf := zap.NewProductionEncoderConfig()
	zapConf.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConf.EncodeLevel = zapcore.CapitalLevelEncoder
	enc := zapcore.NewConsoleEncoder(zapConf)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./lumber.jack.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	ws := zapcore.AddSync(lumberJackLogger)

	enab := zapcore.DebugLevel

	//创建日志对象
	var core zapcore.Core
	if mode == "dev" {
		//开发模式,日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(enc, ws, enab),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewCore(enc, ws, enab)
	}
	//显示调用者 调用堆栈信息
	Logger = zap.New(core, zap.AddCaller())

	//替换zap类库中的全局logger
	zap.ReplaceGlobals(Logger)
	return
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(
			path,
			zap.Int("status", c.Writer.Status()),
			//zap.Any("header", c.Request.Header),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
