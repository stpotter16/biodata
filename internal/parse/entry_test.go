package parse

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stpotter16/biodata/internal/types"
)

func TestParseBPString(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantErr       bool
		wantSystolic  *float64
		wantDiastolic *float64
	}{
		{
			name:          "valid BP",
			input:         "120/80",
			wantErr:       false,
			wantSystolic:  floatPtr(120),
			wantDiastolic: floatPtr(80),
		},
		{
			name:          "valid BP with decimals",
			input:         "120.5/80.3",
			wantErr:       false,
			wantSystolic:  floatPtr(120.5),
			wantDiastolic: floatPtr(80.3),
		},
		{
			name:    "missing slash",
			input:   "12080",
			wantErr: true,
		},
		{
			name:    "too many slashes",
			input:   "120/80/70",
			wantErr: true,
		},
		{
			name:    "non-numeric systolic",
			input:   "abc/80",
			wantErr: true,
		},
		{
			name:    "non-numeric diastolic",
			input:   "120/xyz",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "slash only",
			input:   "/",
			wantErr: true,
		},
		{
			name:    "spaces around values",
			input:   " 120 / 80 ",
			wantErr: true, // Current implementation doesn't trim
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBPString(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("parseBPString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Systolic == nil || *got.Systolic != *tt.wantSystolic {
					t.Errorf("parseBPString() systolic = %v, want %v", got.Systolic, tt.wantSystolic)
				}
				if got.Diastolic == nil || *got.Diastolic != *tt.wantDiastolic {
					t.Errorf("parseBPString() diastolic = %v, want %v", got.Diastolic, tt.wantDiastolic)
				}
			}
		})
	}
}

