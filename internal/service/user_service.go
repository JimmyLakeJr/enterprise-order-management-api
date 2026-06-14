package service

import (
	"context"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/repository"

	"github.com/jackc/pgx/v5"
)

type UserService interface {
	Me(ctx context.Context, userID int64) (*dto.UserResponse, error)
	List(ctx context.Context) ([]dto.UserResponse, error)
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	users repository.UserRepository
}

func NewUserService(users repository.UserRepository) UserService {
	return &userService{users: users}
}

func (s *userService) Me(ctx context.Context, userID int64) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}
	res := ToUserResponse(user)
	return &res, nil
}

func (s *userService) List(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.users.List(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		res = append(res, ToUserResponse(&users[i]))
	}
	return res, nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	if err := s.users.SoftDelete(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return apperror.NotFound("User not found")
		}
		return err
	}
	return nil
}
