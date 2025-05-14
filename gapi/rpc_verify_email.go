package gapi

import (
	"context"

	db "github.com/Cell6969/go_bank/db/sqlc"
	"github.com/Cell6969/go_bank/pb"
	"github.com/Cell6969/go_bank/valid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, request *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(request)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    request.GetEmailId(),
		SecretCode: request.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify email into database")
	}

	response := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}

	return response, nil
}

func validateVerifyEmailRequest(request *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := valid.ValidateEmailId(request.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}

	if err := valid.ValidateSecretCode(request.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
