package service

import (
	"context"
	"math/rand"
	"regexp"
	"time"
	"user-service/internal/pkg/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	default:
		return nil
	}
}

func getMetadataToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md["token"]
	if len(token) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("token"), token[0])

	return ctx, nil
}

func getMetadata(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	userID := md["user_id"]
	if len(userID) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "user_id is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("userID"), userID[0])

	companyID := md["company_id"]
	if len(companyID) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "company_id is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("companyID"), companyID[0])

	return ctx, nil
}

func generateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 10
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)

	return str
}

func checkStrongPassword(password string) error {
	if len(password) < 10 {
		return status.Error(codes.InvalidArgument, "password min 10 character")
	}

	num := `[0-9]{1}`
	az := `[a-z]{1}`
	AZ := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`

	if b, err := regexp.MatchString(num, password); !b || err != nil {
		return status.Errorf(codes.InvalidArgument, "password need num :%v", err)
	}
	if b, err := regexp.MatchString(az, password); !b || err != nil {
		return status.Errorf(codes.InvalidArgument, "password need a_z :%v", err)
	}
	if b, err := regexp.MatchString(AZ, password); !b || err != nil {
		return status.Errorf(codes.InvalidArgument, "password need A_Z :%v", err)
	}
	if b, err := regexp.MatchString(symbol, password); !b || err != nil {
		return status.Errorf(codes.InvalidArgument, "password need symbol :%v", err)
	}

	return nil
}
