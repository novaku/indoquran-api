package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdIndonesian(t *testing.T) {
	tests := []struct {
		name     string
		indo     IdIndonesian
		expected IdIndonesian
	}{
		{
			name: "Valid Indonesian translation entry",
			indo: IdIndonesian{
				ID:    1,
				Surat: 1,
				Ayat:  1,
				Text:  "Dengan nama Allah Yang Maha Pengasih, Maha Penyayang",
			},
			expected: IdIndonesian{
				ID:    1,
				Surat: 1,
				Ayat:  1,
				Text:  "Dengan nama Allah Yang Maha Pengasih, Maha Penyayang",
			},
		},
		{
			name: "Large surah and ayat numbers",
			indo: IdIndonesian{
				ID:    2,
				Surat: 114,
				Ayat:  6,
				Text:  "dari (golongan) jin dan manusia",
			},
			expected: IdIndonesian{
				ID:    2,
				Surat: 114,
				Ayat:  6,
				Text:  "dari (golongan) jin dan manusia",
			},
		},
		{
			name: "Long translation text",
			indo: IdIndonesian{
				ID:    3,
				Surat: 2,
				Ayat:  255,
				Text:  "Allah, tidak ada tuhan selain Dia. Yang Maha Hidup, Yang terus menerus mengurus (makhluk-Nya), tidak mengantuk dan tidak tidur. Milik-Nya apa yang ada di langit dan apa yang ada di bumi. Tidak ada yang dapat memberi syafaat di sisi-Nya tanpa izin-Nya.",
			},
			expected: IdIndonesian{
				ID:    3,
				Surat: 2,
				Ayat:  255,
				Text:  "Allah, tidak ada tuhan selain Dia. Yang Maha Hidup, Yang terus menerus mengurus (makhluk-Nya), tidak mengantuk dan tidak tidur. Milik-Nya apa yang ada di langit dan apa yang ada di bumi. Tidak ada yang dapat memberi syafaat di sisi-Nya tanpa izin-Nya.",
			},
		},
		{
			name: "Zero values",
			indo: IdIndonesian{
				ID:    0,
				Surat: 0,
				Ayat:  0,
				Text:  "",
			},
			expected: IdIndonesian{
				ID:    0,
				Surat: 0,
				Ayat:  0,
				Text:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.ID, tt.indo.ID)
			assert.Equal(t, tt.expected.Surat, tt.indo.Surat)
			assert.Equal(t, tt.expected.Ayat, tt.indo.Ayat)
			assert.Equal(t, tt.expected.Text, tt.indo.Text)
		})
	}
}
