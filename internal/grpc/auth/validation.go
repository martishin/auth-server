package auth

import (
	"github.com/go-playground/validator/v10"
	ssov1 "github.com/tty-monkey/auth-server-schemas/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyAppID  = 0
	emptyUserID = 0
)

func validateLogin(validate *validator.Validate, req *ssov1.LoginRequest) error {
	if err := validateEmail(validate, req.GetEmail()); err != nil {
		return err
	}

	if err := validatePassword(validate, req.GetEmail()); err != nil {
		return err
	}

	if req.GetAppId() == emptyAppID {
		return status.Error(codes.InvalidArgument, "app id is required")
	}

	return nil
}

func validateRegister(validate *validator.Validate, req *ssov1.RegisterRequest) error {
	if err := validateEmail(validate, req.GetEmail()); err != nil {
		return err
	}

	if err := validatePassword(validate, req.GetEmail()); err != nil {
		return err
	}

	return nil
}

func validateIsAdmin(_ *validator.Validate, req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyUserID {
		return status.Error(codes.InvalidArgument, "user id is required")
	}

	return nil
}

func validateEmail(validate *validator.Validate, email string) error {
	if err := validate.Var(email, "required,email"); err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) {
			switch validationErr.Tag() {
			case "required":
				return status.Error(codes.InvalidArgument, "email is required")
			case "email":
				return status.Error(codes.InvalidArgument, "invalid email format")
			default:
				return status.Error(codes.InvalidArgument, "invalid email")
			}
		}
	}

	return nil
}

func validatePassword(validate *validator.Validate, password string) error {
	if err := validate.Var(password, "required"); err != nil {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}
