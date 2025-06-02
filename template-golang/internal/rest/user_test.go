package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"{{ package_name }}/config"
	"{{ package_name }}/domain"
	"{{ package_name }}/internal/repository/postgres"
	"{{ package_name }}/internal/rest"
	"{{ package_name }}/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func startServer(dbPool *pgxpool.Pool) *httptest.Server {
	e := echo.New()
	e.HideBanner = true

    config.LoadEnv()

	userRepo := postgres.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepo)
	apiV1 := e.Group("/api/v1")
	rest.NewUserHandler(apiV1, userService)

	return httptest.NewServer(e)
}

func TestUserHappyPath(t *testing.T) {
	t.Parallel()

	config.LoadEnv()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Fatal("DATABASE_URL environment variable not set")
	}
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer dbPool.Close()

	ts := startServer(dbPool)
	defer ts.Close()

	baseURL := ts.URL + "/api/v1"

	// --- Create User
	var created domain.ResponseSingleData[domain.User]
	{
		payload := domain.CreateUserRequest{
			Name:  "John Doe",
			Email: "john@example.com",
            Password: "Password1234",
		}
		body, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("marshal create payload: %v", err)
		}

		resp, err := http.Post(
			baseURL+"/users",
			"application/json",
			bytes.NewBuffer(body),
		)
		if err != nil {
			t.Fatalf("create request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected status 201; got %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decode create response: %v", err)
		}

        newData := created.Data
		if newData.ID == "" ||
			newData.Name != payload.Name ||
			newData.Email != payload.Email {
			t.Errorf("unexpected created user: %+v", created)
		}
	}

	// --- Update User
	{
		update := domain.User{
			Name:  "Jane Doe",
			Email: "jane@example.com",
		}
		body, _ := json.Marshal(update)
		url := fmt.Sprintf("%s/users/%s", baseURL, created.Data.ID)
		req, err := http.NewRequest(
			http.MethodPut,
			url,
			bytes.NewBuffer(body),
		)
		if err != nil {
			t.Fatalf("new update request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("update request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200; got %d", resp.StatusCode)
		}

		var updated domain.ResponseSingleData[domain.User]
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decode update response: %v", err)
		}
		if updated.Data.ID != created.Data.ID ||
			updated.Data.Name != update.Name ||
			updated.Data.Email != update.Email {
			t.Errorf("unexpected updated user: %+v", updated)
		}
	}

	// --- Get User
	{
		url := fmt.Sprintf("%s/users/%s", baseURL, created.Data.ID)
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("get request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200; got %d", resp.StatusCode)
		}

		var fetched domain.ResponseSingleData[domain.User]
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decode get response: %v", err)
		}
		if fetched.Data.ID != created.Data.ID ||
			fetched.Data.Name != "Jane Doe" ||
			fetched.Data.Email != "jane@example.com" {
			t.Errorf("unexpected fetched user: %+v", fetched)
		}
	}

	// --- Delete User
	{
		url := fmt.Sprintf("%s/users/%s", baseURL, created.Data.ID)
		req, err := http.NewRequest(
			http.MethodDelete,
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("new delete request: %v", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("delete request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("expected status 204; got %d", resp.StatusCode)
		}
	}

	// --- Verify Deletion
	{
		url := fmt.Sprintf("%s/users/%s", baseURL, created.Data.ID)
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("get after deletion error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected status 404; got %d", resp.StatusCode)
		}
	}
}

func TestUserUnhappyPath(t *testing.T) {
	t.Parallel()

	dbURL := os.Getenv("DATABASE_URL")
	dbPool, _ := pgxpool.New(context.Background(), dbURL)
	defer dbPool.Close()

	ts := startServer(dbPool)
	defer ts.Close()

	baseURL := ts.URL

	// --- Delete Non-Existent User
	{
		url := baseURL + "/users/999"
		req, err := http.NewRequest(
			http.MethodDelete,
			url,
			nil,
		)
		if err != nil {
			t.Fatalf("new delete request: %v", err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("delete request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404; got %d", resp.StatusCode)
		}
	}

	// --- Get Non-Existent User
	{
		url := baseURL + "/users/999"
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("get request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404; got %d", resp.StatusCode)
		}
	}

	// --- Update Non-Existent User
	{
		update := domain.User{
			Name:  "Ghost",
			Email: "ghost@invalid",
		}
		body, _ := json.Marshal(update)
		url := baseURL + "/users/999"
		req, err := http.NewRequest(
			http.MethodPut,
			url,
			bytes.NewBuffer(body),
		)
		if err != nil {
			t.Fatalf("new update request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("update request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("expected 404; got %d", resp.StatusCode)
		}
	}

	// --- Create Invalid JSON
	{
		resp, err := http.Post(
			baseURL+"/users",
			"application/json",
			bytes.NewBufferString(`{"name":123}`),
		)
		if err != nil {
			t.Fatalf("create invalid JSON error: %v", err)
		}
		defer resp.Body.Close()
	}
}
