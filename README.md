Library - Indoseiki

A fullstack Library Management System built with Vue.js (frontend) and Go (backend).

---

## Project Structure

```
├── backend/
│   ├── cmd/
│   │   ├── api/
│   │   │   ├── main.go                 # Main entry point
│   │   │   └── handlers/               # HTTP request handlers
│   │   │       ├── auth_handler.go     # Authentication endpoints
│   │   │       ├── book_handler.go     # Book management endpoints
│   │   │       ├── borrowing_handler.go # Borrowing management endpoints
│   │   │       └── user_handler.go     # User management endpoints
│   │   └── frontend/
│   │       └── public/
│   ├── config/                         # Configuration files
│   ├── internal/
│   │   ├── auth/                       # Authentication logic
│   │   │   ├── jwt.go                  # JWT token handling
│   │   │   ├── middleware.go           # Auth middleware
│   │   │   └── ratelimit.go            # Rate limiting
│   │   ├── config/
│   │   │   └── config.go               # App configuration
│   │   ├── database/
│   │   │   ├── db.go                   # Database connection
│   │   │   └── schema.sql              # Database schema
│   │   ├── models/                     # Data models
│   │   │   ├── book.go                 # Book model
│   │   │   ├── borrowing.go            # Borrowing model
│   │   │   └── user.go                 # User model
│   │   └── server/
│   │       └── server.go               # Server setup & routes
│   ├── uploads/                        # File uploads directory
│   │   └── cover/                      # Book cover images
│   ├── go.mod                          # Go module definition
│   ├── go.sum                          # Go module checksums
│   └── README.md                       # Backend documentation
├── frontend/
│   ├── public/                         # Static files
│   │   ├── favicon.png                 # Site favicon
│   │   └── index.html                  # HTML entry point
│   ├── src/
│   │   ├── assets/                     # Static assets
│   │   │   └── tailwind.css            # Tailwind styles
│   │   ├── components/                 # Reusable components
│   │   ├── router/
│   │   │   └── index.js                # Vue Router setup
│   │   ├── store/                      # Vuex store
│   │   │   ├── index.js                # Store configuration
│   │   │   └── modules/                # Store modules
│   │   │       ├── auth.js             # Auth state management
│   │   │       ├── books.js            # Books state management
│   │   │       ├── borrowings.js       # Borrowings state management
│   │   │       └── users.js            # Users state management
│   │   ├── utils/                      # Utility functions
│   │   │   ├── api.js                  # API client setup
│   │   │   └── sessionTimeout.js       # Session management
│   │   ├── views/                      # Page components
│   │   │   ├── BookDetails.vue         # Book details page
│   │   │   ├── Books.vue              # Books list page
│   │   │   ├── Borrowings.vue         # Borrowings management
│   │   │   ├── ForgotPassword.vue     # Password recovery
│   │   │   ├── Home.vue               # Home page
│   │   │   ├── Login.vue              # Login page
│   │   │   ├── MemberDashboard.vue    # Member dashboard
│   │   │   ├── NotFound.vue           # 404 page
│   │   │   ├── Register.vue           # Registration page
│   │   │   └── Users.vue              # User management
│   │   ├── App.vue                     # Root component
│   │   └── main.js                     # App entry point
│   ├── .env.example                    # Environment variables template
│   ├── .eslintrc.js                    # ESLint configuration
│   ├── babel.config.js                 # Babel configuration
│   ├── package.json                    # NPM package definition
│   ├── package-lock.json               # NPM package lock
│   ├── postcss.config.js               # PostCSS configuration
│   └── tailwind.config.js              # Tailwind CSS configuration
├── .gitignore                          # Git ignore rules
└── README.md                           # Project documentation
```

## Features

- User authentication (register, login, forgot password)
- Role-based dashboard (admin & member)
- Book management (CRUD)
  - Book cover image upload and display
  - ISBN validation (10 digits)
  - Author name validation (letters only)
  - Category selection
  - Stock management
- Borrowing management
  - Due date selection
  - Return tracking
  - Borrowing history
- Responsive, modern UI with Tailwind CSS
- Animated transitions and feedback
- File upload handling
  - Book cover images stored in `backend/uploads/cover/`
  - Images served via `/static/cover/` endpoint
  - Automatic image resizing and optimization

## Getting Started

### Prerequisites
- Node.js version 16 & npm (for frontend)
- Go (for backend)

### Backend Setup
1. `cd backend`
2. Copy `.env.example` to `.env` and configure:
   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=library
   JWT_SECRET=your_jwt_secret
   SERVER_PORT=8081
   ```
3. Install Go modules:
   ```sh
   go mod tidy
   ```
4. Run the backend:
   ```sh
   go run ./cmd/api/main.go
   ```
   The backend will be available at http://localhost:8081

### Frontend Setup
1. `cd frontend`
2. Copy `.env.example` to `.env` and configure:
   ```
   VUE_APP_API_URL=http://localhost:8081
   ```
3. Install dependencies:
   ```sh
   npm install
   ```
4. Run the frontend:
   ```sh
   npm run serve
   ```
   The frontend will be available at http://localhost:8080

## File Upload System

The application includes a file upload system for book cover images:

1. **Storage Location**
   - Images are stored in `backend/uploads/cover/`
   - Each image is saved with a unique filename
   - The filename is stored in the book's `cover` field in the database

2. **Access**
   - Images are served via `/static/cover/` endpoint
   - Full URL format: `http://localhost:8081/static/cover/filename.jpg`
   - Images are publicly accessible but stored outside the web root

3. **Frontend Integration**
   - Book list shows thumbnails (16x12 pixels)
   - Book details shows full-size cover
   - Placeholder shown when no cover is available
   - File input accepts common image formats

## Customization

- Change web title and favicon: edit `frontend/public/index.html`
- Modify API endpoints: edit `frontend/src/utils/api.js`
- Adjust image upload settings: edit `backend/cmd/api/handlers/book_handler.go`
- Customize UI: edit `frontend/src/assets/tailwind.css`

## License

MIT 