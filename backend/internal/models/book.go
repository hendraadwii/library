package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Book represents a book in the library
type Book struct {
	ID            int64     `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Stock         int       `json:"stock"`
	ISBN          string    `json:"isbn,omitempty"`
	PublishedYear int       `json:"published_year,omitempty"`
	Category      string    `json:"category,omitempty"`
	Description   string    `json:"description,omitempty"`
	Cover         string    `json:"cover,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BookFilters represents filters for searching books
type BookFilters struct {
	Title    string
	Author   string
	Category string
	Limit    int
	Offset   int
}

// BookService handles book-related database operations
type BookService struct {
	DB *sql.DB
}

// NewBookService creates a new BookService
func NewBookService(db *sql.DB) *BookService {
	return &BookService{DB: db}
}

// Create adds a new book to the database
func (s *BookService) Create(book *Book) error {
	query := `
		INSERT INTO books 
		(title, author, stock, isbn, published_year, category, description, cover) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.DB.Exec(
		query,
		book.Title,
		book.Author,
		book.Stock,
		book.ISBN,
		book.PublishedYear,
		book.Category,
		book.Description,
		book.Cover,
	)
	if err != nil {
		return fmt.Errorf("error creating book: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}

	book.ID = id
	return nil
}

// GetByID retrieves a book by its ID
func (s *BookService) GetByID(id int64) (*Book, error) {
	query := `SELECT id, title, author, stock, isbn, published_year, category, description, COALESCE(cover, ''), created_at, updated_at 
              FROM books WHERE id = ?`

	var book Book
	err := s.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.Stock,
		&book.ISBN,
		&book.PublishedYear,
		&book.Category,
		&book.Description,
		&book.Cover,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("book not found with ID %d", id)
		}
		return nil, fmt.Errorf("error getting book: %w", err)
	}

	return &book, nil
}

// Update updates an existing book
func (s *BookService) Update(book *Book) error {
	query := `
		UPDATE books 
		SET title = ?, author = ?, stock = ?, isbn = ?, 
		    published_year = ?, category = ?, description = ? 
		WHERE id = ?
	`

	_, err := s.DB.Exec(
		query,
		book.Title,
		book.Author,
		book.Stock,
		book.ISBN,
		book.PublishedYear,
		book.Category,
		book.Description,
		book.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating book: %w", err)
	}

	return nil
}

// Delete removes a book by its ID
func (s *BookService) Delete(id int64) error {
	query := "DELETE FROM books WHERE id = ?"

	result, err := s.DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found with ID %d", id)
	}

	return nil
}

// List returns a list of books with pagination and filters
func (s *BookService) List(filters BookFilters) ([]Book, error) {
	// Start with base query
	baseQuery := `
		SELECT id, title, author, stock, isbn, published_year, category, description, COALESCE(cover, ''), created_at, updated_at 
		FROM books
		WHERE 1=1
	`

	// Build where clause and params
	var whereClause string
	var params []interface{}

	if filters.Title != "" {
		whereClause += " AND title LIKE ?"
		params = append(params, "%"+filters.Title+"%")
	}

	if filters.Author != "" {
		whereClause += " AND author LIKE ?"
		params = append(params, "%"+filters.Author+"%")
	}

	if filters.Category != "" {
		whereClause += " AND category = ?"
		params = append(params, filters.Category)
	}

	// Add pagination
	if filters.Limit <= 0 {
		filters.Limit = 10 // Default limit
	}

	if filters.Offset < 0 {
		filters.Offset = 0
	}

	query := baseQuery + whereClause + " ORDER BY id DESC LIMIT ? OFFSET ?"
	params = append(params, filters.Limit, filters.Offset)

	// Execute query
	rows, err := s.DB.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error listing books: %w", err)
	}
	defer rows.Close()

	// Parse results
	var books []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.Stock,
			&book.ISBN,
			&book.PublishedYear,
			&book.Category,
			&book.Description,
			&book.Cover,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning book row: %w", err)
		}
		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating book rows: %w", err)
	}

	return books, nil
}

// UpdateStock updates the stock of a book by its ID
func (s *BookService) UpdateStock(id int64, change int) error {
	query := "UPDATE books SET stock = stock + ? WHERE id = ?"

	result, err := s.DB.Exec(query, change, id)
	if err != nil {
		return fmt.Errorf("error updating book stock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found with ID %d", id)
	}

	return nil
}

// Count returns the total number of books matching the filters
func (s *BookService) Count(filters BookFilters) (int, error) {
	// Start with base query
	baseQuery := "SELECT COUNT(*) FROM books WHERE 1=1"

	// Build where clause and params
	var whereClause string
	var params []interface{}

	if filters.Title != "" {
		whereClause += " AND title LIKE ?"
		params = append(params, "%"+filters.Title+"%")
	}

	if filters.Author != "" {
		whereClause += " AND author LIKE ?"
		params = append(params, "%"+filters.Author+"%")
	}

	if filters.Category != "" {
		whereClause += " AND category = ?"
		params = append(params, filters.Category)
	}

	query := baseQuery + whereClause

	// Execute query
	var count int
	err := s.DB.QueryRow(query, params...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting books: %w", err)
	}

	return count, nil
}

// ISBNExists checks if a book with the given ISBN already exists
func (s *BookService) ISBNExists(isbn string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM books WHERE isbn = ?)"
	err := s.DB.QueryRow(query, isbn).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking ISBN: %w", err)
	}
	return exists, nil
}
