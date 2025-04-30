package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdMuntakhab(t *testing.T) {
	tests := []struct {
		name      string
		muntakhab IdMuntakhab
		expected  IdMuntakhab
	}{
		{
			name: "Valid Id Muntakhab entry",
			muntakhab: IdMuntakhab{
				Index: 1,
				Surat: 1,
				Ayat:  1,
				Text:  "Dengan nama Allah Yang Maha Pengasih, Maha Penyayang",
			},
			expected: IdMuntakhab{
				Index: 1,
				Surat: 1,
				Ayat:  1,
				Text:  "Dengan nama Allah Yang Maha Pengasih, Maha Penyayang",
			},
		},
		{
			name: "With large index and values",
			muntakhab: IdMuntakhab{
				Index: 6236,
				Surat: 114,
				Ayat:  6,
				Text:  "dari (golongan) jin dan manusia",
			},
			expected: IdMuntakhab{
				Index: 6236,
				Surat: 114,
				Ayat:  6,
				Text:  "dari (golongan) jin dan manusia",
			},
		},
		{
			name: "Zero values except index",
			muntakhab: IdMuntakhab{
				Index: 1,
				Surat: 0,
				Ayat:  0,
				Text:  "",
			},
			expected: IdMuntakhab{
				Index: 1,
				Surat: 0,
				Ayat:  0,
				Text:  "",
			},
		},
		{
			name: "With long text content",
			muntakhab: IdMuntakhab{
				Index: 2,
				Surat: 2,
				Ayat:  255,
				Text:  "Allah, tidak ada tuhan selain Dia. Yang Maha Hidup, Yang terus menerus mengurus (makhluk-Nya), tidak mengantuk dan tidak tidur. Milik-Nya apa yang ada di langit dan apa yang ada di bumi. Tidak ada yang dapat memberi syafaat di sisi-Nya tanpa izin-Nya. Dia mengetahui apa yang di hadapan mereka dan apa yang di belakang mereka, dan mereka tidak mengetahui sesuatu apa pun tentang ilmu-Nya melainkan apa yang Dia kehendaki. Kursi-Nya meliputi langit dan bumi. Dan Dia tidak merasa berat memelihara keduanya, dan Dia Mahatinggi, Mahabesar.",
			},
			expected: IdMuntakhab{
				Index: 2,
				Surat: 2,
				Ayat:  255,
				Text:  "Allah, tidak ada tuhan selain Dia. Yang Maha Hidup, Yang terus menerus mengurus (makhluk-Nya), tidak mengantuk dan tidak tidur. Milik-Nya apa yang ada di langit dan apa yang ada di bumi. Tidak ada yang dapat memberi syafaat di sisi-Nya tanpa izin-Nya. Dia mengetahui apa yang di hadapan mereka dan apa yang di belakang mereka, dan mereka tidak mengetahui sesuatu apa pun tentang ilmu-Nya melainkan apa yang Dia kehendaki. Kursi-Nya meliputi langit dan bumi. Dan Dia tidak merasa berat memelihara keduanya, dan Dia Mahatinggi, Mahabesar.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected.Index, tt.muntakhab.Index)
			assert.Equal(t, tt.expected.Surat, tt.muntakhab.Surat)
			assert.Equal(t, tt.expected.Ayat, tt.muntakhab.Ayat)
			assert.Equal(t, tt.expected.Text, tt.muntakhab.Text)
		})
	}
}
