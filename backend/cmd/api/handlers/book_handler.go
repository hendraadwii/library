package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hendraadwii/library/internal/models"
)

// BookHandler handles book-related requests
type BookHandler struct {
	BookService *models.BookService
}

// BookResponse represents a book response
type BookResponse struct {
	ID            int64  `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Stock         int    `json:"stock"`
	ISBN          string `json:"isbn,omitempty"`
	PublishedYear int    `json:"published_year,omitempty"`
	Category      string `json:"category,omitempty"`
	Description   string `json:"description,omitempty"`
	Cover         string `json:"cover,omitempty"`
}

// CreateBookRequest represents a request to create a book
type CreateBookRequest struct {
	Title         string `form:"title" binding:"required"`
	Author        string `form:"author" binding:"required"`
	Stock         int    `form:"stock" binding:"min=0"`
	ISBN          string `form:"isbn" binding:"omitempty"`
	PublishedYear int    `form:"published_year" binding:"omitempty"`
	Category      string `form:"category" binding:"omitempty"`
	Description   string `form:"description" binding:"omitempty"`
}

// UpdateBookRequest represents a request to update a book
type UpdateBookRequest struct {
	Title         string `json:"title" binding:"omitempty"`
	Author        string `json:"author" binding:"omitempty"`
	Stock         int    `json:"stock" binding:"omitempty,min=0"`
	ISBN          string `json:"isbn" binding:"omitempty"`
	PublishedYear int    `json:"published_year" binding:"omitempty"`
	Category      string `json:"category" binding:"omitempty"`
	Description   string `json:"description" binding:"omitempty"`
}

// PaginatedBooksResponse represents a paginated list of books
type PaginatedBooksResponse struct {
	Books      []BookResponse `json:"books"`
	TotalCount int            `json:"total_count"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewBookHandler creates a new BookHandler
func NewBookHandler(bookService *models.BookService) *BookHandler {
	return &BookHandler{BookService: bookService}
}

// GetBooks handles retrieval of all books with pagination
// @Summary Get all books
// @Description Get a paginated list of books with optional filters
// @Tags books
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Param title query string false "Filter by title (substring match)"
// @Param author query string false "Filter by author (substring match)"
// @Param category query string false "Filter by exact category"
// @Success 200 {object} PaginatedBooksResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Parse filter parameters
	filters := models.BookFilters{
		Title:    c.Query("title"),
		Author:   c.Query("author"),
		Category: c.Query("category"),
		Limit:    perPage,
		Offset:   offset,
	}

	// Get books from database
	books, err := h.BookService.List(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting books: " + err.Error()})
		return
	}

	// Get total count for pagination
	totalCount, err := h.BookService.Count(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting books: " + err.Error()})
		return
	}

	// Calculate total pages
	totalPages := totalCount / perPage
	if totalCount%perPage > 0 {
		totalPages++
	}

	// Map to response format
	var bookResponses []BookResponse
	for _, book := range books {
		bookResponses = append(bookResponses, BookResponse{
			ID:            book.ID,
			Title:         book.Title,
			Author:        book.Author,
			Stock:         book.Stock,
			ISBN:          book.ISBN,
			PublishedYear: book.PublishedYear,
			Category:      book.Category,
			Description:   book.Description,
			Cover:         book.Cover,
		})
	}

	c.JSON(http.StatusOK, PaginatedBooksResponse{
		Books:      bookResponses,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// GetBook handles retrieval of a specific book by ID
// @Summary Get a book
// @Description Get a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} BookResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {
	// Parse book ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Get book from database
	book, err := h.BookService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return the book
	c.JSON(http.StatusOK, BookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		Stock:         book.Stock,
		ISBN:          book.ISBN,
		PublishedYear: book.PublishedYear,
		Category:      book.Category,
		Description:   book.Description,
		Cover:         book.Cover,
	})
}

// CreateBook handles creation of a new book
// @Summary Create a book
// @Description Create a new book
// @Tags books
// @Accept json
// @Produce json
// @Param book body CreateBookRequest true "Book details"
// @Success 201 {object} BookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate ISBN if provided
	if req.ISBN != "" {
		// Check if ISBN already exists
		exists, err := h.BookService.ISBNExists(req.ISBN)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking ISBN: " + err.Error()})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ISBN already exists"})
			return
		}
	}

	// Handle file upload
	file, err := c.FormFile("cover")
	var coverFileName string
	if err == nil && file != nil {
		// Generate unique filename
		coverFileName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		coverPath := "./uploads/cover/" + coverFileName
		if err := c.SaveUploadedFile(file, coverPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cover image"})
			return
		}
	}

	// Create book in database
	book := &models.Book{
		Title:         req.Title,
		Author:        req.Author,
		Stock:         req.Stock,
		ISBN:          req.ISBN,
		PublishedYear: req.PublishedYear,
		Category:      req.Category,
		Description:   req.Description,
		Cover:         coverFileName,
	}

	if err := h.BookService.Create(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating book: " + err.Error()})
		return
	}

	// Return the created book
	c.JSON(http.StatusCreated, BookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		Stock:         book.Stock,
		ISBN:          book.ISBN,
		PublishedYear: book.PublishedYear,
		Category:      book.Category,
		Description:   book.Description,
		Cover:         book.Cover,
	})
}

