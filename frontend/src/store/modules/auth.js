import axios from 'axios'
import jwt_decode from 'jwt-decode'

// Utility functions
const getTokens = () => {
  const accessToken = localStorage.getItem('access_token')
  const refreshToken = localStorage.getItem('refresh_token')
  return { accessToken, refreshToken }
}

const setTokens = (tokens) => {
  localStorage.setItem('access_token', tokens.access_token)
  localStorage.setItem('refresh_token', tokens.refresh_token)
}

const clearTokens = () => {
  localStorage.removeItem('access_token')
  localStorage.removeItem('refresh_token')
}

// Create API instance
const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL || 'http://localhost:8081/api/v1',
  headers: {
    'Content-Type': 'application/json',
  }
})

// Configure request interceptor for API calls
api.interceptors.request.use(
  config => {
    const { accessToken } = getTokens()
    if (accessToken) {
      config.headers.Authorization = `Bearer ${accessToken}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Configure response interceptor for API calls
api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config
    
    // If error is unauthorized and not a retry
    if (error.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      try {
        // Get tokens
        const { refreshToken } = getTokens()
        if (!refreshToken) {
          return Promise.reject(error)
        }
        
        // Try to refresh tokens
        const response = await axios.post(`${api.defaults.baseURL}/auth/refresh`, {
          refresh_token: refreshToken
        })
        
        // If successful, update tokens
        if (response.data.access_token && response.data.refresh_token) {
          setTokens(response.data)
          
          // Update authorization header
          originalRequest.headers.Authorization = `Bearer ${response.data.access_token}`
          return api(originalRequest)
        }
      } catch (refreshError) {
        // Handle refresh error (usually by redirecting to login)
        clearTokens()
        return Promise.reject(refreshError)
      }
    }
    
    return Promise.reject(error)
  }
)

export default {
  namespaced: true,
  state: {
    user: null,
    accessToken: null,
    refreshToken: null,
  },
  mutations: {
    SET_USER(state, user) {
      state.user = user
    },
    SET_TOKENS(state, { accessToken, refreshToken }) {
      state.accessToken = accessToken
      state.refreshToken = refreshToken
    },
    CLEAR_AUTH(state) {
      state.user = null
      state.accessToken = null
      state.refreshToken = null
    }
  },
  actions: {
    // Try to restore authentication from localStorage
    initAuth({ commit, dispatch }) {
      const { accessToken, refreshToken } = getTokens()
      
      if (accessToken && refreshToken) {
        try {
          // Decode JWT to get user info
          const decodedToken = jwt_decode(accessToken)
          
          // Check if token is expired
          const currentTime = Date.now() / 1000
          if (decodedToken.exp < currentTime) {
            // Token expired, try to refresh
            dispatch('refreshTokens')
            return
          }
          
          // Set user and tokens
          commit('SET_USER', {
            id: decodedToken.user_id,
            email: decodedToken.email,
            role: decodedToken.role
          })
          commit('SET_TOKENS', { accessToken, refreshToken })
        } catch (error) {
          console.error('Auth initialization error:', error)
          dispatch('logout')
        }
      }
    },
    
    // Login user
    async login({ commit }, credentials) {
      const response = await api.post('/auth/login', credentials)
      
      // Save tokens to localStorage
      setTokens(response.data)
      
      // Decode JWT to get user info
      const decodedToken = jwt_decode(response.data.access_token)
      
      // Set user and tokens in state
      commit('SET_USER', {
        id: decodedToken.user_id,
        email: decodedToken.email,
        role: decodedToken.role
      })
      commit('SET_TOKENS', {
        accessToken: response.data.access_token,
        refreshToken: response.data.refresh_token
      })
      
      return response.data
    },
    
    // Register user
    async register({ commit }, userData) {
      const response = await api.post('/auth/register', userData)
      
      // Save tokens to localStorage
      setTokens(response.data)
      
      // Decode JWT to get user info
      const decodedToken = jwt_decode(response.data.access_token)
      
      // Set user and tokens in state
      commit('SET_USER', {
        id: decodedToken.user_id,
        email: decodedToken.email,
        role: decodedToken.role
      })
      commit('SET_TOKENS', {
        accessToken: response.data.access_token,
        refreshToken: response.data.refresh_token
      })
      
      return response.data
    },
    
    // Forgot password
    async forgotPassword(_, { email }) {
      const response = await api.post('/auth/forgot-password', { email })
      return response.data
    },
    
    // Verify PIN
    async verifyPin(_, { email, pin }) {
      const response = await api.post('/auth/verify-pin', { email, pin })
      return response.data
    },
    
    // Reset password
    async resetPassword(_, { email, pin, newPassword }) {
      const response = await api.post('/auth/reset-password', { email, pin, new_password: newPassword })
      return response.data
    },
    
    // Refresh tokens
    async refreshTokens({ commit }) {
      try {
        const { refreshToken } = getTokens()
        if (!refreshToken) {
          throw new Error('No refresh token available')
        }
        
        const response = await axios.post(`${api.defaults.baseURL}/auth/refresh`, {
          refresh_token: refreshToken
        })
        
        // Save tokens to localStorage
        setTokens(response.data)
        
        // Decode JWT to get user info
        const decodedToken = jwt_decode(response.data.access_token)
        
        // Set user and tokens in state
        commit('SET_USER', {
          id: decodedToken.user_id,
          email: decodedToken.email,
          role: decodedToken.role
        })
        commit('SET_TOKENS', {
          accessToken: response.data.access_token,
          refreshToken: response.data.refresh_token
        })
        
        return response.data
      } catch (error) {
        console.error('Token refresh error:', error)
        commit('CLEAR_AUTH')
        clearTokens()
        throw error
      }
    },
    
    // Logout user
    logout({ commit }) {
      // Clear auth state
      commit('CLEAR_AUTH')
      
      // Clear tokens from localStorage
      clearTokens()
    }
  },
  getters: {
    isAuthenticated: state => !!state.user,
    currentUser: state => state.user,
    userRole: state => state.user ? state.user.role : null,
    getAPI: () => api
  }
} 