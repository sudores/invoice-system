package user

import (
	context "context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/sudores/invoice-system/pkg/api/auth"
	userRepo "github.com/sudores/invoice-system/pkg/repo/user"
	"golang.org/x/crypto/bcrypt"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type UsersGrpcService struct {
	UnimplementedUserServiceServer
	log  *zerolog.Logger
	repo UserRepo
	// TODO: Switch to unified auth way through the dependency inversion
	jwm *auth.JwtManager
}

type UserRepo interface {
	CreateUser(ctx context.Context, u *userRepo.User) (*userRepo.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, u *userRepo.User) (*userRepo.User, error)
	GetUserByEmail(ctx context.Context, email string) (*userRepo.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*userRepo.User, error)
}

func NewUsersGrpcService(log *zerolog.Logger, repo UserRepo, jwm *auth.JwtManager) UsersGrpcService {
	return UsersGrpcService{
		log:  log,
		repo: repo,
		jwm:  jwm,
	}
}

func (ugs UsersGrpcService) RegisterUser(ctx context.Context, req *RegisterUserRequest) (*UserResponse, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}
	created, err := ugs.repo.CreateUser(ctx, &userRepo.User{
		PasswordHash: string(passwordHash),
		Email:        req.Email,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		return &UserResponse{
			Id:      "",
			Email:   req.Email,
			Message: "Failed to register user " + err.Error(),
		}, err
	}
	return &UserResponse{
		Id:      created.Id.String(),
		Email:   req.Email,
		Message: "User successfully created!",
	}, nil
}

func (ugs UsersGrpcService) LogIn(ctx context.Context, req *LogInRequest) (*LogInResponse, error) {
	u, err := ugs.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ugs.log.Trace().Str("user_email", req.Email).Msg("Login failed for user. No password found")
		return &LogInResponse{
			Message: "Failed to get the user in database. Check you email",
		}, errors.New("Failed to get get the user in database " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		ugs.log.Trace().Str("user_email", req.Email).Msg("Login failed for user. Password mismatch")
		return &LogInResponse{
			Message: "Password does not match for user. Please try again",
		}, errors.New("Password does not match for user: " + err.Error())
	}

	ugs.log.Trace().Str("user_email", req.Email).Msg("Login successful for user")

	jwtToken, err := ugs.jwm.GenerateJwt(u.Id)
	if err != nil {
		ugs.log.Trace().Str("user_email", req.Email).Str("user_id", u.Id.String()).Msg("Login failed. Failed to generate JWT token")
		return &LogInResponse{
			Message: "Failed to generate JWT token",
		}, errors.New("Failed to generate JWT token: " + err.Error())
	}

	return &LogInResponse{
		Jwt:          jwtToken,
		RefreshToken: "", // TODO: Implement
		Message:      "Login successful. Good luck!",
	}, nil
}

func (ugs UsersGrpcService) LogOut(ctx context.Context, req *LogOutRequest) (*LogOutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogOut not implemented")
}
func (ugs UsersGrpcService) GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*UserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
