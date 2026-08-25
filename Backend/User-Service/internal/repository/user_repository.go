package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kidkon/ecommerce/user-service/internal/model"
)

var ErrNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	FindByLogin(ctx context.Context, login string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
}

type pgUserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &pgUserRepository{db: db}
}

func (r *pgUserRepository) Create(ctx context.Context, u *model.User) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO users (email, username, password, full_name, provider)
		 VALUES ($1, $2, $3, $4, 'local')
		 RETURNING id::text, role, status, active_flag, email_verified, provider, created_at, updated_at`,
		u.Email, u.Username, u.Password, u.FullName,
	).Scan(&u.ID, &u.Role, &u.Status, &u.ActiveFlag, &u.EmailVerified, &u.Provider, &u.CreatedAt, &u.UpdatedAt)
}

// FindByLogin: ใช้ตอน login — ต้องมี password (เทียบ) + provider/email_verified (ฝังใน token)
func (r *pgUserRepository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id::text, email, username, password, full_name, role, status, active_flag, provider, email_verified
		 FROM users
		 WHERE (lower(email) = lower($1) OR username = $1) AND active_flag = true
		 LIMIT 1`,
		login,
	).Scan(&u.ID, &u.Email, &u.Username, &u.Password, &u.FullName, &u.Role, &u.Status, &u.ActiveFlag, &u.Provider, &u.EmailVerified)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// FindByID: ใช้ตอนดึง profile ละเอียดจาก DB สด (ไม่เอา password)
func (r *pgUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id::text, email, username, full_name, role, status, active_flag, provider, email_verified, created_at, updated_at
		 FROM users
		 WHERE id = $1::uuid AND active_flag = true`,
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.FullName, &u.Role, &u.Status, &u.ActiveFlag, &u.Provider, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}
