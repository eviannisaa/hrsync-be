package repository

import (
	"context"
	"errors"
	"hrsync-backend/internal/db"
	"hrsync-backend/internal/dto"
	"hrsync-backend/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	GetByEmail(ctx context.Context, email string) (*db.UserModel, error)
	Create(ctx context.Context, email, hashedPassword string, role db.Role, employeeId *string) (*db.UserModel, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	GeneratePassword(ctx context.Context, employeeId string) (string, error)
}

type authRepository struct {
	client *db.PrismaClient
}

func NewAuthRepository(client *db.PrismaClient) AuthRepository {
	return &authRepository{client: client}
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*db.UserModel, error) {
	user, err := r.client.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepository) Create(ctx context.Context, email, hashedPassword string, role db.Role, employeeId *string) (*db.UserModel, error) {
	optionalParams := []db.UserSetParam{
		db.User.Role.Set(role),
	}
	if employeeId != nil {
		optionalParams = append(optionalParams, db.User.EmployeeID.Set(*employeeId))
	}

	user, err := r.client.User.CreateOne(
		db.User.Email.Set(email),
		db.User.Password.Set(hashedPassword),
		optionalParams...,
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepository) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Cek apakah email sudah terdaftar
	existing, _ := r.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Parse role
	var role db.Role
	switch req.Role {
	case "ADMIN":
		role = db.RoleAdmin
	default:
		role = db.RoleEmployee
	}

	// Buat user baru
	user, err := r.Create(ctx, req.Email, string(hashed), role, req.EmployeeId)
	if err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	res := &dto.AuthResponse{
		Token: token,
		User: dto.UserInfoResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}

	return res, nil
}

func (r *authRepository) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Cari user berdasarkan email
	user, err := r.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Bandingkan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	res := &dto.AuthResponse{
		Token: token,
		User: dto.UserInfoResponse{
			ID:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
		},
	}

	return res, nil
}

func (r *authRepository) GeneratePassword(ctx context.Context, employeeId string) (string, error) {
	// 1. Fetch the employee
	emp, err := r.client.Employee.FindUnique(db.Employee.ID.Equals(employeeId)).Exec(ctx)
	if err != nil {
		return "", errors.New("employee not found")
	}

	// 2. Generate random 8-char password
	plainPassword := utils.GenerateRandomString(8)
	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// 3. Check if user exists
	existing, _ := r.GetByEmail(ctx, emp.Email)
	if existing != nil {
		// Update password
		_, err = r.client.User.FindUnique(db.User.ID.Equals(existing.ID)).Update(
			db.User.Password.Set(string(hashed)),
		).Exec(ctx)
		if err != nil {
			return "", err
		}
	} else {
		// Create new user
		_, err = r.Create(ctx, emp.Email, string(hashed), db.RoleEmployee, &emp.ID)
		if err != nil {
			return "", err
		}
	}

	return plainPassword, nil
}
