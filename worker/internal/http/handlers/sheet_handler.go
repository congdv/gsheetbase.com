package handlers

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gsheetbase/shared/repository"

	"encoding/json"
	"sort"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetHandler struct {
	sheetRepo repository.AllowedSheetRepo
	userRepo  repository.UserRepo
}

func NewSheetHandler(sheetRepo repository.AllowedSheetRepo, userRepo repository.UserRepo) *SheetHandler {
	return &SheetHandler{
		sheetRepo: sheetRepo,
		userRepo:  userRepo,
	}
}

// isMethodAllowed checks if a specific HTTP method is allowed for the sheet
func isMethodAllowed(allowedMethods []string, method string) bool {
	for _, m := range allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// checks if the JSON has the required headers and returns the row slice
func validateAndMap(headers []interface{}, jsonInput map[string]interface{}) ([]interface{}, error) {
	newRow := make([]interface{}, len(headers))
	missingHeaders := []string{}

	for i, h := range headers {
		headerStr, ok := h.(string)
		if !ok {
			return nil, fmt.Errorf("header at column %d is not string", i)
		}

		if val, exists := jsonInput[headerStr]; exists {
			newRow[i] = val
		} else {
			missingHeaders = append(missingHeaders, headerStr)
			newRow[i] = nil
		}
	}

	if len(missingHeaders) > 0 {
		return nil, fmt.Errorf("JSON is missing required fields: %s", strings.Join(missingHeaders, ", "))
	}

	return newRow, nil
}

// Helper: filterRows applies simple equality filter from where JSON
func filterRows(rows []map[string]interface{}, where string) []map[string]interface{} {
	if where == "" {
		return rows
	}
	// Parse where as map[string]interface{}
	var cond map[string]interface{}
	err := json.Unmarshal([]byte(where), &cond)
	if err != nil {
		return rows // ignore invalid filter
	}
	filtered := make([]map[string]interface{}, 0)
	for _, row := range rows {
		match := true
		for k, v := range cond {
			if fmt.Sprintf("%v", row[k]) != fmt.Sprintf("%v", v) {
				match = false
				break
			}
		}
		if match {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

// Helper: selectFields returns only requested fields
func selectFields(rows []map[string]interface{}, fields string) []map[string]interface{} {
	if fields == "" {
		return rows
	}
	keys := strings.Split(fields, ",")
	selected := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		obj := make(map[string]interface{})
		for _, k := range keys {
			k = strings.TrimSpace(k)
			obj[k] = row[k]
		}
		selected = append(selected, obj)
	}
	return selected
}

// Helper: orderRows sorts by a single field
func orderRows(rows []map[string]interface{}, orderBy string) []map[string]interface{} {
	if orderBy == "" {
		return rows
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return fmt.Sprintf("%v", rows[i][orderBy]) < fmt.Sprintf("%v", rows[j][orderBy])
	})
	return rows
}

// Helper: paginateRows returns paginated slice and pagination info
func paginateRows(rows []map[string]interface{}, limitStr, offsetStr string) ([]map[string]interface{}, map[string]interface{}) {
	limit, err1 := strconv.Atoi(limitStr)
	offset, err2 := strconv.Atoi(offsetStr)
	if err1 != nil || err2 != nil || limit <= 0 || offset < 0 {
		return rows, nil
	}
	total := len(rows)
	end := offset + limit
	if end > total {
		end = total
	}
	paginated := rows
	if offset < total {
		paginated = rows[offset:end]
	} else {
		paginated = []map[string]interface{}{}
	}
	pagination := map[string]interface{}{
		"total":      total,
		"limit":      limit,
		"offset":     offset,
		"nextOffset": end,
	}
	return paginated, pagination
}

func (h *SheetHandler) fetchSheetData(ctx context.Context, accessToken, sheetID, rangeStr string) ([][]interface{}, error) {
	// Create OAuth2 token
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	// Create OAuth2 config (scopes don't matter for just creating a client with existing token)
	config := &oauth2.Config{
		ClientID:     "dummy",
		ClientSecret: "dummy",
		Endpoint:     google.Endpoint,
	}

	client := config.Client(ctx, token)

	// Create Sheets service
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create sheets service: %w", err)
	}

	// Get spreadsheet data
	resp, err := srv.Spreadsheets.Values.Get(sheetID, rangeStr).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %w", err)
	}

	return resp.Values, nil
}

