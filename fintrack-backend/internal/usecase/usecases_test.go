package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"fintrack-backend/internal/domain"
	"fintrack-backend/internal/platform/response"

	"github.com/google/uuid"
)

func TestReportRangeDefaultsToCurrentMonth(t *testing.T) {
	start, end, err := reportRange("", "")
	if err != nil {
		t.Fatalf("reportRange returned error: %v", err)
	}

	now := time.Now().UTC()
	expectedStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if !start.Equal(expectedStart) {
		t.Fatalf("expected start %s, got %s", expectedStart, start)
	}
	if !end.Equal(expectedStart.AddDate(0, 1, 0)) {
		t.Fatalf("expected end %s, got %s", expectedStart.AddDate(0, 1, 0), end)
	}
}

func TestReportRangeParsesExplicitRangeAsExclusiveEnd(t *testing.T) {
	start, end, err := reportRange("2026-06-01", "2026-06-30")
	if err != nil {
		t.Fatalf("reportRange returned error: %v", err)
	}

	expectedStart := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
	if !start.Equal(expectedStart) {
		t.Fatalf("expected start %s, got %s", expectedStart, start)
	}
	if !end.Equal(expectedEnd) {
		t.Fatalf("expected exclusive end %s, got %s", expectedEnd, end)
	}
}

func TestReportRangeRejectsPartialOrInvalidRange(t *testing.T) {
	tests := []struct {
		name      string
		startDate string
		endDate   string
	}{
		{name: "missing end", startDate: "2026-06-01", endDate: ""},
		{name: "missing start", startDate: "", endDate: "2026-06-30"},
		{name: "invalid start", startDate: "invalid", endDate: "2026-06-30"},
		{name: "invalid end", startDate: "2026-06-01", endDate: "invalid"},
		{name: "end before start", startDate: "2026-06-30", endDate: "2026-06-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := reportRange(tt.startDate, tt.endDate)
			if !errors.Is(err, response.ErrBadRequest) {
				t.Fatalf("expected ErrBadRequest, got %v", err)
			}
		})
	}
}

func TestCreateTransactionValidationRejectsInvalidPayloadsBeforeRepository(t *testing.T) {
	uc := &Usecases{}
	accountID := uuid.New()
	categoryID := uuid.New()
	transferAccountID := uuid.New()

	tests := []struct {
		name string
		tx   domain.Transaction
	}{
		{
			name: "zero amount",
			tx: domain.Transaction{
				AccountID:  accountID,
				CategoryID: &categoryID,
				Type:       "expense",
				Amount:     0,
			},
		},
		{
			name: "negative amount",
			tx: domain.Transaction{
				AccountID:  accountID,
				CategoryID: &categoryID,
				Type:       "expense",
				Amount:     -1,
			},
		},
		{
			name: "transfer missing destination",
			tx: domain.Transaction{
				AccountID: accountID,
				Type:      "transfer",
				Amount:    100,
			},
		},
		{
			name: "transfer with category",
			tx: domain.Transaction{
				AccountID:         accountID,
				CategoryID:        &categoryID,
				TransferAccountID: &transferAccountID,
				Type:              "transfer",
				Amount:            100,
			},
		},
		{
			name: "transfer same source and destination",
			tx: domain.Transaction{
				AccountID:         accountID,
				TransferAccountID: &accountID,
				Type:              "transfer",
				Amount:            100,
			},
		},
		{
			name: "expense missing category",
			tx: domain.Transaction{
				AccountID: accountID,
				Type:      "expense",
				Amount:    100,
			},
		},
		{
			name: "income with transfer account",
			tx: domain.Transaction{
				AccountID:         accountID,
				CategoryID:        &categoryID,
				TransferAccountID: &transferAccountID,
				Type:              "income",
				Amount:            100,
			},
		},
		{
			name: "unknown type",
			tx: domain.Transaction{
				AccountID:  accountID,
				CategoryID: &categoryID,
				Type:       "refund",
				Amount:     100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uc.CreateTransaction(context.Background(), tt.tx)
			if !errors.Is(err, response.ErrBadRequest) {
				t.Fatalf("expected ErrBadRequest, got %v", err)
			}
		})
	}
}
