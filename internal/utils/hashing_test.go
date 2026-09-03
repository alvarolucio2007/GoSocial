package utils

import (
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/go-openapi/testify/v2/require"
)

type FakePassword struct {
	Password string `faker:"word"`
}

func TestPassword(t *testing.T) {
	fakePassword := FakePassword{}
	err := faker.FakeData(&fakePassword)
	require.NoError(t, err)
	hashedPassword, err := HashPassword(fakePassword.Password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isSame, err := CheckPassword(fakePassword.Password, hashedPassword)
	require.NoError(t, err)
	require.True(t, isSame)

	hashedPassword2, err := HashPassword(fakePassword.Password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword, hashedPassword2)
}
