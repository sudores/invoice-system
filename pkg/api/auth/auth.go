package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JwtManager struct {
	jwtSecret []byte
	tokenTTL  time.Duration
}

func NewJwtManager(config Config) *JwtManager {
	return &JwtManager{
		jwtSecret: []byte(config.JwtSecret),
		tokenTTL:  config.JwtTokenTTL,
	}
}

func (jw JwtManager) GenerateJwt(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(jw.tokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jw.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (jw JwtManager) VerifyJwt(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jw.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("invalid token")
}

func (jw JwtManager) JwtOk(token string) (bool, error) {
	claims, err := jw.VerifyJwt(token)
	if err != nil {
		return false, err
	}
	if claims == nil {
		return false, nil
	}
	return true, nil
}

// UnaryInterceptor returns a gRPC unary interceptor method on JwtManager
func (jw *JwtManager) UnaryInterceptor() grpc.UnaryServerInterceptor {
	skipMethods := map[string]struct{}{
		"/user.UserService/LogIn":        {},
		"/user.UserService/RegisterUser": {},
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := skipMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata not provided")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is required")
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(authHeader[0], "Bearer"))
		if tokenStr == "" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		claims, err := jw.VerifyJwt(tokenStr)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Store claims in context for downstream handlers
		newCtx := context.WithValue(ctx, "jwt_claims", claims)

		return handler(newCtx, req)
	}
}
