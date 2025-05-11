// API utility with axios instance
import axios from 'axios'

// Base API instance
const api = axios.create({
  baseURL: process.env.VUE_APP_API_URL || 'http://localhost:8081/api/v1',
  headers: {
    'Content-Type': 'application/json',
  }
})

// Get API instance with auth headers
export const getAPI = () => {
  // Get access token from localStorage
  const accessToken = localStorage.getItem('access_token')
  
  // If token exists, set Authorization header
  if (accessToken) {
    api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
  } else {
    // Remove Authorization header if no token
    delete api.defaults.headers.common['Authorization']
  }
  
  return api
}

export default api 