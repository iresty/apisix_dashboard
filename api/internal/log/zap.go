/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package log

import (
	"net/url"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/apache/apisix-dashboard/api/internal/conf"
)

var logger *zap.SugaredLogger

// TODO: we should no longer use init() function after remove all handler's integration tests
// ENV=test is for integration tests only, other ENV should call "InitLogger" explicitly
func init() {
	if env := os.Getenv("ENV"); env == conf.EnvTEST {
		InitLogger()
	}
}

func InitLogger() {
	logger = GetLogger(ErrorLog)
}

func GetLogger(logType Type) *zap.SugaredLogger {
	_ = zap.RegisterSink("winfile", newWinFileSink)

	skip := 2
	writeSyncer := fileWriter(logType)
	encoder := getEncoder(logType)
	logLevel := getLogLevel()
	if logType == AccessLog {
		logLevel = zapcore.InfoLevel
		skip = 0
	}

	core := zapcore.NewCore(encoder, writeSyncer, logLevel)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(skip))

	return zapLogger.Sugar()
}

func getLogLevel() zapcore.LevelEnabler {
	level := zapcore.WarnLevel
	switch conf.ErrorLogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	}
	return level
}

func getEncoder(logType Type) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if logType == AccessLog {
		encoderConfig.LevelKey = zapcore.OmitKey
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func fileWriter(logType Type) zapcore.WriteSyncer {
	logPath := conf.ErrorLogPath
	if logType == AccessLog {
		logPath = conf.AccessLogPath
	}
	//standard output
	if logPath == "/dev/stdout" {
		return zapcore.Lock(os.Stdout)
	}
	if logPath == "/dev/stderr" {
		return zapcore.Lock(os.Stderr)
	}

	writer, _, err := zap.Open(logPath)
	if err != nil {
		panic(err)
	}
	return writer
}

func getZapFields(logger *zap.SugaredLogger, fields []interface{}) *zap.SugaredLogger {
	if len(fields) == 0 {
		return logger
	}

	return logger.With(fields)
}

func newWinFileSink(u *url.URL) (zap.Sink, error) {
	return os.OpenFile(u.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
}
