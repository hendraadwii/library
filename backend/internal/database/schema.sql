-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS library_db;
USE library_db;

-- Create books table
CREATE TABLE IF NOT EXISTS books (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    isbn VARCHAR(50) UNIQUE,
    published_year INT,
    category VARCHAR(100),
    description TEXT,
    cover VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role ENUM('admin', 'member') DEFAULT 'member',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Create borrowing_history table
CREATE TABLE IF NOT EXISTS borrowing_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    book_id INT NOT NULL,
    user_id INT NOT NULL,
    status ENUM('borrowed', 'returned', 'overdue') DEFAULT 'borrowed',
    borrow_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (book_id) REFERENCES books(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Add indexes for performance
CREATE INDEX idx_books_title ON books(title);
CREATE INDEX idx_books_author ON books(author);
CREATE INDEX idx_borrowing_status ON borrowing_history(status);
CREATE INDEX idx_borrowing_dates ON borrowing_history(borrow_date, due_date, return_date);

-- Example parameterized query for book insertion (for documentation)
-- INSERT INTO books (title, author, stock, isbn, published_year, category, description) 
-- VALUES (?, ?, ?, ?, ?, ?, ?);

-- Example parameterized query for selecting books with pagination
-- SELECT * FROM books LIMIT ? OFFSET ?;

-- Example query for most borrowed books
-- SELECT b.id, b.title, b.author, COUNT(bh.id) as borrow_count 
-- FROM books b
-- JOIN borrowing_history bh ON b.id = bh.book_id
-- GROUP BY b.id, b.title, b.author
-- ORDER BY borrow_count DESC
-- LIMIT 10; 