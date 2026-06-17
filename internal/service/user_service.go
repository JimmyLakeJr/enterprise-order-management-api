package service

import (
	"context"
	"math"

	"enterprise-order-management-api/internal/dto"
	"enterprise-order-management-api/internal/pkg/apperror"
	"enterprise-order-management-api/internal/pkg/response"
	"enterprise-order-management-api/internal/repository"

	"github.com/jackc/pgx/v5"
)

type UserService interface {
	Me(ctx context.Context, userID int64) (*dto.UserResponse, error)
	List(ctx context.Context, query dto.UserListQuery) ([]dto.UserResponse, response.Meta, error)
	FindByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	Delete(ctx context.Context, id int64, currentUserID int64) error
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

func (s *userService) List(ctx context.Context, query dto.UserListQuery) ([]dto.UserResponse, response.Meta, error) {
	query = normalizeUserListQuery(query)

	users, total, err := s.users.List(ctx, query)
	if err != nil {
		return nil, response.Meta{}, err
	}
	res := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		res = append(res, ToUserResponse(&users[i]))
	}

	meta := response.Meta{
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(query.Limit))),
	}
	return res, meta, nil
}

func (s *userService) FindByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}
	res := ToUserResponse(user)
	return &res, nil
}

func (s *userService) Update(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.users.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFound("User not found")
	}

	exists, err := s.users.ExistsByEmailOtherUser(ctx, req.Email, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperror.Conflict("Email already exists")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Role = req.Role

	if err := s.users.Update(ctx, user); err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperror.NotFound("User not found")
		}
		return nil, err
	}

	return s.FindByID(ctx, id)
}

func (s *userService) Delete(ctx context.Context, id int64, currentUserID int64) error {
	if id == currentUserID {
		return apperror.BadRequest("Admin cannot delete own account")
	}

	if err := s.users.SoftDelete(ctx, id); err != nil {
		if err == pgx.ErrNoRows {
			return apperror.NotFound("User not found")
		}
		return err
	}
	return nil
}

func normalizeUserListQuery(query dto.UserListQuery) dto.UserListQuery {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
	return query
}
