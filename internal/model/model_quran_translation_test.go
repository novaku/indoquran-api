package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuranTranslation(t *testing.T) {
	tests := []struct {
		name        string
		translation QuranTranslation
		expected    QuranTranslation
	}{
		{
			name: "Valid translation entry",
			translation: QuranTranslation{
				TranslationID:  1,
				AyatKey:        "1_1",
				TranslationKey: "en_sahih",
				Surat:          1,
				Ayat:           1,
				Translation:    "In the name of Allah, the Entirely Merciful, the Especially Merciful",
			},
			expected: QuranTranslation{
				TranslationID:  1,
				AyatKey:        "1_1",
				TranslationKey: "en_sahih",
				Surat:          1,
				Ayat:           1,
				Translation:    "In the name of Allah, the Entirely Merciful, the Especially Merciful",
			},
		},
		{
			name: "Different language translation",
			translation: QuranTranslation{
				TranslationID:  2,
				AyatKey:        "2_255",
				TranslationKey: "id_indonesian",
				Surat:          2,
				Ayat:           255,
				Translation:    "Allah, tidak ada Tuhan selain Dia. Yang Maha Hidup, Yang terus menerus mengurus makhluk-Nya",
			},
			expected: QuranTranslation{
				TranslationID:  2,
				AyatKey:        "2_255",
				TranslationKey: "id_indonesian",
				Surat:          2,
				Ayat:           255,
				Translation:    "Allah, tidak ada Tuhan selain Dia. Yang Maha Hidup, Yang terus menerus mengurus makhluk-Nya",
			},
		},
		{
			name: "Maximum field lengths",
			translation: QuranTranslation{
				TranslationID:  999999,
				AyatKey:        "114_6",
				TranslationKey: "fr_french_max_chars",
				Surat:          114,
				Ayat:           6,
				Translation:    "Long translation text with multiple paragraphs and special characters: é à ç",
			},
			expected: QuranTranslation{
				TranslationID:  999999,
				AyatKey:        "114_6",
				TranslationKey: "fr_french_max_chars",
				Surat:          114,
				Ayat:           6,
				Translation:    "Long translation text with multiple paragraphs and special characters: é à ç",
			},
		},
		{
			name: "Zero values",
			translation: QuranTranslation{
				TranslationID:  0,
				AyatKey:        "",
				TranslationKey: "",
				Surat:          0,
				Ayat:           0,
				Translation:    "",
			},
			expected: QuranTranslation{
				TranslationID:  0,
				AyatKey:        "",
				TranslationKey: "",
				Surat:          0,
				Ayat:           0,
				Translation:    "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.TranslationID, tt.translation.TranslationID)
			assert.Equal(t, tt.expected.AyatKey, tt.translation.AyatKey)
			assert.Equal(t, tt.expected.TranslationKey, tt.translation.TranslationKey)
			assert.Equal(t, tt.expected.Surat, tt.translation.Surat)
			assert.Equal(t, tt.expected.Ayat, tt.translation.Ayat)
			assert.Equal(t, tt.expected.Translation, tt.translation.Translation)
		})
	}
}
