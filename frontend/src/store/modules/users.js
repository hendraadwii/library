import { getAPI } from '../../utils/api'

const state = {
  users: [],
  totalUsers: 0,
  currentPage: 1,
  perPage: 10,
  totalPages: 0,
  loading: false,
  error: null
}

const mutations = {
  SET_USERS(state, { users, totalCount, page, perPage, totalPages }) {
    state.users = users
    state.totalUsers = totalCount
    state.currentPage = page
    state.perPage = perPage
    state.totalPages = totalPages
  },
  SET_USERS_LOADING(state, loading) {
    state.loading = loading
  },
  SET_USERS_ERROR(state, error) {
    state.error = error
  },
  ADD_USER(state, user) {
    state.users.unshift(user)
    state.totalUsers++
  },
  UPDATE_USER(state, updatedUser) {
    const index = state.users.findIndex(user => user.id === updatedUser.id)
    if (index !== -1) {
      state.users.splice(index, 1, updatedUser)
    }
  },
  REMOVE_USER(state, userId) {
    state.users = state.users.filter(user => user.id !== userId)
    state.totalUsers--
  }
}

const actions = {
  // Get users with pagination and filters
  async getUsers({ commit }, { page = 1, perPage = 10, search = '', role = '', status = '' } = {}) {
    commit('SET_USERS_LOADING', true)
    commit('SET_USERS_ERROR', null)
    try {
      const api = getAPI()
      const response = await api.get('/users', {
        params: { page, per_page: perPage, search, role, status }
      })
      
      commit('SET_USERS', {
        users: response.data.users,
        totalCount: response.data.total_count,
        page: response.data.page,
        perPage: response.data.per_page,
        totalPages: response.data.total_pages
      })
      
      return response.data
    } catch (error) {
      commit('SET_USERS_ERROR', error.response?.data?.error || error.message)
      throw error
    } finally {
      commit('SET_USERS_LOADING', false)
    }
  },

  // Create new user
  async createUser({ commit }, userData) {
    commit('SET_USERS_LOADING', true)
    commit('SET_USERS_ERROR', null)
    try {
      const api = getAPI()
      const response = await api.post('/users', userData)
      commit('ADD_USER', response.data)
      return response.data
    } catch (error) {
      commit('SET_USERS_ERROR', error.response?.data?.error || error.message)
      throw error
    } finally {
      commit('SET_USERS_LOADING', false)
    }
  },

  // Update user
  async updateUser({ commit }, { id, ...userData }) {
    commit('SET_USERS_LOADING', true)
    commit('SET_USERS_ERROR', null)
    try {
      const api = getAPI()
      const response = await api.put(`/users/${id}`, userData)
      commit('UPDATE_USER', response.data)
      return response.data
    } catch (error) {
      commit('SET_USERS_ERROR', error.response?.data?.error || error.message)
      throw error
    } finally {
      commit('SET_USERS_LOADING', false)
    }
  },

  // Delete user
  async deleteUser({ commit }, userId) {
    commit('SET_USERS_LOADING', true)
    commit('SET_USERS_ERROR', null)
    try {
      const api = getAPI()
      await api.delete(`/users/${userId}`)
      commit('REMOVE_USER', userId)
    } catch (error) {
      commit('SET_USERS_ERROR', error.response?.data?.error || error.message)
      throw error
    } finally {
      commit('SET_USERS_LOADING', false)
    }
  },

  // Reset user password
  async resetUserPassword({ commit }, { userId, password }) {
    commit('SET_USERS_LOADING', true)
    commit('SET_USERS_ERROR', null)
    try {
      const api = getAPI()
      await api.post(`/users/${userId}/reset-password`, { password })
    } catch (error) {
      commit('SET_USERS_ERROR', error.response?.data?.error || error.message)
      throw error
    } finally {
      commit('SET_USERS_LOADING', false)
    }
  },

  // Toggle user status (lock/unlock)
  async toggleUserStatus({ commit }, { userId, status }) {
    commit('SET_USERS_LOADING', true)
    commit('SET_USERS_ERROR', null)
    try {
      const api = getAPI()
      const response = await api.put(`/users/${userId}/status`, { status })
      commit('UPDATE_USER', response.data)
      return response.data
    } catch (error) {
      commit('SET_USERS_ERROR', error.response?.data?.error || error.message)
      throw error
    } finally {
      commit('SET_USERS_LOADING', false)
    }
  },

  // Get total members count (for dashboard)
  async fetchTotalMembers({ commit }) {
    commit('SET_USERS_LOADING', true)
    commit('SET_USERS_ERROR', null)
    try {
      const api = getAPI()
      const response = await api.get('/users/count')
      commit('SET_TOTAL_MEMBERS', response.data.count)
      return response.data
    } catch (error) {
      commit('SET_USERS_ERROR', error.response?.data?.error || error.message)
      throw error
    } finally {
      commit('SET_USERS_LOADING', false)
    }
  }
}

const getters = {
  allUsers: state => state.users,
  totalUsers: state => state.totalUsers,
  currentPage: state => state.currentPage,
  perPage: state => state.perPage,
  totalPages: state => state.totalPages,
  usersLoading: state => state.loading,
  usersError: state => state.error
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
} 