func TestParseEntryPost(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		wantErr     bool
		checkResult func(t *testing.T, entry types.Entry)
	}{
		{
			name:        "all fields valid",
			requestBody: `{"date": "2024-01-15", "weight": "180.5", "waist": "34.2", "bp": "120/80"}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Date.Format("2006-01-02") != "2024-01-15" {
					t.Errorf("date = %v, want 2024-01-15", entry.Date)
				}
				if entry.Weight.Value == nil || *entry.Weight.Value != 180.5 {
					t.Errorf("weight = %v, want 180.5", entry.Weight.Value)
				}
				if entry.Waist.Value == nil || *entry.Waist.Value != 34.2 {
					t.Errorf("waist = %v, want 34.2", entry.Waist.Value)
				}
				if entry.BP.Systolic == nil || *entry.BP.Systolic != 120 {
					t.Errorf("systolic = %v, want 120", entry.BP.Systolic)
				}
				if entry.BP.Diastolic == nil || *entry.BP.Diastolic != 80 {
					t.Errorf("diastolic = %v, want 80", entry.BP.Diastolic)
				}
			},
		},
		{
			name:        "only required field (date)",
			requestBody: `{"date": "2024-01-15", "weight": "", "waist": "", "bp": ""}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Date.Format("2006-01-02") != "2024-01-15" {
					t.Errorf("date = %v, want 2024-01-15", entry.Date)
				}
				if entry.Weight.Value != nil {
					t.Errorf("weight should be nil, got %v", entry.Weight.Value)
				}
				if entry.Waist.Value != nil {
					t.Errorf("waist should be nil, got %v", entry.Waist.Value)
				}
				if entry.BP.Valid() {
					t.Errorf("BP should be invalid")
				}
			},
		},
		{
			name:        "optional weight only",
			requestBody: `{"date": "2024-01-15", "weight": "180", "waist": "", "bp": ""}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Weight.Value == nil || *entry.Weight.Value != 180 {
					t.Errorf("weight = %v, want 180", entry.Weight.Value)
				}
				if entry.Waist.Value != nil {
					t.Errorf("waist should be nil")
				}
			},
		},
		{
			name:        "optional waist only",
			requestBody: `{"date": "2024-01-15", "weight": "", "waist": "34", "bp": ""}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Waist.Value == nil || *entry.Waist.Value != 34 {
					t.Errorf("waist = %v, want 34", entry.Waist.Value)
				}
				if entry.Weight.Value != nil {
					t.Errorf("weight should be nil")
				}
			},
		},
		{
			name:        "optional BP only",
			requestBody: `{"date": "2024-01-15", "weight": "", "waist": "", "bp": "120/80"}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.BP.Systolic == nil || *entry.BP.Systolic != 120 {
					t.Errorf("systolic = %v, want 120", entry.BP.Systolic)
				}
				if entry.Weight.Value != nil || entry.Waist.Value != nil {
					t.Errorf("weight and waist should be nil")
				}
			},
		},
		{
			name:        "invalid date format",
			requestBody: `{"date": "01/15/2024", "weight": "180", "waist": "", "bp": ""}`,
			wantErr:     true,
		},
		{
			name:        "invalid weight",
			requestBody: `{"date": "2024-01-15", "weight": "abc", "waist": "", "bp": ""}`,
			wantErr:     true,
		},
		{
			name:        "invalid waist",
			requestBody: `{"date": "2024-01-15", "weight": "", "waist": "xyz", "bp": ""}`,
			wantErr:     true,
		},
		{
			name:        "invalid BP format",
			requestBody: `{"date": "2024-01-15", "weight": "", "waist": "", "bp": "120-80"}`,
			wantErr:     true,
		},
		{
			name:        "malformed JSON",
			requestBody: `{"date": "2024-01-15", "weight": "180"`,
			wantErr:     true,
		},
		{
			name:        "missing date field",
			requestBody: `{"weight": "180", "waist": "", "bp": ""}`,
			wantErr:     true,
		},
		{
			name:        "empty request body",
			requestBody: ``,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/entry", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			got, err := ParseEntryPost(req)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEntryPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkResult != nil {
				tt.checkResult(t, got)
			}
		})
	}
}

func TestParseEntryPut(t *testing.T) {
	tests := []struct {
		name        string
		pathDate    string
		requestBody string
		wantErr     bool
		checkResult func(t *testing.T, entry types.Entry)
	}{
		{
			name:        "all fields valid",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "180.5", "waist": "34.2", "bp": "120/80"}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Date.Format("2006-01-02") != "2024-01-15" {
					t.Errorf("date = %v, want 2024-01-15", entry.Date)
				}
				if entry.Weight.Value == nil || *entry.Weight.Value != 180.5 {
					t.Errorf("weight = %v, want 180.5", entry.Weight.Value)
				}
			},
		},
		{
			name:        "all fields empty (setting to null)",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "", "waist": "", "bp": ""}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Weight.Value != nil {
					t.Errorf("weight should be nil")
				}
				if entry.Waist.Value != nil {
					t.Errorf("waist should be nil")
				}
				if entry.BP.Valid() {
					t.Errorf("BP should be invalid")
				}
			},
		},
		{
			name:        "only weight",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "180", "waist": "", "bp": ""}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Weight.Value == nil || *entry.Weight.Value != 180 {
					t.Errorf("weight = %v, want 180", entry.Weight.Value)
				}
			},
		},
		{
			name:        "only waist",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "", "waist": "34", "bp": ""}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Waist.Value == nil || *entry.Waist.Value != 34 {
					t.Errorf("waist = %v, want 34", entry.Waist.Value)
				}
			},
		},
		{
			name:        "only BP",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "", "waist": "", "bp": "120/80"}`,
			wantErr:     false,
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.BP.Systolic == nil || *entry.BP.Systolic != 120 {
					t.Errorf("systolic = %v, want 120", entry.BP.Systolic)
				}
			},
		},
		{
			name:        "invalid weight",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "abc", "waist": "", "bp": ""}`,
			wantErr:     true,
		},
		{
			name:        "invalid waist",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "", "waist": "xyz", "bp": ""}`,
			wantErr:     true,
		},
		{
			name:        "invalid BP format",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "", "waist": "", "bp": "120-80"}`,
			wantErr:     true,
		},
		{
			name:        "malformed JSON",
			pathDate:    "2024-01-15",
			requestBody: `{"weight": "180"`,
			wantErr:     true,
		},
		{
			name:        "invalid date in path",
			pathDate:    "01/15/2024",
			requestBody: `{"weight": "180", "waist": "", "bp": ""}`,
			wantErr:     false, // Note: current implementation returns nil error, might want to change
		},
		{
			name:        "empty request body",
			pathDate:    "2024-01-15",
			requestBody: ``,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/entries/"+tt.pathDate, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			req.SetPathValue("date", tt.pathDate)

			got, err := ParseEntryPut(req)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseEntryPut() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.checkResult != nil {
				tt.checkResult(t, got)
			}
		})
	}
}

func TestParseEntryDTO(t *testing.T) {
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		dto         types.EntryDTO
		checkResult func(t *testing.T, entry types.Entry)
	}{
		{
			name: "all fields valid",
			dto: types.EntryDTO{
				Id:           1,
				Date:         testDate,
				Weight:       sql.NullFloat64{Float64: 180.5, Valid: true},
				Waist:        sql.NullFloat64{Float64: 34.2, Valid: true},
				Bp:           sql.NullString{String: "120/80", Valid: true},
				Created:      "2024-01-15T00:00:00Z",
				LastModified: "2024-01-15T00:00:00Z",
			},
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Weight.Value == nil || *entry.Weight.Value != 180.5 {
					t.Errorf("weight = %v, want 180.5", entry.Weight.Value)
				}
				if entry.Waist.Value == nil || *entry.Waist.Value != 34.2 {
					t.Errorf("waist = %v, want 34.2", entry.Waist.Value)
				}
				if entry.BP.Systolic == nil || *entry.BP.Systolic != 120 {
					t.Errorf("systolic = %v, want 120", entry.BP.Systolic)
				}
			},
		},
		{
			name: "all fields null",
			dto: types.EntryDTO{
				Id:           1,
				Date:         testDate,
				Weight:       sql.NullFloat64{Valid: false},
				Waist:        sql.NullFloat64{Valid: false},
				Bp:           sql.NullString{Valid: false},
				Created:      "2024-01-15T00:00:00Z",
				LastModified: "2024-01-15T00:00:00Z",
			},
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Weight.Value != nil {
					t.Errorf("weight should be nil")
				}
				if entry.Waist.Value != nil {
					t.Errorf("waist should be nil")
				}
				if entry.BP.Valid() {
					t.Errorf("BP should be invalid")
				}
			},
		},
		{
			name: "valid weight, null waist and BP",
			dto: types.EntryDTO{
				Id:           1,
				Date:         testDate,
				Weight:       sql.NullFloat64{Float64: 180, Valid: true},
				Waist:        sql.NullFloat64{Valid: false},
				Bp:           sql.NullString{Valid: false},
				Created:      "2024-01-15T00:00:00Z",
				LastModified: "2024-01-15T00:00:00Z",
			},
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Weight.Value == nil || *entry.Weight.Value != 180 {
					t.Errorf("weight = %v, want 180", entry.Weight.Value)
				}
				if entry.Waist.Value != nil {
					t.Errorf("waist should be nil")
				}
			},
		},
		{
			name: "valid waist, null weight and BP",
			dto: types.EntryDTO{
				Id:           1,
				Date:         testDate,
				Weight:       sql.NullFloat64{Valid: false},
				Waist:        sql.NullFloat64{Float64: 34, Valid: true},
				Bp:           sql.NullString{Valid: false},
				Created:      "2024-01-15T00:00:00Z",
				LastModified: "2024-01-15T00:00:00Z",
			},
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.Waist.Value == nil || *entry.Waist.Value != 34 {
					t.Errorf("waist = %v, want 34", entry.Waist.Value)
				}
				if entry.Weight.Value != nil {
					t.Errorf("weight should be nil")
				}
			},
		},
		{
			name: "valid BP, null weight and waist",
			dto: types.EntryDTO{
				Id:           1,
				Date:         testDate,
				Weight:       sql.NullFloat64{Valid: false},
				Waist:        sql.NullFloat64{Valid: false},
				Bp:           sql.NullString{String: "120/80", Valid: true},
				Created:      "2024-01-15T00:00:00Z",
				LastModified: "2024-01-15T00:00:00Z",
			},
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.BP.Systolic == nil || *entry.BP.Systolic != 120 {
					t.Errorf("systolic = %v, want 120", entry.BP.Systolic)
				}
				if entry.Weight.Value != nil || entry.Waist.Value != nil {
					t.Errorf("weight and waist should be nil")
				}
			},
		},
		{
			name: "invalid BP string in DTO (should handle gracefully)",
			dto: types.EntryDTO{
				Id:           1,
				Date:         testDate,
				Weight:       sql.NullFloat64{Valid: false},
				Waist:        sql.NullFloat64{Valid: false},
				Bp:           sql.NullString{String: "invalid", Valid: true},
				Created:      "2024-01-15T00:00:00Z",
				LastModified: "2024-01-15T00:00:00Z",
			},
			checkResult: func(t *testing.T, entry types.Entry) {
				if entry.BP.Valid() {
					t.Errorf("BP should be invalid for malformed string")
				}
			},
		},
		{
			name: "date only (minimum valid DTO)",
			dto: types.EntryDTO{
				Id:           1,
				Date:         testDate,
				Weight:       sql.NullFloat64{Valid: false},
				Waist:        sql.NullFloat64{Valid: false},
				Bp:           sql.NullString{Valid: false},
				Created:      "2024-01-15T00:00:00Z",
				LastModified: "2024-01-15T00:00:00Z",
			},
			checkResult: func(t *testing.T, entry types.Entry) {
				if !entry.Date.Equal(testDate) {
					t.Errorf("date = %v, want %v", entry.Date, testDate)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseEntryDTO(tt.dto)

			if err != nil {
				t.Errorf("ParseEntryDTO() unexpected error = %v", err)
				return
			}

			if tt.checkResult != nil {
				tt.checkResult(t, got)
			}
		})
	}
}

// Helper function to create float64 pointers
func floatPtr(f float64) *float64 {
	return &f
}
