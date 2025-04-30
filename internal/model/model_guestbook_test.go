package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGuestbook(t *testing.T) {
	tests := []struct {
		name     string
		entry    Guestbook
		expected Guestbook
	}{
		{
			name: "Valid guestbook entry",
			entry: Guestbook{
				ID:        1,
				Name:      "John Doe",
				Message:   "Hello World!",
				CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			expected: Guestbook{
				ID:        1,
				Name:      "John Doe",
				Message:   "Hello World!",
				CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "With deleted at",
			entry: Guestbook{
				ID:        2,
				Name:      "Jane Doe",
				Message:   "Test Message",
				CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				DeletedAt: gorm.DeletedAt{Time: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
			},
			expected: Guestbook{
				ID:        2,
				Name:      "Jane Doe",
				Message:   "Test Message",
				CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				DeletedAt: gorm.DeletedAt{Time: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Valid: true},
			},
		},
		{
			name: "Zero values",
			entry: Guestbook{
				ID:      0,
				Name:    "",
				Message: "",
			},
			expected: Guestbook{
				ID:      0,
				Name:    "",
				Message: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.ID, tt.entry.ID)
			assert.Equal(t, tt.expected.Name, tt.entry.Name)
			assert.Equal(t, tt.expected.Message, tt.entry.Message)
			assert.Equal(t, tt.expected.CreatedAt, tt.entry.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, tt.entry.UpdatedAt)
			assert.Equal(t, tt.expected.DeletedAt, tt.entry.DeletedAt)
		})
	}
}
