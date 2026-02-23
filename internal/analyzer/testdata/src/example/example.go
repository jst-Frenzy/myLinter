package example

import (
	"go.uber.org/zap"
	"log/slog"
	"os"
)

func test() {
	slogLogger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	sugarZap := zapLogger.Sugar()

	password := ""
	apiKey := ""
	token := ""

	password += "password"
	apiKey += "apiKey"
	token += "token"

	slogLogger.Info("Starting server on port 8080") // want "log has issues: starting with capital letter"
	sugarZap.Error("Failed to connect to database") // want "log has issues: starting with capital letter"

	slogLogger.Info("starting server on port 8080") // +
	sugarZap.Error("failed to connect to database") // +

	slogLogger.Info("запуск сервера")                  // want "log has issues: non-English"
	sugarZap.Error("ошибка подключения к базе данных") // want "log has issues: non-English"

	slogLogger.Info("starting server")             //+
	zap.L().Error("failed to connect to database") // +

	slog.With().Info("server started!🚀")             // want "log has issues: special char or emoji"
	slogLogger.Error("connection failed!!!")         // want "log has issues: special char or emoji"
	zap.L().Warn("warning: something went wrong...") // want "log has issues: special char or emoji"

	slogLogger.Info("server started")     // +
	slogLogger.Error("connection failed") // +
	zap.L().Warn("something went wrong")  // +

	slogLogger.Info("user password: " + password) // want "log has issues: sensitive data, special char or emoji"
	slogLogger.Debug("api_key= " + apiKey)        // want "log has issues: sensitive data, special char or emoji"
	zap.L().Info("token: " + token)               // want "log has issues: sensitive data, special char or emoji"

	slogLogger.Info("user authenticated successfully") // +
	slogLogger.Debug("api request completed")          // +
	zap.L().Info("token validated")                    // +

}
