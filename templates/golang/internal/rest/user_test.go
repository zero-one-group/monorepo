package rest_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"{{ package_name }}/domain"
	"{{ package_name }}/internal/repository/postgres"
	"{{ package_name }}/internal/rest"
	"{{ package_name }}/service"

	"github.com/stretchr/testify/require"
)

func TestUserCRUD_E2E(t *testing.T) {
	kit := NewTestKit(t)

	// Wire the routes and services
	userRepo := postgres.NewUserRepository(kit.DB, kit.Metrics)
	userSvc := service.NewUserService(userRepo)
	rest.NewUserHandler(kit.Echo.Group("/api/v1"), userSvc)

	// Now start the test server
	kit.Start(t)

	// Create
	createReq := domain.CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "Password1234",
	}
	type CreateType domain.ResponseSingleData[domain.User]
	cre, code := doRequest[CreateType](
		t, http.MethodPost,
		kit.BaseURL+"/api/v1/users",
		createReq,
	)
	require.Equal(t, http.StatusCreated, code)
	user := cre.Data
	require.NotEmpty(t, user.ID)

	// Get
	type GetType domain.ResponseSingleData[domain.User]
	getE, code := doRequest[GetType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
		nil,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, user.ID, getE.Data.ID)

	// Update
	updPayload := domain.User{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}
	type UpdType domain.ResponseSingleData[domain.User]
	updE, code := doRequest[UpdType](
		t, http.MethodPut,
		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
		updPayload,
	)
	require.Equal(t, http.StatusOK, code)
	require.Equal(t, "Jane Doe", updE.Data.Name)

	// Delete
	req, err := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
		nil,
	)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Get after delete
	type ErrType domain.ResponseSingleData[domain.Empty]
	errE, code := doRequest[ErrType](
		t, http.MethodGet,
		fmt.Sprintf("%s/api/v1/users/%s", kit.BaseURL, user.ID),
		nil,
	)
	require.Equal(t, http.StatusNotFound, code)
	require.Equal(t, "User not found", errE.Message)

	// Hard delete, since delete API uses soft delete
	_, err = kit.DB.Exec(context.Background(), "DELETE from users where id = $1", user.ID)
	require.NoError(t, err)
}
