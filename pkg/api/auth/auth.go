package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const claimsKey = "jwt_claims"

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

// GenerateJwt creates a signed JWT string for a given user ID
func (jw JwtManager) GenerateJwt(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(jw.tokenTTL).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jw.jwtSecret)
}

// VerifyAndParse extracts claims from a JWT string
func (jw *JwtManager) VerifyAndParse(tokenStr string) (*jwt.MapClaims, error) {
	tokenStr = strings.TrimSpace(strings.TrimPrefix(tokenStr, "Bearer"))
	if tokenStr == "" {
		return nil, errors.New("empty token")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Ensure HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jw.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}

// UnaryInterceptor for native gRPC requests
func (jw *JwtManager) UnaryInterceptor() grpc.UnaryServerInterceptor {
	skipMethods := map[string]struct{}{
		"/user.UserService/Login":  {},
		"/user.UserService/Signup": {},
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

		claims, err := jw.VerifyAndParse(authHeader[0])
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		newCtx := context.WithValue(ctx, claimsKey, claims)
		return handler(newCtx, req)
	}
}

// GetClaimsFromContext retrieves JWT claims from context
func GetClaimsFromContext(ctx context.Context) (*jwt.MapClaims, error) {
	val := ctx.Value(claimsKey)
	if val == nil {
		return nil, status.Error(codes.Unauthenticated, "jwt claims not found in context")
	}
	claims, ok := val.(*jwt.MapClaims)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid jwt claims type")
	}
	return claims, nil
}

// GetUUIDFromContext retrieves the user ID from JWT claims in context
func GetUUIDFromContext(ctx context.Context) (*uuid.UUID, error) {
	claims, err := GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	sub, ok := (*claims)["sub"].(string)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "uuid claim missing")
	}

	uid, err := uuid.Parse(sub)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid uuid format")
	}
	return &uid, nil
}

// GatewayMiddleware returns a grpc-gateway Middleware for HTTP requests
func (jw *JwtManager) GatewayMiddleware() runtime.Middleware {
	skipPaths := map[string]struct{}{
		"/api/v1/login":  {},
		"/api/v1/signup": {},
	}

	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return runtime.HandlerFunc(func(w http.ResponseWriter, r *http.Request, p map[string]string) {
			if _, ok := skipPaths[r.URL.Path]; ok {
				next(w, r, p)
				return
			}
			fmt.Println("GatewayMiddleware here!")

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := jw.VerifyAndParse(authHeader)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next(w, r.WithContext(ctx), p)
		})
	}
}
