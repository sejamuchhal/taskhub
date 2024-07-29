package util_test

import (
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/sejamuchhal/taskhub/auth/storage"
	"github.com/sejamuchhal/taskhub/auth/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	handler := util.NewTokenHandler("secret", db)

	user := &storage.User{
		ID: "test-user-id",
		Name: "test-user",
		Email: "user@mail.com",
	}
	tokenString, claims, err := handler.CreateToken(user, time.Hour, "access")
	require.NoError(t, err)
	require.NotNil(t, tokenString)

	mock.ExpectGet(tokenString).SetVal("")

	verifiedClaims, err := handler.VerifyToken(tokenString, "access")
	require.NoError(t, err)
	require.NotNil(t, verifiedClaims)

	assert.Equal(t, claims.UserID, verifiedClaims.UserID)
	assert.Equal(t, claims.TokenType, verifiedClaims.TokenType)
}

func TestVerifyBlacklistedToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()
	handler := util.NewTokenHandler("secret", db)

	user := &storage.User{
		ID: "test-user-id",
		Name: "test-user",
		Email: "user@mail.com",
	}

	tokenString, _, err := handler.CreateToken(user, 24*time.Hour, "refresh")
	require.NoError(t, err)
	require.NotNil(t, tokenString)

	mock.ExpectGet(tokenString).SetVal("blacklisted")

	_, err = handler.VerifyToken(tokenString, "refresh")
	require.Error(t, err)
	assert.Equal(t, "token is blacklisted", err.Error())

}

func TestBlacklistToken(t *testing.T) {
	db, mock := redismock.NewClientMock()
	defer db.Close()

	handler := util.NewTokenHandler("secret", db)
	user := &storage.User{
		ID: "test-user-id",
		Name: "test-user",
		Email: "user@mail.com",
	}

	
	tokenString, _, err := handler.CreateToken(user, 24*time.Hour, "refresh")
	require.NoError(t, err)
	require.NotNil(t, tokenString)

	mock.ExpectSet(tokenString, "blacklisted", time.Hour).SetVal("OK")
	err = handler.BlacklistToken(tokenString, time.Hour)
	require.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())

}