// transformToJSON converts a 2D array to JSON objects using first row as headers
func transformToJSON(data [][]interface{}) []map[string]interface{} {
	if len(data) == 0 {
		return []map[string]interface{}{}
	}

	// First row is the header
	headers := data[0]
	result := make([]map[string]interface{}, 0, len(data)-1)

	// Process remaining rows
	for i := 1; i < len(data); i++ {
		row := data[i]
		obj := make(map[string]interface{})

		for j, header := range headers {
			headerStr := fmt.Sprintf("%v", header)
			if j < len(row) {
				obj[headerStr] = row[j]
			} else {
				obj[headerStr] = nil
			}
		}

		result = append(result, obj)
	}

	return result
}

func (h *SheetHandler) appendSheetData(ctx context.Context, accessToken, sheetID, rangeStr string, data [][]interface{}) (*sheets.AppendValuesResponse, error) {
	token := &oauth2.Token{AccessToken: accessToken}
	config := &oauth2.Config{
		ClientID:     "dummy",
		ClientSecret: "dummy",
		Endpoint:     google.Endpoint,
	}
	client := config.Client(ctx, token)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to create sheets service: %w", err)
	}

	valueRange := &sheets.ValueRange{
		Values: data,
	}

	resp, err := srv.Spreadsheets.Values.Append(sheetID, rangeStr, valueRange).
		ValueInputOption("USER_ENTERED").
		InsertDataOption("INSERT_ROWS").
		IncludeValuesInResponse(true).
		Do()
	if err != nil {
		return nil, fmt.Errorf("unable to append data to sheet: %w", err)
	}

	return resp, nil
}

func (h *SheetHandler) updateSheetData(ctx context.Context, accessToken, sheetID, rangeStr string, data [][]interface{}) error {
	token := &oauth2.Token{AccessToken: accessToken}
	config := &oauth2.Config{
		ClientID:     "dummy",
		ClientSecret: "dummy",
		Endpoint:     google.Endpoint,
	}
	client := config.Client(ctx, token)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to create sheets service: %w", err)
	}

	valueRange := &sheets.ValueRange{
		Values: data,
	}

	_, err = srv.Spreadsheets.Values.Update(sheetID, rangeStr, valueRange).
		ValueInputOption("USER_ENTERED").
		Do()
	if err != nil {
		return fmt.Errorf("unable to update sheet data: %w", err)
	}

	return nil
}

func (h *SheetHandler) deleteSheetDataAtRowIndex(ctx context.Context, accessToken, spreadsheetId, rangeStr string, rowIndexToDelete int64) error {
	token := &oauth2.Token{AccessToken: accessToken}
	config := &oauth2.Config{
		ClientID:     "dummy",
		ClientSecret: "dummy",
		Endpoint:     google.Endpoint,
	}
	client := config.Client(ctx, token)

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to create sheets service: %w", err)
	}

	spreadsheet, err := srv.Spreadsheets.Get(spreadsheetId).Do()
	if err != nil {
		return err
	}
	sheetName, _, err := parseRange(rangeStr)
	if err != nil {
		return err
	}
	var sheetID int64
	found := false
	for _, s := range spreadsheet.Sheets {
		if s.Properties.Title == sheetName {
			sheetID = s.Properties.SheetId
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("sheet named '%s' not found", sheetName)
	}

	batchUpdate := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{
			{
				DeleteDimension: &sheets.DeleteDimensionRequest{
					Range: &sheets.DimensionRange{
						SheetId:    sheetID,
						Dimension:  "ROWS",
						StartIndex: rowIndexToDelete,
						EndIndex:   rowIndexToDelete + 1,
					},
				},
			},
		},
	}

	_, err = srv.Spreadsheets.BatchUpdate(spreadsheetId, batchUpdate).Do()

	if err != nil {
		return fmt.Errorf("unable to delete row data: %w", err)
	}

	return nil
}

// This helper takes "Sheet1!A5" and returns ("Sheet1", 4)
func parseRange(rangeStr string) (string, int64, error) {
	parts := strings.Split(rangeStr, "!")
	if len(parts) != 2 {
		return rangeStr, -1, nil
	}

	sheetName := parts[0]
	// Use regex to find just the numbers in "A2" or "A2:B2"
	re := regexp.MustCompile(`[0-9]+`)
	rowStr := re.FindString(parts[1])

	rowNum, err := strconv.ParseInt(rowStr, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("could not parse row number: %v", err)
	}

	// Convert to 0-indexed (Row 2 becomes Index 1)
	return sheetName, rowNum - 1, nil
}
