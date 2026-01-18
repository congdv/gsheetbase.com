package services

import (
	"context"
	"fmt"
	"strings"

	"gsheetbase/web/internal/config"
	"gsheetbase/web/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// SheetService provides access to Google Sheets data with application-level access control.
// Uses spreadsheets.readonly OAuth scope (clean consent screen) with permission verification.
//
// Security model:
// - Users see "View your Google Spreadsheets" consent (no scary edit/delete warnings)
// - Application-level control: users explicitly register sheets they want to access
// - Only registered sheets in allowed_sheets table can be accessed
// - Read-only enforced by OAuth scope
type SheetService interface {
	GetSheetData(ctx context.Context, userID uuid.UUID, accessToken, sheetID, sheetRange string) ([][]interface{}, error)
}

type sheetService struct {
	cfg         *config.Config
	allowedRepo repository.AllowedSheetRepo
}

func NewSheetService(cfg *config.Config, allowedRepo repository.AllowedSheetRepo) SheetService {
	return &sheetService{
		cfg:         cfg,
		allowedRepo: allowedRepo,
	}
}

func (s *sheetService) GetSheetData(ctx context.Context, userID uuid.UUID, accessToken, sheetID, sheetRange string) ([][]interface{}, error) {
	// Verify user has registered this sheet
	allowed, err := s.allowedRepo.IsAllowed(ctx, userID, sheetID)
	if err != nil {
		return nil, fmt.Errorf("failed to check sheet permission: %w", err)
	}
	if !allowed {
		return nil, fmt.Errorf("access denied: this sheet has not been registered. Please register the sheet first via POST /api/sheets/register")
	}

	// Create Google Sheets client with access token
	oauth2Token := &oauth2.Token{AccessToken: accessToken}
	oauth2Config := &oauth2.Config{
		ClientID:     s.cfg.GoogleClientID,
		ClientSecret: s.cfg.GoogleClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{sheets.SpreadsheetsReadonlyScope},
	}
	client := oauth2Config.Client(ctx, oauth2Token)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create sheets service: %w", err)
	}

	// Read sheet data
	resp, err := srv.Spreadsheets.Values.Get(sheetID, sheetRange).Do()
	if err != nil {
		return nil, handleSheetError(err)
	}

	return resp.Values, nil
}

// handleSheetError converts Google Sheets API errors into user-friendly messages
func handleSheetError(err error) error {
	if apiErr, ok := err.(*googleapi.Error); ok {
		switch apiErr.Code {
		case 403:
			if strings.Contains(strings.ToLower(apiErr.Message), "permission") {
				return fmt.Errorf("permission denied: you need to share this Google Sheet with your account that authenticated with this app. The sheet must be shared with read access")
			}
			return fmt.Errorf("access forbidden: %s", apiErr.Message)
		case 404:
			return fmt.Errorf("sheet not found: please check the sheet ID and ensure the sheet exists")
		case 401:
			return fmt.Errorf("authentication failed: please re-authenticate with Google")
		case 400:
			return fmt.Errorf("invalid request: %s", apiErr.Message)
		default:
			return fmt.Errorf("google sheets API error (code %d): %s", apiErr.Code, apiErr.Message)
		}
	}
	return fmt.Errorf("unable to retrieve data from sheet: %w", err)
}
