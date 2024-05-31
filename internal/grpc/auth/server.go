package auth

import (
	"context"
	"errors"
	ssov1 "github.com/KVSH-user/protos_viz/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"vizapSSO/internal/storage"
)

const (
	emptyValue = 0
)

type Auth interface {
	Login(ctx context.Context, phone string, password string,
		appID int32) (accessToken, refreshToken string, err error)
	RegisterNewUser(ctx context.Context, phone string, password string,
	) (userID int64, err error)
	ValidateSession(ctx context.Context, accessToken string) (isValid bool, uid int64, err error)
	RefreshSession(ctx context.Context, accessToken, refreshToken string) (newAccessToken, newRefreshToken string, err error)
	RequestPasswordReset(ctx context.Context, email string) (response string, err error)
	PerformPasswordReset(ctx context.Context, token, newPassword string) (success bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.auth.Login(ctx, req.GetPhone(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.InvalidArgument, "Invalid phone or password")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetPhone(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) ValidateSession(ctx context.Context, req *ssov1.ValidateRequest,
) (*ssov1.ValidateResponse, error) {
	isValid, uid, err := s.auth.ValidateSession(ctx, req.GetAccessToken())
	if err != nil {
		//TODO: errors...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.ValidateResponse{
		IsValid: isValid,
		Uid:     uid,
	}, nil
}

func (s *serverAPI) RefreshSession(ctx context.Context, req *ssov1.RefreshRequest,
) (*ssov1.RefreshResponse, error) {
	newAccessToken, newRefreshToken, err := s.auth.RefreshSession(ctx, req.AccessToken, req.GetRefreshToken())
	if err != nil {
		if errors.Is(err, storage.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.InvalidArgument, "invalid refresh token")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RefreshResponse{
		NewAccessToken:  newAccessToken,
		NewRefreshToken: newRefreshToken,
	}, nil
}

func (s *serverAPI) RequestPasswordReset(ctx context.Context, req *ssov1.PasswordResetRequest,
) (*ssov1.PasswordResetResponse, error) {
	response, err := s.auth.RequestPasswordReset(ctx, req.GetEmail())
	if err != nil {
		//TODO: errors...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.PasswordResetResponse{
		Message: response,
	}, nil
}

func (s *serverAPI) PerformPasswordReset(ctx context.Context, req *ssov1.PerformPasswordResetRequest,
) (*ssov1.PerformPasswordResetResponse, error) {
	success, err := s.auth.PerformPasswordReset(ctx, req.GetToken(), req.GetNewPassword())
	if err != nil {
		//TODO: errors...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.PerformPasswordResetResponse{
		Success: success,
	}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.GetPhone() == "" {
		return status.Error(codes.InvalidArgument, "phone is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetPhone() == "" {
		return status.Error(codes.InvalidArgument, "phone is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
