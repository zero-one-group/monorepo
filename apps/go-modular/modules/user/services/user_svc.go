package services

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
	"go-modular/modules/user/models"
	"go-modular/modules/user/repository"
)

// UserServiceInterface defines the contract for user business logic.
type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUsers(ctx context.Context, filter *models.FilterUser) ([]*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)

// UserService implements user business logic using a UserRepositoryInterface.
type UserService struct {
	userRepo repository.UserRepositoryInterface
}

type UserServiceOpts struct {
	UserRepo repository.UserRepositoryInterface
}

// NewUserService creates a new UserService.
func NewUserService(opts UserServiceOpts) *UserService {
	return &UserService{
		userRepo: opts.UserRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.Must(uuid.NewV7())
	}

	// Generate username from email if not provided
	if user.Email != "" && (user.Username == nil || *user.Username == "") {
		base := strings.SplitN(user.Email, "@", 2)[0]
		// Remove non-alphanumeric and non-underscore, make lowercase
		re := regexp.MustCompile(`[^a-z0-9_]+`)
		sanitized := re.ReplaceAllString(strings.ToLower(base), "")
		if sanitized == "" {
			sanitized = "user"
		}
		username := sanitized
		// Ensure username is unique, add suffix if needed
		suffix := 1
		for {
			exists, err := s.userRepo.UsernameExists(ctx, username)
			if err != nil {
				return err
			}
			if !exists {
				break
			}
			username = sanitized + "_" + strconv.Itoa(suffix)
			suffix++
		}
		user.Username = &username
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, filter *models.FilterUser) ([]*models.User, error) {
	return s.userRepo.ListUsers(ctx, filter)
}

func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, id)
}
