package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-mongodb-api/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "test-secret-key"

func makeToken(t *testing.T, role, email, userID string, expired bool) string {
	t.Helper()
	expiry := time.Now().Add(24 * time.Hour)
	if expired {
		expiry = time.Now().Add(-1 * time.Hour)
	}
	claims := &middleware.Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestAuthenticate_ValidToken(t *testing.T) {
	token := makeToken(t, "user", "alice@example.com", "user-id-123", false)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler))
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthenticate_MissingHeader(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler := middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler))
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_InvalidBearerFormat(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Token somevalue")
	w := httptest.NewRecorder()

	handler := middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler))
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer not-a-real-token")
	w := httptest.NewRecorder()

	handler := middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler))
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_ExpiredToken(t *testing.T) {
	token := makeToken(t, "user", "alice@example.com", "user-id-123", true)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := middleware.Authenticate(testSecret)(http.HandlerFunc(okHandler))
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthenticate_WrongSecret(t *testing.T) {
	token := makeToken(t, "user", "alice@example.com", "user-id-123", false)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := middleware.Authenticate("wrong-secret")(http.HandlerFunc(okHandler))
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRequireRoles_Allowed(t *testing.T) {
	token := makeToken(t, "admin", "alice@example.com", "user-id-123", false)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	chain := middleware.Authenticate(testSecret)(
		middleware.RequireRoles("admin", "user")(http.HandlerFunc(okHandler)),
	)
	chain.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRequireRoles_Forbidden(t *testing.T) {
	token := makeToken(t, "candidate", "bob@example.com", "user-id-456", false)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	chain := middleware.Authenticate(testSecret)(
		middleware.RequireRoles("admin")(http.HandlerFunc(okHandler)),
	)
	chain.ServeHTTP(w, r)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestRequireRoles_NoClaimsInContext(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Call RequireRoles directly without Authenticate middleware
	middleware.RequireRoles("admin")(http.HandlerFunc(okHandler)).ServeHTTP(w, r)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetClaims_Present(t *testing.T) {
	token := makeToken(t, "user", "alice@example.com", "uid-1", false)

	var capturedClaims *middleware.Claims
	captureHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r.Context())
		assert.True(t, ok)
		capturedClaims = claims
		w.WriteHeader(http.StatusOK)
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	handler := middleware.Authenticate(testSecret)(captureHandler)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, capturedClaims)
	assert.Equal(t, "alice@example.com", capturedClaims.Email)
	assert.Equal(t, "user", capturedClaims.Role)
	assert.Equal(t, "uid-1", capturedClaims.UserID)
}