// UpdateBook handles updating an existing book
// @Summary Update a book
// @Description Update an existing book
// @Tags books
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Book ID"
// @Param title formData string false "Book title"
// @Param author formData string false "Book author"
// @Param stock formData int false "Book stock"
// @Param isbn formData string false "Book ISBN"
// @Param published_year formData int false "Published year"
// @Param category formData string false "Book category"
// @Param description formData string false "Book description"
// @Param cover formData file false "Book cover image"
// @Success 200 {object} BookResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	// Parse book ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Get existing book
	book, err := h.BookService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Handle file upload if provided
	file, err := c.FormFile("cover")
	if err == nil && file != nil {
		// Delete old cover if exists
		if book.Cover != "" {
			oldCoverPath := "./uploads/cover/" + book.Cover
			os.Remove(oldCoverPath)
		}

		// Generate unique filename
		coverFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		coverPath := "./uploads/cover/" + coverFileName
		if err := c.SaveUploadedFile(file, coverPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cover image"})
			return
		}
		book.Cover = coverFileName
	}

	// Update other fields if provided
	if title := c.PostForm("title"); title != "" {
		book.Title = title
	}
	if author := c.PostForm("author"); author != "" {
		book.Author = author
	}
	if stock := c.PostForm("stock"); stock != "" {
		stockNum, err := strconv.Atoi(stock)
		if err == nil && stockNum >= 0 {
			book.Stock = stockNum
		}
	}
	if isbn := c.PostForm("isbn"); isbn != "" {
		book.ISBN = isbn
	}
	if publishedYear := c.PostForm("published_year"); publishedYear != "" {
		year, err := strconv.Atoi(publishedYear)
		if err == nil && year > 0 {
			book.PublishedYear = year
		}
	}
	if category := c.PostForm("category"); category != "" {
		book.Category = category
	}
	if description := c.PostForm("description"); description != "" {
		book.Description = description
	}

	// Update book in database
	if err := h.BookService.Update(book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating book: " + err.Error()})
		return
	}

	// Return the updated book
	c.JSON(http.StatusOK, BookResponse{
		ID:            book.ID,
		Title:         book.Title,
		Author:        book.Author,
		Stock:         book.Stock,
		ISBN:          book.ISBN,
		PublishedYear: book.PublishedYear,
		Category:      book.Category,
		Description:   book.Description,
		Cover:         book.Cover,
	})
}

// DeleteBook handles deletion of a book
// @Summary Delete a book
// @Description Delete a book by its ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	// Parse book ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Delete book from database
	if err := h.BookService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return success with no content
	c.Status(http.StatusNoContent)
}
