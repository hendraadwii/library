package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// BorrowStatus represents the status of a book borrowing
type BorrowStatus string

const (
	StatusBorrowed BorrowStatus = "borrowed"
	StatusReturned BorrowStatus = "returned"
	StatusOverdue  BorrowStatus = "overdue"
)

// Borrowing represents a book borrowing record
type Borrowing struct {
	ID         int64        `json:"id"`
	BookID     int64        `json:"book_id"`
	UserID     int64        `json:"user_id"`
	Status     BorrowStatus `json:"status"`
	BorrowDate time.Time    `json:"borrow_date"`
	DueDate    time.Time    `json:"due_date"`
	ReturnDate *time.Time   `json:"return_date,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`

	// Joined fields
	BookTitle  string `json:"book_title,omitempty"`
	BookAuthor string `json:"book_author,omitempty"`
	UserName   string `json:"user_name,omitempty"`
	UserEmail  string `json:"user_email,omitempty"`
}

// BorrowingFilters represents filters for searching borrowings
type BorrowingFilters struct {
	BookID   int64
	UserID   int64
	Status   BorrowStatus
	Overdue  bool
	FromDate time.Time
	ToDate   time.Time
	Limit    int
	Offset   int
}

// MostBorrowedBook represents a book with its borrow count
type MostBorrowedBook struct {
	BookID      int64  `json:"book_id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	BorrowCount int    `json:"borrow_count"`
}

// BorrowingService handles borrowing-related database operations
type BorrowingService struct {
	DB *sql.DB
}

// NewBorrowingService creates a new BorrowingService
func NewBorrowingService(db *sql.DB) *BorrowingService {
	return &BorrowingService{DB: db}
}

// Create records a new book borrowing
func (s *BorrowingService) Create(borrowing *Borrowing) error {
	// Start a transaction to ensure consistency
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Check if book exists and has stock
	var stock int
	err = tx.QueryRow("SELECT stock FROM books WHERE id = ?", borrowing.BookID).Scan(&stock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("book not found with ID %d", borrowing.BookID)
		}
		return fmt.Errorf("error checking book stock: %w", err)
	}

	if stock <= 0 {
		return fmt.Errorf("book is out of stock")
	}

	// Check if user exists
	var userID int64
	err = tx.QueryRow("SELECT id FROM users WHERE id = ?", borrowing.UserID).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found with ID %d", borrowing.UserID)
		}
		return fmt.Errorf("error checking user: %w", err)
	}

	// Insert borrowing record
	query := `
		INSERT INTO borrowing_history 
		(book_id, user_id, status, borrow_date, due_date) 
		VALUES (?, ?, ?, ?, ?)
	`

	if borrowing.BorrowDate.IsZero() {
		borrowing.BorrowDate = time.Now()
	}

	if borrowing.Status == "" {
		borrowing.Status = StatusBorrowed
	}

	result, err := tx.Exec(
		query,
		borrowing.BookID,
		borrowing.UserID,
		borrowing.Status,
		borrowing.BorrowDate,
		borrowing.DueDate,
	)
	if err != nil {
		return fmt.Errorf("error creating borrowing record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert ID: %w", err)
	}
	borrowing.ID = id

	// Update book stock
	_, err = tx.Exec("UPDATE books SET stock = stock - 1 WHERE id = ?", borrowing.BookID)
	if err != nil {
		return fmt.Errorf("error updating book stock: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// GetByID retrieves a borrowing record by ID
func (s *BorrowingService) GetByID(id int64) (*Borrowing, error) {
	query := `
		SELECT 
			bh.id, bh.book_id, bh.user_id, bh.status, 
			bh.borrow_date, bh.due_date, bh.return_date, 
			bh.created_at, bh.updated_at,
			b.title, b.author, u.full_name, u.email
		FROM borrowing_history bh
		JOIN books b ON bh.book_id = b.id
		JOIN users u ON bh.user_id = u.id
		WHERE bh.id = ?
	`

	var borrowing Borrowing
	var returnDate sql.NullTime

	err := s.DB.QueryRow(query, id).Scan(
		&borrowing.ID,
		&borrowing.BookID,
		&borrowing.UserID,
		&borrowing.Status,
		&borrowing.BorrowDate,
		&borrowing.DueDate,
		&returnDate,
		&borrowing.CreatedAt,
		&borrowing.UpdatedAt,
		&borrowing.BookTitle,
		&borrowing.BookAuthor,
		&borrowing.UserName,
		&borrowing.UserEmail,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("borrowing record not found with ID %d", id)
		}
		return nil, fmt.Errorf("error getting borrowing record: %w", err)
	}

	if returnDate.Valid {
		borrowing.ReturnDate = &returnDate.Time
	}

	return &borrowing, nil
}

// ReturnBook marks a book as returned and updates the stock
func (s *BorrowingService) ReturnBook(id int64) error {
	// Start a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Get the borrowing record with book ID
	var bookID int64
	var status BorrowStatus
	err = tx.QueryRow(
		"SELECT book_id, status FROM borrowing_history WHERE id = ?",
		id,
	).Scan(&bookID, &status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("borrowing record not found with ID %d", id)
		}
		return fmt.Errorf("error getting borrowing record: %w", err)
	}

	// Check if already returned
	if status == StatusReturned {
		return fmt.Errorf("book already returned")
	}

	// Update the borrowing record
	_, err = tx.Exec(
		"UPDATE borrowing_history SET status = ?, return_date = ? WHERE id = ?",
		StatusReturned,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("error updating borrowing record: %w", err)
	}

	// Update book stock
	_, err = tx.Exec("UPDATE books SET stock = stock + 1 WHERE id = ?", bookID)
	if err != nil {
		return fmt.Errorf("error updating book stock: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// UpdateOverdueStatus checks for and updates overdue borrowings
func (s *BorrowingService) UpdateOverdueStatus() error {
	query := `
		UPDATE borrowing_history 
		SET status = ? 
		WHERE status = ? AND due_date < ? AND return_date IS NULL
	`

	_, err := s.DB.Exec(query, StatusOverdue, StatusBorrowed, time.Now())
	if err != nil {
		return fmt.Errorf("error updating overdue status: %w", err)
	}

	return nil
}

// List returns a list of borrowing records with pagination and filters
func (s *BorrowingService) List(filters BorrowingFilters) ([]Borrowing, error) {
	baseQuery := `
		SELECT 
			bh.id, bh.book_id, bh.user_id, bh.status, 
			bh.borrow_date, bh.due_date, bh.return_date, 
			bh.created_at, bh.updated_at,
			b.title, b.author, u.full_name, u.email
		FROM borrowing_history bh
		JOIN books b ON bh.book_id = b.id
		JOIN users u ON bh.user_id = u.id
		WHERE 1=1
	`

	var whereClause string
	var params []interface{}

	if filters.BookID > 0 {
		whereClause += " AND bh.book_id = ?"
		params = append(params, filters.BookID)
	}

	if filters.UserID > 0 {
		whereClause += " AND bh.user_id = ?"
		params = append(params, filters.UserID)
	}

	if filters.Status != "" {
		whereClause += " AND bh.status = ?"
		params = append(params, filters.Status)
	}

	if filters.Overdue {
		whereClause += " AND bh.due_date < ? AND (bh.status = ? OR bh.status = ?)"
		params = append(params, time.Now(), StatusBorrowed, StatusOverdue)
	}

	if !filters.FromDate.IsZero() {
		whereClause += " AND bh.borrow_date >= ?"
		params = append(params, filters.FromDate)
	}

	if !filters.ToDate.IsZero() {
		whereClause += " AND bh.borrow_date <= ?"
		params = append(params, filters.ToDate)
	}

	// Add pagination
	if filters.Limit <= 0 {
		filters.Limit = 10 // Default limit
	}

	if filters.Offset < 0 {
		filters.Offset = 0
	}

	query := baseQuery + whereClause + " ORDER BY bh.id DESC LIMIT ? OFFSET ?"
	params = append(params, filters.Limit, filters.Offset)

	// Execute query
	rows, err := s.DB.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("error listing borrowing records: %w", err)
	}
	defer rows.Close()

	// Parse results
	var borrowings []Borrowing
	for rows.Next() {
		var borrowing Borrowing
		var returnDate sql.NullTime

		err := rows.Scan(
			&borrowing.ID,
			&borrowing.BookID,
			&borrowing.UserID,
			&borrowing.Status,
			&borrowing.BorrowDate,
			&borrowing.DueDate,
			&returnDate,
			&borrowing.CreatedAt,
			&borrowing.UpdatedAt,
			&borrowing.BookTitle,
			&borrowing.BookAuthor,
			&borrowing.UserName,
			&borrowing.UserEmail,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning borrowing row: %w", err)
		}

		if returnDate.Valid {
			borrowing.ReturnDate = &returnDate.Time
		}

		borrowings = append(borrowings, borrowing)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating borrowing rows: %w", err)
	}

	return borrowings, nil
}

// Count returns the total number of borrowing records matching the filters
func (s *BorrowingService) Count(filters BorrowingFilters) (int, error) {
	baseQuery := "SELECT COUNT(*) FROM borrowing_history bh WHERE 1=1"

	var whereClause string
	var params []interface{}

	if filters.BookID > 0 {
		whereClause += " AND bh.book_id = ?"
		params = append(params, filters.BookID)
	}

	if filters.UserID > 0 {
		whereClause += " AND bh.user_id = ?"
		params = append(params, filters.UserID)
	}

	if filters.Status != "" {
		whereClause += " AND bh.status = ?"
		params = append(params, filters.Status)
	}

	if filters.Overdue {
		whereClause += " AND bh.due_date < ? AND (bh.status = ? OR bh.status = ?)"
		params = append(params, time.Now(), StatusBorrowed, StatusOverdue)
	}

	if !filters.FromDate.IsZero() {
		whereClause += " AND bh.borrow_date >= ?"
		params = append(params, filters.FromDate)
	}

	if !filters.ToDate.IsZero() {
		whereClause += " AND bh.borrow_date <= ?"
		params = append(params, filters.ToDate)
	}

	query := baseQuery + whereClause

	var count int
	err := s.DB.QueryRow(query, params...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting borrowing records: %w", err)
	}

	return count, nil
}

// GetMostBorrowedBooks returns a list of the most borrowed books
func (s *BorrowingService) GetMostBorrowedBooks(limit int) ([]MostBorrowedBook, error) {
	if limit <= 0 {
		limit = 10
	}

	query := `
		SELECT b.id, b.title, b.author, COUNT(bh.id) as borrow_count 
		FROM books b
		JOIN borrowing_history bh ON b.id = bh.book_id
		GROUP BY b.id, b.title, b.author
		ORDER BY borrow_count DESC
		LIMIT ?
	`

	rows, err := s.DB.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("error getting most borrowed books: %w", err)
	}
	defer rows.Close()

	var result []MostBorrowedBook
	for rows.Next() {
		var book MostBorrowedBook
		if err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.BorrowCount); err != nil {
			return nil, fmt.Errorf("error scanning most borrowed book row: %w", err)
		}
		result = append(result, book)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating most borrowed book rows: %w", err)
	}

	return result, nil
}
