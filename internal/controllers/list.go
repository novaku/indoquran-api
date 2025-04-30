package controllers

import (
	"strconv"

	"indoquran-api/internal/services/list"
	"indoquran-api/pkg/cache"
	"indoquran-api/pkg/database"

	"github.com/gin-gonic/gin"
)

// ListControllerInterface defines the contract for list operations
type ListControllerInterface interface {
	GetSuratList(c *gin.Context)
	GetAyatList(c *gin.Context)
}

// ListController handles list-related operations
type ListController struct {
	suratService list.SuratService
	ayatService  list.AyatService
}

// NewListController creates a new instance of ListController
func NewListController(suratService list.SuratService, ayatService list.AyatService) ListControllerInterface {
	return &ListController{
		suratService: suratService,
		ayatService:  ayatService,
	}
}

// GetSuratList handles GET /surat endpoint
func (l *ListController) GetSuratList(c *gin.Context) {
	suratID := c.DefaultQuery("surat", "")
	surats, err := l.suratService.GetSuratList(suratID)
	WriteResponse(c, surats, err)
}

// GetAyatList handles GET /surat/:id/ayat endpoint
func (l *ListController) GetAyatList(c *gin.Context) {
	suratID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("n", "10"))
	ayatList, err := l.ayatService.GetAyatList(suratID, page, pageSize)
	WriteResponse(c, ayatList, err)
}

// DefaultListController returns a new ListController with default implementations
func DefaultListController() ListControllerInterface {
	return NewListController(
		list.NewSurat(list.NewRedisCacheService(), list.NewGormDatabaseService()),
		list.NewAyat(list.NewDatabaseRepository(database.GetDB()), list.NewRedisCache(cache.GetRedis())),
	)
}

// ListSurat is deprecated, use DefaultListController().GetSuratList instead
func ListSurat(c *gin.Context) {
	DefaultListController().GetSuratList(c)
}

// ListAyatInSurat is deprecated, use DefaultListController().GetAyatList instead
func ListAyatInSurat(c *gin.Context) {
	DefaultListController().GetAyatList(c)
}
