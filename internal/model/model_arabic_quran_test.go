package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArabicQuran(t *testing.T) {
	tests := []struct {
		name     string
		quran    ArabicQuran
		expected ArabicQuran
	}{
		{
			name: "Valid Arabic Quran entry",
			quran: ArabicQuran{
				Surat: 1,
				Ayat:  1,
				Text:  "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ",
			},
			expected: ArabicQuran{
				Surat: 1,
				Ayat:  1,
				Text:  "بِسْمِ اللَّهِ الرَّحْمَٰنِ الرَّحِيمِ",
			},
		},
		{
			name: "With ID field",
			quran: ArabicQuran{
				ID:    1,
				Surat: 2,
				Ayat:  255,
				Text:  "اللَّهُ لَا إِلَٰهَ إِلَّا هُوَ الْحَيُّ الْقَيُّومُ",
			},
			expected: ArabicQuran{
				ID:    1,
				Surat: 2,
				Ayat:  255,
				Text:  "اللَّهُ لَا إِلَٰهَ إِلَّا هُوَ الْحَيُّ الْقَيُّومُ",
			},
		},
		{
			name: "Zero values",
			quran: ArabicQuran{
				ID:    0,
				Surat: 0,
				Ayat:  0,
				Text:  "",
			},
			expected: ArabicQuran{
				ID:    0,
				Surat: 0,
				Ayat:  0,
				Text:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.ID, tt.quran.ID)
			assert.Equal(t, tt.expected.Surat, tt.quran.Surat)
			assert.Equal(t, tt.expected.Ayat, tt.quran.Ayat)
			assert.Equal(t, tt.expected.Text, tt.quran.Text)
		})
	}
}
