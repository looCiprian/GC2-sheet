package authentication

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/api/drive/v3"
	"log"
)

func AuthenticateSheet(credential string) (context.Context, *sheets.Service) {


	ctx := context.Background()
	client, err := sheets.NewService(ctx, option.WithCredentialsFile(credential))
	if err != nil {
		log.Fatal("[-] Authentication failed")
	}

	return ctx, client
}

func AuthenticateDrive(credential string) (context.Context, *drive.Service) {


	ctx := context.Background()
	client, err := drive.NewService(ctx, option.WithCredentialsFile(credential))
	if err != nil {
		log.Fatal("[-] Authentication failed")
	}

	return ctx, client
}