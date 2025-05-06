package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/Cell6969/go_bank/db/sqlc"
	"github.com/Cell6969/go_bank/pb"
	"github.com/Cell6969/go_bank/util"
	"github.com/Cell6969/go_bank/valid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateUserParams{
		Username: request.GetUsername(),
		FullName: sql.NullString{
			String: request.GetFullName(),
			Valid:  request.FullName != nil,
		},
		Email: sql.NullString{
			String: request.GetEmail(),
			Valid:  request.Email != nil,
		},
	}

	if request.Password != nil {
		hashedPassword, err := util.HashPassword(request.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to hash password: %s", err)
		}

		arg.Password = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found:%s", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user:%s", err)
	}

	response := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return response, nil
}

func validateUpdateUserRequest(request *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := valid.ValidateUsername(request.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if request.Password != nil {
		if err := valid.ValidatePassword(request.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if request.FullName != nil {
		if err := valid.ValidateFullName(request.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if request.Email != nil {
		if err := valid.ValidateEmail(request.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
