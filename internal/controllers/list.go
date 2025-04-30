package controllers

import (
	"strconv"

	"indoquran-api/internal/services/list"
	"indoquran-api/pkg/cache"
	"indoquran-api/pkg/database"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
)

// List handles list-related operations
type List struct {
	suratService list.SuratService
	ayatService  list.AyatService
}

// NewList creates a new instance of List
func NewList() *List {
	return &List{
		suratService: list.NewSurat(
			list.NewRedisCacheService(),
			list.NewGormDatabaseService(),
		),
		ayatService: list.NewAyat(
			list.NewDatabaseRepository(database.NewMySQLDatabase().GetConnection()),
			list.NewRedisCache(redis.NewClient(&redis.Options{
				Addr:     cache.NewRedisConfig().GetAddress(),
				Password: cache.NewRedisConfig().GetPassword(),
				DB:       cache.NewRedisConfig().GetDB(),
			})),
		),
	}
}

// GetSuratList handles GET /surat endpoint
func (l *List) GetSuratList(c *gin.Context) {
	suratID := c.DefaultQuery("surat", "")
	surats, err := l.suratService.GetSuratList(suratID)
	WriteResponse(c, surats, err)
}

// GetAyatList handles GET /surat/:id/ayat endpoint
func (l *List) GetAyatList(c *gin.Context) {
	suratID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("n", "10"))
	ayatList, err := l.ayatService.GetAyatList(suratID, page, pageSize)
	WriteResponse(c, ayatList, err)
}

// ListSurat is deprecated, use NewList().GetSuratList instead
func ListSurat(c *gin.Context) {
	NewList().GetSuratList(c)
}

// ListAyatInSurat is deprecated, use NewList().GetAyatList instead
func ListAyatInSurat(c *gin.Context) {
	NewList().GetAyatList(c)
}
