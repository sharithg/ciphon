package auth

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/token"
	"github.com/sharithg/siphon/internal/storage/minio"
)

type Auth struct {
	Service *auth.Service
}

type SlogLogger struct {
	slogger *slog.Logger
}

func (l SlogLogger) Info(msg string, args ...any) {
	l.slogger.Info(msg, args...)
}

func (l SlogLogger) Logf(format string, args ...any) {
	l.slogger.Info(fmt.Sprintf(format, args...))
}

func New(clientId, clientSecret string, store *minio.Storage, addr string) *Auth {
	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) {
			return "secret", nil
		}),
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		Issuer:         "ciphon",
		URL:            fmt.Sprintf("http://127.0.0.1%s", addr),
		AvatarStore:    NewFileStore("/tmp/avatars", store),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
		}),
		Logger: SlogLogger{slogger: slogger},
	}

	service := auth.NewService(options)
	service.AddProvider("github", clientId, clientSecret)

	return &Auth{
		Service: service,
	}
}
