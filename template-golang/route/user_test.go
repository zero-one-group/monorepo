package router_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-app/domain"
	router "go-app/route"
	"go-app/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func startServer() *httptest.Server {
	e := echo.New()
	e.HideBanner = true
	apiV1 := e.Group("/api/v1")

	svc := service.NewUserService()
	router.RegisterUserRoutes(apiV1, svc)
	return httptest.NewServer(e)
}

func TestUserHappyPath(t *testing.T) {
	t.Parallel()
	ts := startServer()
	defer ts.Close()

	baseURL := ts.URL + "/api/v1"

	// --- Create User
	var created domain.User
	{
		payload := domain.User{
			Name:  "John Doe",
			Email: "john@example.com",
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
		if created.ID == 0 ||
			created.Name != payload.Name ||
			created.Email != payload.Email {
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
		url := fmt.Sprintf("%s/users/%d", baseURL, created.ID)
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

		var updated domain.User
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			t.Fatalf("decode update response: %v", err)
		}
		if updated.ID != created.ID ||
			updated.Name != update.Name ||
			updated.Email != update.Email {
			t.Errorf("unexpected updated user: %+v", updated)
		}
	}

	// --- Get User
	{
		url := fmt.Sprintf("%s/users/%d", baseURL, created.ID)
		resp, err := http.Get(url)
		if err != nil {
			t.Fatalf("get request error: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200; got %d", resp.StatusCode)
		}

		var fetched domain.User
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decode get response: %v", err)
		}
		if fetched.ID != created.ID ||
			fetched.Name != "Jane Doe" ||
			fetched.Email != "jane@example.com" {
			t.Errorf("unexpected fetched user: %+v", fetched)
		}
	}

	// --- Delete User
	{
		url := fmt.Sprintf("%s/users/%d", baseURL, created.ID)
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
		url := fmt.Sprintf("%s/users/%d", baseURL, created.ID)
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
	// NOTE: this will make tests run parallely
	t.Parallel()
	ts := startServer()
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
