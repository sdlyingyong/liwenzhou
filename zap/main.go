package main

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

var (
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
	url_baiud   = "http://www.baidu.com"
)

func main() {
	showZapLog()
	showZapLogFile()
	showSugarLog()
	showSugarLogFile()
	showZapWithLum()
}

func showZapWithLum() {
	//准备配置参数
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	enc := zapcore.NewConsoleEncoder(cfg)
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
	core := zapcore.NewCore(enc, ws, enab)
	//显示调用者 调用堆栈信息
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel))
	sugarLogger = logger.Sugar()

	for i := 0; i < 10000; i++ {
		sugarSimpleHttpGet(url_baiud)
	}

}

func InitLogger() {
	logger, _ = zap.NewProduction()
}

func InitFileLogger() *zap.Logger {
	//准备配置参数
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	file, _ := os.Create("./zap.log")
	ws := zapcore.AddSync(file)
	enab := zapcore.DebugLevel

	//创建日志对象
	core := zapcore.NewCore(enc, ws, enab)
	logger := zap.New(core)
	return logger
}

func InitSugarLogger() *zap.SugaredLogger {
	logger, _ = zap.NewProduction()
	sugarLogger = logger.Sugar()
	return sugarLogger
}

func InitSugarFileLogger() *zap.SugaredLogger {
	//准备配置参数
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	enc := zapcore.NewConsoleEncoder(cfg)
	file, _ := os.Create("./zap.log")
	ws := zapcore.AddSync(file)
	enab := zapcore.DebugLevel

	//创建日志对象
	core := zapcore.NewCore(enc, ws, enab)
	//显示调用者 调用堆栈信息
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.DebugLevel))
	sugarLogger = logger.Sugar()

	return sugarLogger
}

func showSugarLog() {
	sugarLogger = InitSugarLogger()
	sugarSimpleHttpGet(url_baiud)
}

func showZapLogFile() {
	logger = InitFileLogger()
	simpleHttpGet(url_baiud)
}

func showSugarLogFile() {
	sugarLogger = InitSugarFileLogger()
	sugarSimpleHttpGet(url_baiud)
}

func showZapLog() {
	InitLogger()
	defer logger.Sync()
	simpleHttpGet("http://www.baidu.com")
	simpleHttpGet("http://www.google.com")
}

func sugarSimpleHttpGet(url string) {
	sugarLogger.Debugf("Trying to hit GET request for %s", url)
	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Errorf("Error fetching URL %s : Error=%s", url, err)
	} else {
		sugarLogger.Infof("Success! statuscode is %s for URL %s", resp.StatusCode, url)
		resp.Body.Close()
	}
}

func simpleHttpGet(url string) {
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("Error fetching url...",
			zap.String("url", url),
			zap.Error(err))
		return
	}
	defer resp.Body.Close()
	logger.Info("Success...",
		zap.String("statusCode", resp.Status),
		zap.String("url", url),
	)
}
