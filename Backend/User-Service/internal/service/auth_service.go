package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kidkon/ecommerce/common/response"
	"github.com/kidkon/ecommerce/user-service/internal/model"
	"github.com/kidkon/ecommerce/user-service/internal/repository"
	"github.com/kidkon/ecommerce/user-service/internal/security"
	"github.com/kidkon/ecommerce/user-service/internal/token"
)

type AuthService struct {
	repo      repository.UserRepository
	jwtSecret string
	tokenTTL  time.Duration
}

func NewAuthService(repo repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: jwtSecret, tokenTTL: 1 * time.Hour}
}

type RegisterInput struct {
	Email, Username, Password, FullName string
}

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*model.User, error) {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.FullName = strings.TrimSpace(in.FullName)
	if in.Email == "" || in.Password == "" || in.FullName == "" {
		return nil, response.BadRequest("email, password, full_name are required")
	}

	hash, err := security.HashPassword(in.Password)
	if err != nil {
		return nil, response.Internal("failed to hash password").Wrap(err)
	}

	u := &model.User{Email: in.Email, FullName: in.FullName, Password: hash}
	if in.Username != "" {
		u.Username = &in.Username
	}

	if err := s.repo.Create(ctx, u); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, response.Conflict("email or username already exists")
		}
		return nil, response.Internal("failed to create user").Wrap(err)
	}
	return u, nil
}

type LoginResult struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

// Login: ออก JWT ที่ฝัง username/full_name/email + provider/email_verified
func (s *AuthService) Login(ctx context.Context, login, password string) (*LoginResult, error) {
	login = strings.TrimSpace(login)
	if login == "" || password == "" {
		return nil, response.BadRequest("login and password are required")
	}

	u, err := s.repo.FindByLogin(ctx, login)
	if err != nil || !security.CheckPassword(u.Password, password) {
		return nil, response.Unauthorized("invalid credentials")
	}

	username := ""
	if u.Username != nil {
		username = *u.Username
	}
	tok, err := token.Generate(token.Claims{
		UserID: u.ID, Username: username, FullName: u.FullName, Email: u.Email, Role: u.Role,
		Provider: u.Provider, EmailVerified: u.EmailVerified,
	}, s.jwtSecret, s.tokenTTL)
	if err != nil {
		return nil, response.Internal("failed to sign token").Wrap(err)
	}
	return &LoginResult{Token: tok, User: u}, nil
}

// GetProfile: ดึงข้อมูล user ละเอียดจาก DB สด (ด้วย userId จาก token)
func (s *AuthService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	u, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, response.NotFound("user not found")
		}
		return nil, response.Internal("failed to load profile").Wrap(err)
	}
	return u, nil
}
