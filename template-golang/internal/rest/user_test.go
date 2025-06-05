package rest_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"{{ package_name }}/domain"
	"{{ package_name }}/internal/rest"
	"{{ package_name }}/internal/rest/mocks"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserHappyPath(t *testing.T) {
	t.Parallel()

	mockUserService := new(mocks.UserService)

    newUser := domain.User{
	    ID:    "d4b8583d-5038-4838-bcd7-3d8dddfedd6a",
	    Name:  "John Doe",
	    Email: "john@example.com",
    }

    // data for update
	updatedUser := newUser
	updatedUser.Name = "Jane Doe"
	updatedUser.Email = "jane@example.com"


	handler := rest.UserHandler{
		Service: mockUserService,
	}

    // --- Create User
	t.Run("CreateUser", func(t *testing.T) {
		createReq := domain.CreateUserRequest{
			Name: newUser.Name,
			Email: newUser.Email,
            Password: "Password1234",
		}
		mockUserService.
			On("CreateUser", mock.Anything, &createReq).
			Return(&newUser, nil).
			Once()

		body, err := json.Marshal(createReq)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err = handler.CreateUser(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp domain.ResponseSingleData[domain.User]
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, &newUser, &resp.Data)

		mockUserService.AssertExpectations(t)
	})

    // --- Update User
    t.Run("UpdateUser", func(t *testing.T) {
        id, err := uuid.Parse(newUser.ID)
		require.NoError(t, err)

		mockUserService.
			On("UpdateUser", mock.Anything, id, mock.MatchedBy(func(u *domain.User) bool {
				return u.Name == updatedUser.Name && u.Email == updatedUser.Email
			})).
			Return(&updatedUser, nil).
			Once()

		body, err := json.Marshal(updatedUser)
		require.NoError(t, err)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+newUser.ID, bytes.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(newUser.ID)

		err = handler.UpdateUser(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.ResponseSingleData[domain.User]
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, &updatedUser, &resp.Data)

		mockUserService.AssertExpectations(t)
	})

// 	// --- Get User
    t.Run("GetUser", func(t *testing.T) {
        id, err := uuid.Parse(newUser.ID)
		require.NoError(t, err)

		mockUserService.
			On("GetUser", mock.Anything, id).
			Return(&updatedUser, nil).
			Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+newUser.ID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(newUser.ID)

        err = handler.GetUser(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp domain.ResponseSingleData[domain.User]
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.Equal(t, "success", resp.Status)
		assert.Equal(t, updatedUser, resp.Data)

		mockUserService.AssertExpectations(t)
	})

   // --- Delete User
	t.Run("DeleteUser", func(t *testing.T) {
        id, err := uuid.Parse(newUser.ID)
		require.NoError(t, err)
		mockUserService.
			On("DeleteUser", mock.Anything, id).
			Return(nil).
			Once()

		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+newUser.ID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(newUser.ID)

		err = handler.DeleteUser(c)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockUserService.AssertExpectations(t)
	})

    // --- Verify Deletion
	t.Run("GetUserAfterDeletion", func(t *testing.T) {
        id, err := uuid.Parse(newUser.ID)
		require.NoError(t, err)
		mockUserService.
            On("GetUser", mock.Anything, id).
            Return(nil, sql.ErrNoRows).
            Once()

		e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+newUser.ID, nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newUser.ID)

        err = handler.GetUser(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusNotFound, rec.Code)

        var resp domain.ResponseSingleData[domain.Empty]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)

        assert.Equal(t, "error", resp.Status)
        assert.Equal(t, "User not found", resp.Message)
        assert.Equal(t, http.StatusNotFound, resp.Code)

        mockUserService.AssertExpectations(t)
	})
}



func TestUserUnhappyPath(t *testing.T) {
    mockUserService := new(mocks.UserService)

    newUser := domain.User{
	    ID:    "d4b8583d-5038-4838-bcd7-3d8dddfedd6a",
	    Name:  "John Doe",
	    Email: "john@example.com",
    }

    handler := rest.UserHandler{
		Service: mockUserService,
	}


	// --- Get Non-Existent User
	t.Run("GetNonExistingUser", func(t *testing.T) {
        id, err := uuid.Parse(newUser.ID)
		require.NoError(t, err)
		mockUserService.
            On("GetUser", mock.Anything, id).
            Return(nil, sql.ErrNoRows).
            Once()

		e := echo.New()
        req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+newUser.ID, nil)
        rec := httptest.NewRecorder()
        c := e.NewContext(req, rec)
        c.SetParamNames("id")
        c.SetParamValues(newUser.ID)

        err = handler.GetUser(c)
        require.NoError(t, err)

        assert.Equal(t, http.StatusNotFound, rec.Code)

        var resp domain.ResponseSingleData[domain.Empty]
        err = json.Unmarshal(rec.Body.Bytes(), &resp)
        require.NoError(t, err)

        assert.Equal(t, "error", resp.Status)
        assert.Equal(t, "User not found", resp.Message)
        assert.Equal(t, http.StatusNotFound, resp.Code)

        mockUserService.AssertExpectations(t)
	})

	// // --- Create Invalid JSON
	t.Run("CreateUser_InvalidNameType", func(t *testing.T) {
	    body := []byte(`{
		    "Name": 12345,
		    "Email": "test@example.com",
		    "Password": "Password1234"
	    }`)

	    e := echo.New()
	    req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(body))
	    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	    rec := httptest.NewRecorder()
	    c := e.NewContext(req, rec)

	    err := handler.CreateUser(c)
	    require.NoError(t, err)

	    assert.Equal(t, http.StatusBadRequest, rec.Code)

	    var resp domain.ResponseSingleData[domain.Empty]
	    err = json.Unmarshal(rec.Body.Bytes(), &resp)
	    require.NoError(t, err)
	    assert.Equal(t, "error", resp.Status)
	    assert.NotEmpty(t, resp.Message)
})
}
