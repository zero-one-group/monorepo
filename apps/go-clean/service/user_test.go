package service_test

import (
	"context"
	"errors"
	"go-clean/domain"
	"go-clean/service"
	"go-clean/service/mocks"

	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	userService := service.NewUserService(mockUserRepo)

	ctx := context.Background()
	req := &domain.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
	}
	expectedUser := &domain.User{
		ID:    uuid.New().String(),
		Name:  "Test User",
		Email: "test@example.com",
	}

	t.Run("Successfully creates a user", func(t *testing.T) {
		mockUserRepo.On("CreateUser", mock.Anything, req).Return(expectedUser, nil).Once()

		user, err := userService.CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Name, user.Name)
		assert.Equal(t, expectedUser.Email, user.Email)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		repoErr := errors.New("database error")
		mockUserRepo.On("CreateUser", mock.Anything, req).Return(nil, repoErr).Once()

		user, err := userService.CreateUser(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, repoErr, err)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userService := service.NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := &domain.User{
		ID:    userID.String(),
		Name:  "Fetched User",
		Email: "fetched@example.com",
	}

	t.Run("Successfully fetches a user", func(t *testing.T) {
		mockUserRepo.On("GetUser", mock.Anything, userID).Return(expectedUser, nil).Once()

		user, err := userService.GetUser(ctx, userID)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Name, user.Name)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		repoErr := errors.New("network error")
		mockUserRepo.On("GetUser", mock.Anything, userID).Return(nil, repoErr).Once()

		user, err := userService.GetUser(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, repoErr, err)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns nil when user not found in repository", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUser", mock.Anything, userID).Return(nil, nil).Once()

		user, err := userService.GetUser(ctx, userID)

		assert.NoError(t, err)
		assert.Nil(t, user)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userService := service.NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	existingUser := &domain.User{
		ID:    userID.String(),
		Name:  "Old Name",
		Email: "old@example.com",
	}
	updateReq := &domain.User{
		Name:  "New Name",
		Email: "new@example.com",
	}

	t.Run("Successfully updates a user", func(t *testing.T) {
		mockUserRepo.On("GetUser", mock.Anything, userID).Return(existingUser, nil).Once()

		expectedUpdatedUser := &domain.User{
			ID:    userID.String(),
			Name:  updateReq.Name,
			Email: updateReq.Email,
		}
		mockUserRepo.On("UpdateUser", mock.Anything, userID, expectedUpdatedUser).Return(expectedUpdatedUser, nil).Once()

		user, err := userService.UpdateUser(ctx, userID, updateReq)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUpdatedUser.Name, user.Name)
		assert.Equal(t, expectedUpdatedUser.Email, user.Email)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrUserNotFound if user does not exist", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUser", mock.Anything, userID).Return(nil, nil).Once()

		user, err := userService.UpdateUser(ctx, userID, updateReq)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		assert.Nil(t, user)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetUser fails", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		repoErr := errors.New("get user repo error")
		mockUserRepo.On("GetUser", mock.Anything, userID).Return(nil, repoErr).Once()

		user, err := userService.UpdateUser(ctx, userID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, repoErr, err)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns error if UpdateUser fails", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUser", mock.Anything, userID).Return(existingUser, nil).Once()

		repoErr := errors.New("update user repo error")
		expectedUpdatedUser := &domain.User{
			ID:    userID.String(),
			Name:  updateReq.Name,
			Email: updateReq.Email,
		}
		mockUserRepo.On("UpdateUser", mock.Anything, userID, expectedUpdatedUser).Return(nil, repoErr).Once()

		user, err := userService.UpdateUser(ctx, userID, updateReq)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, repoErr, err)

		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userService := service.NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	existingUser := &domain.User{
		ID:    userID.String(),
		Name:  "User to delete",
		Email: "delete@example.com",
	}

	t.Run("Successfully deletes a user", func(t *testing.T) {
		mockUserRepo.On("GetUser", mock.Anything, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("DeleteUser", mock.Anything, userID).Return(nil).Once()

		err := userService.DeleteUser(ctx, userID)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns ErrUserNotFound if user does not exist", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUser", mock.Anything, userID).Return(nil, nil).Once()

		err := userService.DeleteUser(ctx, userID)

		assert.ErrorIs(t, err, domain.ErrUserNotFound)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns error if GetUser fails", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		repoErr := errors.New("get user repo error during delete")
		mockUserRepo.On("GetUser", mock.Anything, userID).Return(nil, repoErr).Once()

		err := userService.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns error if DeleteUser fails", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUser", mock.Anything, userID).Return(existingUser, nil).Once()
		repoErr := errors.New("delete user repo error")
		mockUserRepo.On("DeleteUser", mock.Anything, userID).Return(repoErr).Once()

		err := userService.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, repoErr, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserList(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	userService := service.NewUserService(mockUserRepo)

	ctx := context.Background()
	filter := &domain.UserFilter{
		Search: "test",
	}
	expectedUsers := []domain.User{
		{ID: uuid.New().String(), Name: "Test User One"},
		{ID: uuid.New().String(), Name: "Another Test User"},
	}

	t.Run("Successfully fetches user list", func(t *testing.T) {
		mockUserRepo.On("GetUserList", mock.Anything, filter).Return(expectedUsers, nil).Once()

		users, err := userService.GetUserList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 2)
		assert.Equal(t, expectedUsers[0].Name, users[0].Name)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns empty list when no users found", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		mockUserRepo.On("GetUserList", mock.Anything, filter).Return([]domain.User{}, nil).Once()

		users, err := userService.GetUserList(ctx, filter)

		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Len(t, users, 0)

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Returns error when repository fails", func(t *testing.T) {
		mockUserRepo = new(mocks.UserRepository)
		userService = service.NewUserService(mockUserRepo)

		repoErr := errors.New("get user list database error")
		mockUserRepo.On("GetUserList", mock.Anything, filter).Return(nil, repoErr).Once()

		users, err := userService.GetUserList(ctx, filter)

		assert.Error(t, err)
		assert.Nil(t, users)
		assert.Equal(t, repoErr, err)

		mockUserRepo.AssertExpectations(t)
	})
}
