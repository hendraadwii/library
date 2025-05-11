package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/aksaaaraa/library/internal/models"
	"github.com/gin-gonic/gin"
)

// BorrowingHandler handles borrowing-related requests
type BorrowingHandler struct {
	BorrowingService *models.BorrowingService
	BookService      *models.BookService
}

// BorrowRequest represents a request to borrow a book
type BorrowRequest struct {
	BookID   int64     `json:"book_id" binding:"required"`
	DueDate  time.Time `json:"due_date" binding:"required"`
}

// BorrowingResponse represents a borrowing response
type BorrowingResponse struct {
	ID         int64        `json:"id"`
	BookID     int64        `json:"book_id"`
	UserID     int64        `json:"user_id"`
	Status     string       `json:"status"`
	BorrowDate time.Time    `json:"borrow_date"`
	DueDate    time.Time    `json:"due_date"`
	ReturnDate *time.Time   `json:"return_date,omitempty"`
	BookTitle  string       `json:"book_title,omitempty"`
	BookAuthor string       `json:"book_author,omitempty"`
	UserName   string       `json:"user_name,omitempty"`
	UserEmail  string       `json:"user_email,omitempty"`
}

// PaginatedBorrowingsResponse represents a paginated list of borrowings
type PaginatedBorrowingsResponse struct {
	Borrowings []BorrowingResponse `json:"borrowings"`
	TotalCount int                 `json:"total_count"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	TotalPages int                 `json:"total_pages"`
}

// NewBorrowingHandler creates a new BorrowingHandler
func NewBorrowingHandler(borrowingService *models.BorrowingService, bookService *models.BookService) *BorrowingHandler {
	return &BorrowingHandler{
		BorrowingService: borrowingService,
		BookService:      bookService,
	}
}

// GetBorrowings handles retrieval of all borrowings with pagination
// @Summary Get all borrowings
// @Description Get a paginated list of borrowings with optional filters
// @Tags borrowings
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Param status query string false "Filter by status (borrowed, returned, overdue)"
// @Param book_id query int false "Filter by book ID"
// @Param user_id query int false "Filter by user ID"
// @Param overdue query bool false "Filter by overdue status"
// @Success 200 {object} PaginatedBorrowingsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /borrowings [get]
func (h *BorrowingHandler) GetBorrowings(c *gin.Context) {
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
	bookID, _ := strconv.ParseInt(c.Query("book_id"), 10, 64)
	userID, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	status := models.BorrowStatus(c.Query("status"))
	overdue, _ := strconv.ParseBool(c.Query("overdue"))

	// Create filters
	filters := models.BorrowingFilters{
		BookID:  bookID,
		UserID:  userID,
		Status:  status,
		Overdue: overdue,
		Limit:   perPage,
		Offset:  offset,
	}

	// Get borrowings from database
	borrowings, err := h.BorrowingService.List(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting borrowings: " + err.Error()})
		return
	}

	// Get total count for pagination
	totalCount, err := h.BorrowingService.Count(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting borrowings: " + err.Error()})
		return
	}

	// Calculate total pages
	totalPages := totalCount / perPage
	if totalCount%perPage > 0 {
		totalPages++
	}

	// Map to response format
	var borrowingResponses []BorrowingResponse
	for _, borrowing := range borrowings {
		borrowingResponses = append(borrowingResponses, BorrowingResponse{
			ID:         borrowing.ID,
			BookID:     borrowing.BookID,
			UserID:     borrowing.UserID,
			Status:     string(borrowing.Status),
			BorrowDate: borrowing.BorrowDate,
			DueDate:    borrowing.DueDate,
			ReturnDate: borrowing.ReturnDate,
			BookTitle:  borrowing.BookTitle,
			BookAuthor: borrowing.BookAuthor,
			UserName:   borrowing.UserName,
			UserEmail:  borrowing.UserEmail,
		})
	}

	c.JSON(http.StatusOK, PaginatedBorrowingsResponse{
		Borrowings: borrowingResponses,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// GetBorrowing handles retrieval of a specific borrowing by ID
// @Summary Get a borrowing
// @Description Get a borrowing by its ID
// @Tags borrowings
// @Accept json
// @Produce json
// @Param id path int true "Borrowing ID"
// @Success 200 {object} BorrowingResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /borrowings/{id} [get]
func (h *BorrowingHandler) GetBorrowing(c *gin.Context) {
	// Parse borrowing ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid borrowing ID"})
		return
	}

	// Get borrowing from database
	borrowing, err := h.BorrowingService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Return the borrowing
	c.JSON(http.StatusOK, BorrowingResponse{
		ID:         borrowing.ID,
		BookID:     borrowing.BookID,
		UserID:     borrowing.UserID,
		Status:     string(borrowing.Status),
		BorrowDate: borrowing.BorrowDate,
		DueDate:    borrowing.DueDate,
		ReturnDate: borrowing.ReturnDate,
		BookTitle:  borrowing.BookTitle,
		BookAuthor: borrowing.BookAuthor,
		UserName:   borrowing.UserName,
		UserEmail:  borrowing.UserEmail,
	})
}

// BorrowBook handles borrowing a book
// @Summary Borrow a book
// @Description Borrow a book for the authenticated user
// @Tags borrowings
// @Accept json
// @Produce json
// @Param borrowing body BorrowRequest true "Borrowing details"
// @Success 201 {object} BorrowingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse "Book not found"
// @Failure 409 {object} ErrorResponse "Book out of stock"
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /borrowings [post]
func (h *BorrowingHandler) BorrowBook(c *gin.Context) {
	var req BorrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Validate due date is in the future
	if req.DueDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Due date must be in the future"})
		return
	}

	// Create borrowing record
	borrowing := &models.Borrowing{
		BookID:     req.BookID,
		UserID:     userID.(int64),
		BorrowDate: time.Now(),
		DueDate:    req.DueDate,
		Status:     models.StatusBorrowed,
	}

	// Save to database (which will also check stock and update book stock)
	if err := h.BorrowingService.Create(borrowing); err != nil {
		// Check if error is about stock
		if err.Error() == "book is out of stock" {
			c.JSON(http.StatusConflict, gin.H{"error": "Book is out of stock"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error borrowing book: " + err.Error()})
		return
	}

	// Fetch the complete borrowing record with book and user details
	completeBorrowing, err := h.BorrowingService.GetByID(borrowing.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving borrowing details: " + err.Error()})
		return
	}

	// Return the borrowing record
	c.JSON(http.StatusCreated, BorrowingResponse{
		ID:         completeBorrowing.ID,
		BookID:     completeBorrowing.BookID,
		UserID:     completeBorrowing.UserID,
		Status:     string(completeBorrowing.Status),
		BorrowDate: completeBorrowing.BorrowDate,
		DueDate:    completeBorrowing.DueDate,
		ReturnDate: completeBorrowing.ReturnDate,
		BookTitle:  completeBorrowing.BookTitle,
		BookAuthor: completeBorrowing.BookAuthor,
		UserName:   completeBorrowing.UserName,
		UserEmail:  completeBorrowing.UserEmail,
	})
}

// ReturnBook handles returning a borrowed book
// @Summary Return a book
// @Description Return a previously borrowed book
// @Tags borrowings
// @Accept json
// @Produce json
// @Param id path int true "Borrowing ID"
// @Success 200 {object} BorrowingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Book already returned"
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /borrowings/{id}/return [post]
func (h *BorrowingHandler) ReturnBook(c *gin.Context) {
	// Parse borrowing ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid borrowing ID"})
		return
	}

	// Return the book
	if err := h.BorrowingService.ReturnBook(id); err != nil {
		// Check if already returned
		if err.Error() == "book already returned" {
			c.JSON(http.StatusConflict, gin.H{"error": "Book already returned"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error returning book: " + err.Error()})
		return
	}

	// Get updated borrowing record
	borrowing, err := h.BorrowingService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving borrowing details: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, BorrowingResponse{
		ID:         borrowing.ID,
		BookID:     borrowing.BookID,
		UserID:     borrowing.UserID,
		Status:     string(borrowing.Status),
		BorrowDate: borrowing.BorrowDate,
		DueDate:    borrowing.DueDate,
		ReturnDate: borrowing.ReturnDate,
		BookTitle:  borrowing.BookTitle,
		BookAuthor: borrowing.BookAuthor,
		UserName:   borrowing.UserName,
		UserEmail:  borrowing.UserEmail,
	})
}

// GetOverdueBooks gets all overdue books
// @Summary Get overdue books
// @Description Get all books that are currently overdue
// @Tags borrowings
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Success 200 {object} PaginatedBorrowingsResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /borrowings/overdue [get]
func (h *BorrowingHandler) GetOverdueBooks(c *gin.Context) {
	// Update overdue status
	if err := h.BorrowingService.UpdateOverdueStatus(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating overdue status: " + err.Error()})
		return
	}

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

	// Create overdue filter
	filters := models.BorrowingFilters{
		Overdue: true,
		Limit:   perPage,
		Offset:  offset,
	}

	// Get borrowings from database
	borrowings, err := h.BorrowingService.List(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting overdue books: " + err.Error()})
		return
	}

	// Get total count for pagination
	totalCount, err := h.BorrowingService.Count(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting overdue books: " + err.Error()})
		return
	}

	// Calculate total pages
	totalPages := totalCount / perPage
	if totalCount%perPage > 0 {
		totalPages++
	}

	// Map to response format
	var borrowingResponses []BorrowingResponse
	for _, borrowing := range borrowings {
		borrowingResponses = append(borrowingResponses, BorrowingResponse{
			ID:         borrowing.ID,
			BookID:     borrowing.BookID,
			UserID:     borrowing.UserID,
			Status:     string(borrowing.Status),
			BorrowDate: borrowing.BorrowDate,
			DueDate:    borrowing.DueDate,
			ReturnDate: borrowing.ReturnDate,
			BookTitle:  borrowing.BookTitle,
			BookAuthor: borrowing.BookAuthor,
			UserName:   borrowing.UserName,
			UserEmail:  borrowing.UserEmail,
		})
	}

	c.JSON(http.StatusOK, PaginatedBorrowingsResponse{
		Borrowings: borrowingResponses,
		TotalCount: totalCount,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// GetMostBorrowedBooks gets the most borrowed books
// @Summary Get most borrowed books
// @Description Get a list of the most borrowed books
// @Tags books
// @Accept json
// @Produce json
// @Param limit query int false "Limit (default: 10)"
// @Success 200 {array} object
// @Failure 500 {object} ErrorResponse
// @Router /books/most-borrowed [get]
func (h *BorrowingHandler) GetMostBorrowedBooks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	books, err := h.BorrowingService.GetMostBorrowedBooks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting most borrowed books: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
} 