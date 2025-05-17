import { getAPI } from '../../utils/api'

export default {
  namespaced: true,
  state: {
    books: [],
    book: null,
    totalBooks: 0,
    currentPage: 1,
    perPage: 10,
    totalPages: 0,
    mostBorrowed: []
  },
  mutations: {
    SET_BOOKS(state, { books, totalCount, page, perPage, totalPages }) {
      state.books = Array.isArray(books) ? books : []
      state.totalBooks = totalCount
      state.currentPage = page
      state.perPage = perPage
      state.totalPages = totalPages
    },
    SET_BOOK(state, book) {
      state.book = book
    },
    ADD_BOOK(state, book) {
      state.books.unshift(book)
    },
    UPDATE_BOOK(state, updatedBook) {
      const index = state.books.findIndex(book => book.id === updatedBook.id)
      if (index !== -1) {
        state.books.splice(index, 1, updatedBook)
      }
      if (state.book && state.book.id === updatedBook.id) {
        state.book = updatedBook
      }
    },
    REMOVE_BOOK(state, bookId) {
      state.books = state.books.filter(book => book.id !== bookId)
      if (state.book && state.book.id === bookId) {
        state.book = null
      }
    },
    SET_MOST_BORROWED(state, books) {
      state.mostBorrowed = books
    }
  },
  actions: {
    // Get books with pagination and filters
    async getBooks({ commit }, { page = 1, perPage = 10, title = '', author = '', category = '' } = {}) {
      try {
        const api = getAPI()
        const response = await api.get('/books', {
          params: { page, per_page: perPage, title, author, category }
        })
        
        commit('SET_BOOKS', {
          books: response.data.books,
          totalCount: response.data.total_count,
          page: response.data.page,
          perPage: response.data.per_page,
          totalPages: response.data.total_pages
        })
        
        return response.data
      } catch (error) {
        console.error('Error fetching books:', error)
        throw error
      }
    },
    
    // Get a single book by ID
    async getBook({ commit }, bookId) {
      try {
        const api = getAPI()
        const response = await api.get(`/books/${bookId}`)
        commit('SET_BOOK', response.data)
        return response.data
      } catch (error) {
        console.error(`Error fetching book ${bookId}:`, error)
        throw error
      }
    },
    
    // Create a new book
    async createBook({ commit }, bookData) {
      try {
        const api = getAPI()
        const response = await api.post('/books', bookData)
        commit('ADD_BOOK', response.data)
        return response.data
      } catch (error) {
        console.error('Error creating book:', error)
        throw error
      }
    },
    
    // Update an existing book
    async updateBook({ commit }, { bookId, bookData }) {
      try {
        const api = getAPI()
        let response

        // Check if bookData is FormData (has file)
        if (bookData instanceof FormData) {
          response = await api.put(`/books/${bookId}`, bookData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
          })
        } else {
          response = await api.put(`/books/${bookId}`, bookData)
        }

        commit('UPDATE_BOOK', response.data)
        return response.data
      } catch (error) {
        console.error(`Error updating book ${bookId}:`, error)
        throw error
      }
    },
    
    // Delete a book
    async deleteBook({ commit }, bookId) {
      try {
        const api = getAPI()
        await api.delete(`/books/${bookId}`)
        commit('REMOVE_BOOK', bookId)
      } catch (error) {
        console.error(`Error deleting book ${bookId}:`, error)
        throw error
      }
    },
    
    // Get most borrowed books
    async getMostBorrowedBooks({ commit }, limit = 10) {
      try {
        const api = getAPI()
        const response = await api.get('/books/most-borrowed', {
          params: { limit }
        })
        commit('SET_MOST_BORROWED', response.data)
        return response.data
      } catch (error) {
        console.error('Error fetching most borrowed books:', error)
        throw error
      }
    }
  },
  getters: {
    allBooks: state => Array.isArray(state.books) ? state.books : [],
    bookById: state => id => state.books.find(book => book.id === id),
    currentBook: state => state.book,
    totalBooks: state => state.totalBooks,
    currentPage: state => state.currentPage,
    perPage: state => state.perPage,
    totalPages: state => state.totalPages,
    mostBorrowedBooks: state => state.mostBorrowed
  }
} 