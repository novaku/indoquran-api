package controllers

import (
	"fmt"
	"strconv"

	"indoquran-api/internal/services/search"
	"indoquran-api/pkg/cache"
	"indoquran-api/pkg/database"

	"github.com/go-redis/redis"

	"github.com/gin-gonic/gin"
)

// SearchControllerInterface defines the contract for search operations
type SearchControllerInterface interface {
	SearchHandler(c *gin.Context)
}

// SearchRequestValidatorInterface defines the contract for request validation
type SearchRequestValidatorInterface interface {
	ValidateRequest(c *gin.Context) (*SearchRequest, error)
}

// SearchResponseBuilderInterface defines the contract for building search responses
type SearchResponseBuilderInterface interface {
	BuildResponse(result *search.SearchResult, request *SearchRequest) *ResultJsonFormat
}

// SearchController handles HTTP requests for search operations
type SearchController struct {
	searchService   search.SearchServiceInterface
	validator       SearchRequestValidatorInterface
	responseBuilder SearchResponseBuilderInterface
}

// SearchRequest represents the search request parameters
type SearchRequest struct {
	Query       string `form:"q" binding:"required"`
	Page        int    `form:"p" binding:"min=1"`
	Juz         int    `form:"juz" binding:"min=0"`
	Surat       int    `form:"surat" binding:"min=0"`
	RowsPerPage int    `form:"n" binding:"min=1"`
}

// DefaultSearchRequestValidator implements SearchRequestValidatorInterface
type DefaultSearchRequestValidator struct{}

// NewDefaultSearchRequestValidator creates a new instance of DefaultSearchRequestValidator
func NewDefaultSearchRequestValidator() SearchRequestValidatorInterface {
	return &DefaultSearchRequestValidator{}
}

// ValidateRequest implements SearchRequestValidatorInterface
func (v *DefaultSearchRequestValidator) ValidateRequest(c *gin.Context) (*SearchRequest, error) {
	query := c.DefaultQuery("q", "")
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	juz, _ := strconv.Atoi(c.DefaultQuery("juz", "0"))
	surat, _ := strconv.Atoi(c.DefaultQuery("surat", "0"))
	rowsPerPage, _ := strconv.Atoi(c.DefaultQuery("n", "10"))

	if query == "" {
		return nil, ErrInvalidRequest("search query is required")
	}

	return &SearchRequest{
		Query:       query,
		Page:        page,
		Juz:         juz,
		Surat:       surat,
		RowsPerPage: rowsPerPage,
	}, nil
}

// DefaultSearchResponseBuilder implements SearchResponseBuilderInterface
type DefaultSearchResponseBuilder struct{}

// NewDefaultSearchResponseBuilder creates a new instance of DefaultSearchResponseBuilder
func NewDefaultSearchResponseBuilder() SearchResponseBuilderInterface {
	return &DefaultSearchResponseBuilder{}
}

// BuildResponse implements SearchResponseBuilderInterface
func (b *DefaultSearchResponseBuilder) BuildResponse(result *search.SearchResult, request *SearchRequest) *ResultJsonFormat {
	return &ResultJsonFormat{
		Results: result.Results,
		Pagination: &Pagination{
			CurrentPage: request.Page,
			RowsPerPage: request.RowsPerPage,
			TotalPages:  (int(result.TotalCount) + request.RowsPerPage - 1) / request.RowsPerPage,
			TotalRows:   int(result.TotalCount),
		},
		Aggregate: result.Count,
	}
}

// NewSearchController creates a new instance of SearchController
func NewSearchController(
	searchService search.SearchServiceInterface,
	validator SearchRequestValidatorInterface,
	responseBuilder SearchResponseBuilderInterface,
) SearchControllerInterface {
	return &SearchController{
		searchService:   searchService,
		validator:       validator,
		responseBuilder: responseBuilder,
	}
}

// SearchHandler handles the search request
func (sc *SearchController) SearchHandler(c *gin.Context) {
	// Validate request
	request, err := sc.validator.ValidateRequest(c)
	if err != nil {
		WriteResponse(c, nil, fmt.Errorf("validation error: %w", err))
		return
	}

	// Perform search
	result, err := sc.searchService.FullTextSearch(
		request.Query,
		request.Juz,
		request.Surat,
		request.Page,
		request.RowsPerPage,
	)
	if err != nil {
		WriteResponse(c, nil, fmt.Errorf("search error: %w", err))
		return
	}

	// Build and send response
	response := sc.responseBuilder.BuildResponse(result, request)
	WriteResponse(c, response, nil)
}

// DefaultSearchController returns a new SearchController with default implementations
func DefaultSearchController() SearchControllerInterface {
	return NewSearchController(
		search.NewSearchService(
			database.NewMySQLDatabase().GetConnection(),
			redis.NewClient(&redis.Options{
				Addr:     cache.NewRedisConfig().GetAddress(),
				Password: cache.NewRedisConfig().GetPassword(),
				DB:       cache.NewRedisConfig().GetDB(),
			}),
		),
		NewDefaultSearchRequestValidator(),
		NewDefaultSearchResponseBuilder(),
	)
}

// SearchHandler is the global handler function that uses the controller
func SearchHandler(c *gin.Context) {
	DefaultSearchController().SearchHandler(c)
}

// ErrInvalidRequest represents an invalid request error
type ErrInvalidRequest string

func (e ErrInvalidRequest) Error() string {
	return string(e)
}
