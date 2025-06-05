package service

import (
	"context"
	"go-app/domain"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
)

// UserService holds users in‐memory.
type UserService struct {
	mu     sync.Mutex
	users  map[int]*domain.User
	nextID int
}

// NewUserService constructs a fresh service.
func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]*domain.User),
		nextID: 1,
	}
}

// CreateUser adds a new user.
func (s *UserService) CreateUser(
	ctx context.Context,
	u *domain.User,
) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &domain.User{
		ID:    s.nextID,
		Name:  u.Name,
		Email: u.Email,
	}
	s.users[s.nextID] = user
	s.nextID++
	return user, nil
}

// GetUser fetches a user by ID.
func (s *UserService) GetUser(
	ctx context.Context,
	id int,
) (*domain.User, error) {
	span, svcCtx := opentracing.StartSpanFromContext(
		ctx,
		"UserService.GetUser",
	)
	defer span.Finish()

	s.mu.Lock()
	defer s.mu.Unlock()

	// NOTE: Example span on repository
	spanRepo, _ := opentracing.StartSpanFromContext(
		svcCtx,
		"UserRepository.GetUser",
	)
	defer spanRepo.Finish()
	spanRepo.SetOperationName("ExampleRepoGetUserByID")
	spanRepo.SetTag("db.statement", "SELECT * FROM users WHERE id=$1")
	time.Sleep(time.Second.Abs())

	if u, ok := s.users[id]; ok {
		// return a copy
		return &domain.User{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		}, nil
	}
	spanRepo.SetTag("error", true)
	return nil, domain.ErrUserNotFound
}

// UpdateUser updates name/email of an existing user.
func (s *UserService) UpdateUser(
	ctx context.Context,
	id int,
	u *domain.User,
) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing, ok := s.users[id]; ok {
		existing.Name = u.Name
		existing.Email = u.Email
		// return a copy
		return &domain.User{
			ID:    existing.ID,
			Name:  existing.Name,
			Email: existing.Email,
		}, nil
	}
	return nil, domain.ErrUserNotFound
}

// DeleteUser removes a user by ID.
func (s *UserService) DeleteUser(
	ctx context.Context,
	id int,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; ok {
		delete(s.users, id)
		return nil
	}
	return domain.ErrUserNotFound
}
