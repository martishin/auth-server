package jwt_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tty-monkey/auth-server/internal/domain/models"
	"github.com/tty-monkey/auth-server/internal/lib/clock"
	jwtlib "github.com/tty-monkey/auth-server/internal/lib/jwt"
)

type MockClock struct {
	MockTime time.Time
}

func (m MockClock) Now() time.Time {
	return m.MockTime
}

func setup() (models.User, models.App, clock.Clock) {
	user := models.User{ID: 1, Email: "test@example.com"}
	app := models.App{ID: 1, Name: "test", Secret: "secret"}
	mockClock := MockClock{MockTime: time.Now()}
	return user, app, mockClock
}

func TestNewToken_CanCreate(t *testing.T) {
	// given
	user, app, mockClock := setup()

	// when
	_, err := jwtlib.NewToken(user, app, time.Hour, mockClock)

	// then
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}
}

func TestNewToken_CanParseCreatedToken(t *testing.T) {
	// given
	user, app, mockClock := setup()

	// when
	tokenString, _ := jwtlib.NewToken(user, app, time.Hour, mockClock)

	// then
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}
}

func TestNewToken_CorrectClaims(t *testing.T) {
	// given
	user, app, mockClock := setup()

	// when
	tokenString, _ := jwtlib.NewToken(user, app, time.Hour, mockClock)

	// then
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to read token claims")
	}

	if int64(claims["uid"].(float64)) != user.ID {
		t.Errorf("Expected uid %d, got %v", user.ID, claims["uid"])
	}
	if claims["email"] != user.Email {
		t.Errorf("Expected email %s, got %v", user.Email, claims["email"])
	}
	if int(claims["app_id"].(float64)) != app.ID {
		t.Errorf("Expected app_id %d, got %v", app.ID, claims["app_id"])
	}
}

func TestNewToken_CorrectExpiration(t *testing.T) {
	// given
	user, app, mockClock := setup()

	// when
	tokenString, _ := jwtlib.NewToken(user, app, time.Hour, mockClock)

	// then
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(app.Secret), nil
	})

	claims, _ := token.Claims.(jwt.MapClaims)

	expectedExp := mockClock.Now().Add(time.Hour).Unix()
	if int64(claims["exp"].(float64)) != expectedExp {
		t.Errorf("Expected exp %d, got %v", expectedExp, claims["exp"])
	}
}
