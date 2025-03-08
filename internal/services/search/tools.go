package search

import (
	"fmt"
	"strings"

	"indoquran-api/internal/database"
	"indoquran-api/internal/model"
	"indoquran-api/pkg/logger"

	"gorm.io/gorm"
)

// buildSQLLikeClause builds the LIKE clause for the search query.
func buildSQLLikeClause(input string) string {
	// Split the input by spaces
	words := strings.Split(input, " ")

	// Initialize the LIKE clause
	likeClause := ""

	// Loop through the words to build the LIKE clause
	// limit search for 5 words only
	for i, word := range words {
		if i > 0 && i < 5 {
			likeClause += " AND "
		}
		likeClause += fmt.Sprintf("quran_translation.translation LIKE '%%%s%%'", word)
	}

	return likeClause
}

// querySearch performs the search query on the database.
func queryTotalCount(searchTerm string, juz, surat int) (int64, error) {
	var (
		dbClient          = database.GetDB()
		querySessionCount *gorm.DB
		totalCount        int64
		err               error
	)

	querySessionCount = dbClient.Table("quran_translation").
		Where(searchTerm).
		Joins("INNER JOIN quran_ayat ON quran_ayat.ayat_number = quran_translation.translation_id")
	if juz > 0 {
		querySessionCount = querySessionCount.Where("quran_ayat.juz = ?", juz)
	}
	if surat > 0 {
		querySessionCount = querySessionCount.Where("quran_ayat.surat = ?", surat)
	}
	err = querySessionCount.Count(&totalCount).Error
	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error fetching total count: %#v", err)
		return 0, err
	}

	return totalCount, nil
}

// querySearch performs the search query on the database.
func querySearch(searchTerm string, limit, offset, juz, surat int) ([]*model.AyatDetail, error) {
	var (
		dbClient     = database.GetDB()
		querySession *gorm.DB
		results      []*model.AyatDetail
		err          error
	)

	logger.WriteLog(logger.LogLevelInfo, "querySearch called with searchTerm: %s, limit: %d, offset: %d, juz: %d, surat: %d", searchTerm, limit, offset, juz, surat)

	// Perform the search query with pagination
	querySession = dbClient.
		Table("quran_translation").
		Select("quran_translation.translation AS text_indo, quran_ayat.text AS text_arabic, quran_ayat.juz, quran_translation.surat AS surat, quran_translation.ayat AS ayat, quran_translation.translation_id AS id").
		Where(searchTerm).
		Limit(limit).
		Offset(offset).
		Order("quran_translation.translation_id ASC").
		Joins("JOIN quran_ayat ON quran_ayat.ayat_number = quran_translation.translation_id")

	if juz > 0 {
		querySession = querySession.Where("quran_ayat.juz = ?", juz)
	}
	if surat > 0 {
		querySession = querySession.Where("id_indonesian.surat = ?", surat)
	}
	err = querySession.Scan(&results).Error

	if err != nil {
		logger.WriteLog(logger.LogLevelError, "Error fetching from database: %#v", err)
		return nil, err
	}

	return results, nil
}

// queryAggregate performs the aggregate query on the database.
func queryAggregate(searchTerm string, juz, surat int) []*model.CountResult {
	var (
		dbClient   = database.GetDB()
		queryParts []string
		count      []*model.CountResult
	)

	// Append conditions based on the values of juz and sura
	if juz > 0 {
		queryParts = append(queryParts, "quran_ayat.juz = "+fmt.Sprintf("%d", juz))
	}
	if surat > 0 {
		queryParts = append(queryParts, "quran_ayat.surat = "+fmt.Sprintf("%d", surat))
	}

	// Concatenate all parts with "AND" and a leading space
	addQuery := ""
	if len(queryParts) > 0 {
		addQuery = "AND " + strings.Join(queryParts, " AND ")
	}

	// Execute the query
	query := `
		(
			SELECT 
				'juz' AS type,
				quran_ayat.juz AS identifier,
				COUNT(*) AS count
			FROM quran_translation
			INNER JOIN quran_ayat ON quran_ayat.ayat_number = quran_translation.translation_id
			WHERE ` + searchTerm + `
			` + addQuery + `
			GROUP BY quran_ayat.juz
		)
		UNION ALL
		(
			SELECT 
				'surat' AS type,
				quran_translation.surat AS identifier,
				COUNT(*) AS count
			FROM quran_translation
			INNER JOIN quran_ayat ON quran_ayat.ayat_number = quran_translation.translation_id
			WHERE ` + searchTerm + `
			` + addQuery + `
			GROUP BY quran_translation.surat
		)
		ORDER BY type, identifier;
		`
	dbClient.Raw(query).Scan(&count)

	return count
}
