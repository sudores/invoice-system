package user

import (
	context "context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/sudores/invoice-system/pkg/api/auth"
	userRepo "github.com/sudores/invoice-system/pkg/repo/user"
	"golang.org/x/crypto/bcrypt"
	grpc "google.golang.org/grpc"
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

func (ugs UsersGrpcService) Descriptor() *grpc.ServiceDesc {
	return &UserService_ServiceDesc
}

func (ugs *UsersGrpcService) RegisterHttp(ctx context.Context, mux *runtime.ServeMux) error {
	return RegisterUserServiceHandlerServer(ctx, mux, ugs)
}

type UserRepo interface {
	CreateUser(ctx context.Context, u *userRepo.User) (*userRepo.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, u *userRepo.User) (*userRepo.User, error)
	GetUserByEmail(ctx context.Context, email string) (*userRepo.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (*userRepo.User, error)
}

func NewUsersGrpcService(log *zerolog.Logger, repo UserRepo, jwm *auth.JwtManager) *UsersGrpcService {
	return &UsersGrpcService{
		log:  log,
		repo: repo,
		jwm:  jwm,
	}
}

func (ugs UsersGrpcService) Signup(ctx context.Context, req *SignupReq) (*UserResp, error) {
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
		return nil, err
	}
	return &UserResp{
		Id:      created.Id.String(),
		Email:   req.Email,
		Message: "User successfully created!",
	}, nil
}

func (ugs UsersGrpcService) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	u, err := ugs.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ugs.log.Trace().Str("user_email", req.Email).Msg("Login failed for user. No password found")
		return nil, errors.New("Failed to get get the user in database " + err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		ugs.log.Trace().Str("user_email", req.Email).
			Msg("Login failed for user. Password mismatch")
		return nil, errors.New("Password does not match for user: " + err.Error())
	}

	ugs.log.Trace().Str("user_email", req.Email).Msg("Login successful for user")

	jwtToken, err := ugs.jwm.GenerateJwt(u.Id)
	if err != nil {
		ugs.log.Trace().Str("user_email", req.Email).
			Str("user_id", u.Id.String()).Msg("Login failed. Failed to generate JWT token")
		return nil, errors.New("Failed to generate JWT token: " + err.Error())
	}

	refreshToken, err := ugs.jwm.GenerateRefresh(u.Id)
	if err != nil {
		ugs.log.Trace().Str("user_email", req.Email).
			Str("user_id", u.Id.String()).Msg("Login failed. Failed to generate refresh token")
		return nil, errors.New("Failed to generate refresh token: " + err.Error())
	}

	return &LoginResp{
		Jwt:          jwtToken,
		RefreshToken: refreshToken,
		Message:      "Login successful. Good luck!",
	}, nil
}

func (ugs UsersGrpcService) Refresh(ctx context.Context, req *RefreshReq) (*RefreshResp, error) {
	claims, err := ugs.jwm.VerifyAndParse(req.RefreshToken)
	if err != nil {
		return nil, err
	}
	strUid, ok := (*claims)["sub"].(string)
	if !ok {
		return nil, errors.New("Invalid JWT: sub claim missing or not a string")
	}
	uid, err := uuid.Parse(strUid)
	if err != nil {
		return nil, err
	}
	ugs.log.Trace().Str("user_id", uid.String()).Msg("Login successful for user")

	jwtToken, err := ugs.jwm.GenerateJwt(uid)
	if err != nil {
		ugs.log.Trace().Str("user_id", uid.String()).
			Msg("Refresh failed. Failed to generate JWT token")
		return nil, errors.New("Failed to generate JWT token: " + err.Error())
	}

	return &RefreshResp{
		Jwt: jwtToken,
	}, nil
}

func (ugs UsersGrpcService) GetUserInfo(ctx context.Context, req *GetUserInfoReq) (*UserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}

func (ugs UsersGrpcService) GetSelfInfo(ctx context.Context, req *GetSelfInfoReq) (*UserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSelfInfo not implemented")
}
