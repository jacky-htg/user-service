package service

import (
	"context"

	"github.com/jacky-htg/erp-pkg/app"
	"google.golang.org/grpc/metadata"
)

func setMetadata(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"user_id":    ctx.Value(app.Ctx("userID")).(string),
		"company_id": ctx.Value(app.Ctx("companyID")).(string),
	})

	return metadata.NewOutgoingContext(ctx, md)
}

func setMetadataToken(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"token": ctx.Value(app.Ctx("token")).(string),
	})

	return metadata.NewOutgoingContext(ctx, md)
}
