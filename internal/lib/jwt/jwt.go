package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tty-monkey/auth-server/internal/domain/models"
	"github.com/tty-monkey/auth-server/internal/lib/clock"
)

// NewToken returns a new JWT token for a passed user and app.
func NewToken(user models.User, app models.App, duration time.Duration, clock clock.Clock) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	//nolint:errcheck // linter misunderstands that error cannot be returned here
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = clock.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	return token.SignedString([]byte(app.Secret))
}
