package quran

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"bitbucket.org/indoquran-api/src/handlers"
	"bitbucket.org/indoquran-api/src/models/quran"
	model_quran "bitbucket.org/indoquran-api/src/models/quran"
	"bitbucket.org/indoquran-api/src/models/quran/alqurancloud"
	"bitbucket.org/indoquran-api/src/models/quran/banghasan"
	"bitbucket.org/indoquran-api/src/models/quran/github"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ImportSurat : function to import surat from outsource API
func ImportSurat(c *gin.Context) {
	result := &banghasan.SuratHasil{}

	resp, err := http.Get("https://api.banghasan.com/quran/format/json/surat")
	if err != nil {
		mlog.Error(err)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&result)
	var surats []interface{}

	for _, data := range result.Hasil {
		no, err := strconv.Atoi(data.Nomor)
		if err != nil {
			mlog.Error(err)
		}
		jumAyat, err := strconv.Atoi(data.Ayat)
		if err != nil {
			mlog.Error(err)
		}
		surats = append(surats, model_quran.Surat{
			ID:         primitive.NewObjectID(),
			Nomor:      no,
			Nama:       data.Nama,
			Asma:       data.Asma,
			JumlahAyat: jumAyat,
			Arti:       data.Arti,
			Keterangan: data.Keterangan,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	res, err := suratCollection.InsertMany(c, surats)

	if err != nil {
		mlog.Error(err)
	}

	handlers.DefaultResponse(c, http.StatusBadRequest, "Success Get Data", res.InsertedIDs)
}

// ImportAyat : function to import ayat
func ImportAyat(c *gin.Context) {
	var insertedID []interface{}

	cur, err := suratCollection.Find(c, bson.M{})
	if err != nil {
		mlog.Error(err)
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var result model_quran.Surat
		err := cur.Decode(&result)
		if err != nil {
			mlog.Error(err)
		}

		// if result.Nomor >= 56 {
		for i := 1; i <= result.JumlahAyat; i++ {
			// if result.Nomor == 56 && i >= 46 {
			// 	getAyatHTTPDataInsertDB(c, result.Nomor, i, result.JumlahAyat)
			// } else if result.Nomor > 56 {
			getAyatHTTPDataInsertDB(c, result.Nomor, i, result.JumlahAyat)
			// }
		}
		// }
	}
	if err := cur.Err(); err != nil {
		mlog.Error(err)
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success Get Data", insertedID)
}

func getAyatHTTPDataInsertDB(c *gin.Context, suratID, ayatID, jumAyat int) {
	modelAyat := &banghasan.Ayat{}

	resp, err := http.Get("https://api.banghasan.com/quran/format/json/surat/" + strconv.Itoa(suratID) + "/ayat/" + strconv.Itoa(ayatID))
	if err != nil {
		mlog.Error(err)
	}

	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&modelAyat)

	ayat := &model_quran.Ayat{
		ID:        primitive.NewObjectID(),
		Nomor:     ayatID,
		Surat:     suratID,
		TxtAR:     modelAyat.Ayat.Data.AR[0].Teks,
		TxtID:     modelAyat.Ayat.Data.ID[0].Teks,
		TxtIDT:    modelAyat.Ayat.Data.IDT[0].Teks,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := ayatCollection.InsertOne(c, ayat)
	if err != nil {
		mlog.Error(err)
	}

	mlog.Info("Inserted, Surat: %d, Ayat: %d / %d, DB_ID: %+v", ayat.Surat, ayat.Nomor, jumAyat, res.InsertedID)
}

// ImportCatatan : function to import catatan
func ImportCatatan(c *gin.Context) {
	maxCatatanID := 1610
	catatan := &banghasan.CatatanRoot{}

	for i := 1; i <= maxCatatanID; i++ {
		resp, err := http.Get("https://api.banghasan.com/quran/format/json/catatan/" + strconv.Itoa(i))
		if err != nil {
			mlog.Error(err)
		}

		defer resp.Body.Close()

		json.NewDecoder(resp.Body).Decode(&catatan)

		suratID, err := strconv.Atoi(catatan.Quran.Surat)
		if err != nil {
			mlog.Error(err)
		}
		ayatID, err := strconv.Atoi(catatan.Quran.Ayat)
		if err != nil {
			mlog.Error(err)
		}

		modelCatatan := model_quran.Catatan{
			ID:          primitive.NewObjectID(),
			Nomor:       i,
			Surat:       suratID,
			Ayat:        ayatID,
			TeksCatatan: catatan.Catatan.Teks,
			TeksQuran:   catatan.Quran.Teks,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		res, err := catatanCollection.InsertOne(c, modelCatatan)
		if err != nil {
			mlog.Error(err)
		}

		mlog.Info("Inserted, CatatanID: %d, Surat: %d, Ayat: %d, DB_ID: %+v", i, suratID, ayatID, res.InsertedID)
	}
}

// ImportTafsir : function to import tafsir
func ImportTafsir(c *gin.Context) {
	var surat quran.Surat
	var success []string

	cursor, err := suratCollection.Find(c, bson.D{}, options.Find().SetSort(bson.D{{Key: "nomor", Value: 1}}))
	if err != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed fetch database", err)
		return
	}
	defer cursor.Close(c)
	for cursor.Next(c) {
		if err := cursor.Decode(&surat); err != nil {
			mlog.Error(err)
		}

		tafsirStruct := github.CreateStruct(fmt.Sprintf("%d", surat.Nomor), surat.JumlahAyat)
		json, _ := json.Marshal(tafsirStruct)

		for i := 1; i <= surat.JumlahAyat; i++ {
			tafsirText := gjson.Get(string(json), fmt.Sprintf("%d", surat.Nomor)+".tafsir.id.kemenag.text."+fmt.Sprintf("%d", i))
			modelTafsir := model_quran.Tafsir{
				ID:         primitive.NewObjectID(),
				Surat:      surat.Nomor,
				Ayat:       i,
				TextTafsir: tafsirText.String(),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			res, err := tafsirCollection.InsertOne(c, modelTafsir)
			if err != nil {
				mlog.Error(err)
			}

			mlog.Info("Inserted, Tafsir Surat: %d, Ayat: %d, DB_ID: %+v", surat.Nomor, i, res.InsertedID)
		}
		success = append(success, fmt.Sprintf("Surat: %s, id: %d, sukses!", surat.Nama, surat.Nomor))
	}

	handlers.DefaultResponse(c, http.StatusOK, "Berhasil", success)
}

// ImportJuz : import juz of ayat
func ImportJuz(c *gin.Context) {
	ayatCloud := &alqurancloud.Ayat{}
	var updated = []string{}

	cur, err := ayatCollection.Find(c, bson.M{})
	if err != nil {
		mlog.Error(err)
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var result model_quran.Ayat
		err := cur.Decode(&result)
		if err != nil {
			mlog.Error(err)
		}

		resp, err := http.Get("http://api.alquran.cloud/v1/ayah/" + strconv.Itoa(result.Surat) + ":" + strconv.Itoa(result.Nomor) + "/id.indonesian")
		if err != nil {
			mlog.Error(err)
		}

		defer resp.Body.Close()

		json.NewDecoder(resp.Body).Decode(&ayatCloud)

		filter := bson.M{"_id": bson.M{"$eq": result.ID}}
		update := bson.M{"$set": bson.M{"juz": ayatCloud.Data.Juz}}
		_, err = ayatCollection.UpdateOne(c, filter, update)
		if err != nil {
			mlog.Error(err)
		}
		updated = append(updated, fmt.Sprintf("Updated Surat: %d, Ayat: %d, Juz: %d", result.Surat, result.Nomor, ayatCloud.Data.Juz))
		mlog.Info("Updated Surat: %d, Ayat: %d, Juz: %d", result.Surat, result.Nomor, ayatCloud.Data.Juz)
	}
	if err := cur.Err(); err != nil {
		mlog.Error(err)
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success udate Juz", updated)
}

// ImportTafsirMove : function to move tafsir to ayat collection
func ImportTafsirMove(c *gin.Context) {
	var updated = []string{}

	cur, err := ayatCollection.Find(c, bson.M{})
	if err != nil {
		mlog.Error(err)
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var result model_quran.Ayat
		err := cur.Decode(&result)
		var tafsir model_quran.Tafsir
		if err = tafsirCollection.FindOne(c, bson.M{"surat": result.Surat, "ayat": result.Nomor}).Decode(&tafsir); err != nil {
			mlog.Error(err)
		}

		filter := bson.M{"_id": bson.M{"$eq": result.ID}}
		update := bson.M{"$set": bson.M{"txt_tafsir": tafsir.TextTafsir}}
		_, err = ayatCollection.UpdateOne(c, filter, update)
		if err != nil {
			mlog.Error(err)
		}
		updated = append(updated, fmt.Sprintf("Updated Surat: %d, Ayat: %d", result.Surat, result.Nomor))
		mlog.Info("Updated Surat: %d, Ayat: %d", result.Surat, result.Nomor)
	}
	if err := cur.Err(); err != nil {
		mlog.Error(err)
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success udate Juz", updated)
}

// ImportImage : import image (hardcode image url)
func ImportImage(c *gin.Context) {
	var updated = []string{}

	cur, err := ayatCollection.Find(c, bson.M{})
	if err != nil {
		mlog.Error(err)
	}
	defer cur.Close(c)
	for cur.Next(c) {
		var result model_quran.Ayat
		err := cur.Decode(&result)
		image := fmt.Sprintf(imageURL, result.Surat, result.Nomor)

		filter := bson.M{"_id": bson.M{"$eq": result.ID}}
		update := bson.M{"$set": bson.M{"image": image}}
		_, err = ayatCollection.UpdateOne(c, filter, update)
		if err != nil {
			mlog.Error(err)
		}

		t := fmt.Sprintf("Updated Surat: %d, Ayat: %d, Image: %s", result.Surat, result.Nomor, image)
		updated = append(updated, t)
		mlog.Info(t)
	}
	if err := cur.Err(); err != nil {
		mlog.Error(err)
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success udate Juz", updated)
}

// ImportArabText : function to import Arabic text from scrapping to URL
func ImportArabText(c *gin.Context) {
	// res, err := http.Get("https://kalam.sindonews.com/next/114/surat/0")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer res.Body.Close()

	// if res.StatusCode != 200 {
	// 	log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	// }

	// doc, err := goquery.NewDocumentFromReader(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // NOTE: font-family : Uthmani

	// result := []string{}
	// doc.Find(".ayat-arab").Each(func(i int, sel *goquery.Selection) {
	// 	arTxt := sel.Text()
	// 	log.Println(arTxt)
	// 	result = append(result, arTxt)
	// })

	// handlers.DefaultResponse(c, http.StatusOK, "Success get http", result)

	//======================================================================================================

	cur, err := suratCollection.Find(c, bson.M{})
	if err != nil {
		mlog.Error(err)
	}
	upd := []int{}
	defer cur.Close(c)
	for cur.Next(c) {
		var result model_quran.Surat
		err := cur.Decode(&result)
		if err != nil {
			mlog.Error(err)
		}

		perPage := 20
		div := float64(result.JumlahAyat) / float64(perPage)
		pageNum := int(math.Ceil(div))
		mlog.Info("hasil bagi : %v, page count : %v", div, pageNum)

		for i := 1; i <= pageNum; i++ {
			n := i - 1
			if i > 1 {
				n = (i - 1) * perPage
			}

			urlScrap := "https://kalam.sindonews.com/next/" + strconv.Itoa(result.Nomor) + "/surat/" + strconv.Itoa(n)
			res, err := http.Get(urlScrap)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}

			// NOTE: font-family : Uthmani

			suratNo := result.Nomor
			ayatNo := ((i - 1) * perPage) + 1

			doc.Find(".ayat-arab").Each(func(i int, sel *goquery.Selection) {
				arTxt := sel.Text()

				findFilter := bson.M{
					"$and": []bson.M{
						{"surat": suratNo},
						{"nomor": ayatNo},
					},
				}
				update := bson.M{"$set": bson.M{"txt_ar": arTxt}}
				_, err := ayatCollection.UpdateOne(c, findFilter, update)
				if err != nil {
					mlog.Error(err)
				}

				log.Println(fmt.Sprintf("Updated Surat: %d, Ayat: %d, url: %s, ArabText: %s", result.Nomor, ayatNo, urlScrap, arTxt))

				ayatNo++
			})
		}
		upd = append(upd, result.Nomor)
	}
	if err := cur.Err(); err != nil {
		mlog.Error(err)
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Success update arab text from http", err)
	}

	handlers.DefaultResponse(c, http.StatusOK, "Success update arab text from http", upd)
}

// ImportWiseWords : function to import wise words from local file
func ImportWiseWords(c *gin.Context) {
	file, err := os.Open("resources/wisewords.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		kata := &model_quran.KataBijak{
			ID:   primitive.NewObjectID(),
			Kata: scanner.Text(),
		}

		res, err := kataBijakCollection.InsertOne(c, kata)
		if err != nil {
			mlog.Error(err)
		}
		fmt.Println(res.InsertedID)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// AddAyatID : add ID for every ayat
func AddAyatID(c *gin.Context) {
	var surat quran.Surat
	var ayatID int = 1

	cursor, err := suratCollection.Find(c, bson.D{}, options.Find().SetSort(bson.D{{Key: "nomor", Value: 1}}))
	if err != nil {
		handlers.DefaultResponse(c, http.StatusInternalServerError, "Failed fetch database", err)
		return
	}
	defer cursor.Close(c)
	for cursor.Next(c) {
		if err := cursor.Decode(&surat); err != nil {
			mlog.Error(err)
		}

		for i := 1; i <= surat.JumlahAyat; i++ {
			filter := bson.M{
				"$and": []bson.M{
					{"surat": surat.Nomor},
					{"nomor": i},
				},
			}
			update := bson.M{"$set": bson.M{"ayat_id": ayatID}}
			_, err = ayatCollection.UpdateOne(c, filter, update)
			if err != nil {
				mlog.Error(err)
			}
			mlog.Info("Update Surat: %d, Ayat: %s, ayat_id: %d", surat.Nomor, i, ayatID)
			ayatID++
		}
	}
	handlers.DefaultResponse(c, http.StatusOK, "Success update ayat_id", nil)
}
