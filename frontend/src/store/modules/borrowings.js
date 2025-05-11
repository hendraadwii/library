import { getAPI } from '../../utils/api'

export default {
  namespaced: true,
  state: {
    borrowings: [],
    borrowing: null,
    totalBorrowings: 0,
    currentPage: 1,
    perPage: 10,
    totalPages: 0,
    overdueBorrowings: []
  },
  mutations: {
    SET_BORROWINGS(state, { borrowings, totalCount, page, perPage, totalPages }) {
      state.borrowings = borrowings
      state.totalBorrowings = totalCount
      state.currentPage = page
      state.perPage = perPage
      state.totalPages = totalPages
    },
    SET_BORROWING(state, borrowing) {
      state.borrowing = borrowing
    },
    ADD_BORROWING(state, borrowing) {
      state.borrowings.unshift(borrowing)
    },
    UPDATE_BORROWING(state, updatedBorrowing) {
      const index = state.borrowings.findIndex(b => b.id === updatedBorrowing.id)
      if (index !== -1) {
        state.borrowings.splice(index, 1, updatedBorrowing)
      }
      if (state.borrowing && state.borrowing.id === updatedBorrowing.id) {
        state.borrowing = updatedBorrowing
      }
    },
    SET_OVERDUE(state, borrowings) {
      state.overdueBorrowings = borrowings
    }
  },
  actions: {
    // Get borrowings with pagination and filters
    async getBorrowings({ commit }, { page = 1, perPage = 10, status = '', bookId = null, userId = null, overdue = false } = {}) {
      try {
        const api = getAPI()
        const response = await api.get('/borrowings', {
          params: { 
            page, 
            per_page: perPage, 
            status, 
            book_id: bookId, 
            user_id: userId, 
            overdue 
          }
        })
        
        commit('SET_BORROWINGS', {
          borrowings: response.data.borrowings,
          totalCount: response.data.total_count,
          page: response.data.page,
          perPage: response.data.per_page,
          totalPages: response.data.total_pages
        })
        
        return response.data
      } catch (error) {
        console.error('Error fetching borrowings:', error)
        throw error
      }
    },
    
    // Get a single borrowing by ID
    async getBorrowing({ commit }, borrowingId) {
      try {
        const api = getAPI()
        const response = await api.get(`/borrowings/${borrowingId}`)
        commit('SET_BORROWING', response.data)
        return response.data
      } catch (error) {
        console.error(`Error fetching borrowing ${borrowingId}:`, error)
        throw error
      }
    },
    
    // Borrow a book
    async borrowBook({ commit }, borrowData) {
      try {
        const api = getAPI()
        const response = await api.post('/borrowings', borrowData)
        commit('ADD_BORROWING', response.data)
        return response.data
      } catch (error) {
        console.error('Error borrowing book:', error)
        throw error
      }
    },
    
    // Return a book
    async returnBook({ commit }, borrowingId) {
      try {
        const api = getAPI()
        const response = await api.post(`/borrowings/${borrowingId}/return`)
        commit('UPDATE_BORROWING', response.data)
        return response.data
      } catch (error) {
        console.error(`Error returning book for borrowing ${borrowingId}:`, error)
        throw error
      }
    },
    
    // Get overdue borrowings
    async getOverdueBorrowings({ commit }, { page = 1, perPage = 10 } = {}) {
      try {
        const api = getAPI()
        const response = await api.get('/borrowings/overdue', {
          params: { page, per_page: perPage }
        })
        commit('SET_OVERDUE', response.data.borrowings)
        
        // Also update pagination data
        commit('SET_BORROWINGS', {
          borrowings: response.data.borrowings,
          totalCount: response.data.total_count,
          page: response.data.page,
          perPage: response.data.per_page,
          totalPages: response.data.total_pages
        })
        
        return response.data
      } catch (error) {
        console.error('Error fetching overdue borrowings:', error)
        throw error
      }
    }
  },
  getters: {
    allBorrowings: state => state.borrowings,
    borrowingById: state => id => state.borrowings.find(b => b.id === id),
    currentBorrowing: state => state.borrowing,
    totalBorrowings: state => state.totalBorrowings,
    currentPage: state => state.currentPage,
    perPage: state => state.perPage,
    totalPages: state => state.totalPages,
    overdueBorrowings: state => state.overdueBorrowings,
    
    // Get borrowings by status
    borrowingsByStatus: state => status => {
      return state.borrowings.filter(b => b.status === status)
    },
    
    // Get borrowings by book ID
    borrowingsByBook: state => bookId => {
      return state.borrowings.filter(b => b.book_id === bookId)
    }
  }
} 