package authentication

import (
	"GC2-sheet/internal/utils"
	"context"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func AuthenticateSheet(credential string) (context.Context, *sheets.Service) {

	ctx := context.Background()
	client, err := sheets.NewService(ctx, option.WithCredentialsJSON([]byte(credential)))
	if err != nil {
		utils.LogFatalDebug("[-] Authentication failed Google Sheet")
	}

	return ctx, client
}

func AuthenticateDrive(credential string) (context.Context, *drive.Service) {

	ctx := context.Background()
	client, err := drive.NewService(ctx, option.WithCredentialsJSON([]byte(credential)))
	if err != nil {
		utils.LogFatalDebug("[-] Authentication failed Google Drive")
	}

	return ctx, client
}
