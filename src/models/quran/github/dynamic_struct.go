package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jbrodriguez/mlog"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

func CreateStruct(idSurat string, totalAyat int) interface{} {
	var build dynamicstruct.DynamicStruct

	for i := 1; i <= totalAyat; i++ {
		newText := dynamicstruct.NewStruct().AddField(fmt.Sprintf("Num%d", i), "", `json:"`+fmt.Sprintf("%d", i)+`"`).Build().New()
		build = dynamicstruct.ExtendStruct(newText).Build()
	}
	textSlice := build.NewSliceOfStructs()

	kemenag := dynamicstruct.NewStruct().
		AddField("Name", "", `json:"name"`).
		AddField("Source", "", `json:"source"`).
		AddField("Text", &textSlice, `json:"text"`).
		Build().New()

	id := dynamicstruct.NewStruct().
		AddField("Kemenag", kemenag, `json:"kemenag"`).
		Build().New()

	tafsir := dynamicstruct.NewStruct().
		AddField("ID", id, `json:"id"`).
		Build().New()

	numSurat := dynamicstruct.NewStruct().
		AddField("Tafsir", tafsir, `json:"tafsir"`).
		AddField("NumberOfAyah", "", `json:"number_of_ayah"`).
		Build().New()

	tafsirStruct := dynamicstruct.NewStruct().
		AddField("Num", numSurat, `json:"`+idSurat+`"`).
		Build().New()

	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/rioastamal/quran-json/master/surah/%s.json", idSurat))
	if err != nil {
		mlog.Error(err)
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&tafsirStruct)
	if err != nil {
		mlog.Error(err)
	}

	return tafsirStruct
}
