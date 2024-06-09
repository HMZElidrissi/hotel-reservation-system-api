package utils

import (
	"testing"
	"time"

	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	viper.Set("JWT_SECRET", "test_secret")
	jwtSecret = []byte(viper.GetString("JWT_SECRET"))
}

func TestGenerateJWT(t *testing.T) {
	id := primitive.NewObjectID()
	email := "test@example.com"
	role := "user"

	tokenString, err := GenerateJWT(id, email, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
}

func TestParseToken(t *testing.T) {
	id := primitive.NewObjectID()
	email := "test@example.com"
	role := "user"

	tokenString, err := GenerateJWT(id, email, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims, err := ParseToken(tokenString)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	assert.Equal(t, id.Hex(), claims["id"])
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, role, claims["role"])

	exp := int64(claims["exp"].(float64))
	assert.True(t, exp > time.Now().Unix())
}

func TestParseTokenInvalidSignature(t *testing.T) {
	id := primitive.NewObjectID()
	email := "test@example.com"
	role := "customer"

	tokenString, err := GenerateJWT(id, email, role)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	invalidTokenString := tokenString + "invalid"

	_, err = ParseToken(invalidTokenString)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, jwt.ErrSignatureInvalid))
}
