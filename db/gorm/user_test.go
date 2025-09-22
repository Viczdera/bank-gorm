package gorm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Viczdera/bank-gorm/util"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createTestUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	args := CreateUserParams{
		Username:       util.RandomOwner(),
		PasswordHashed: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testStore.CreateUser(context.Background(), args)

	fmt.Println("user: ", user)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, args.Username, user.Username)
	require.Equal(t, args.PasswordHashed, user.PasswordHashed)
	require.Equal(t, args.FullName, user.FullName)
	require.Equal(t, args.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createTestUser(t)
	userFind, err := testStore.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, userFind)

	require.Equal(t, user1.Username, userFind.Username)
	require.Equal(t, user1.PasswordHashed, userFind.PasswordHashed)
	require.Equal(t, user1.FullName, userFind.FullName)
	require.Equal(t, user1.Email, userFind.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, userFind.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, userFind.CreatedAt, time.Second)
}

func TestGetUserNotFound(t *testing.T) {
	_, err := testStore.GetUser(context.Background(), "nonexistent")
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}
