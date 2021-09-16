package authentication

import (
	"context"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
)

func AuthenticateSheet(credential string) (context.Context, *sheets.Service) {

	ctx := context.Background()
	client, err := sheets.NewService(ctx, option.WithCredentialsFile(credential)) // Comment this line if you have hardcoded the parameters
	// client, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(credential))) // Remove comment from this line if you have hardcoded the parameters
	if err != nil {
		log.Fatal("[-] Authentication failed")
	}

	return ctx, client
}

func AuthenticateDrive(credential string) (context.Context, *drive.Service) {


	ctx := context.Background()
	client, err := drive.NewService(ctx, option.WithCredentialsFile(credential)) // Comment this line if you have hardcoded the parameters
	// client, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(credential))) // Remove comment from this line if you have hardcoded the parameters
	if err != nil {
		log.Fatal("[-] Authentication failed")
	}

	return ctx, client
}