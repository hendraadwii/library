<template>
  <div class="flex flex-col items-center justify-center min-h-[60vh] bg-white rounded-lg shadow p-8">
    <svg class="h-16 w-16 text-indigo-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
    </svg>
    <h2 class="text-2xl font-bold text-gray-800 mb-2">User Management</h2>
    <div v-if="loading" class="text-gray-500 mb-2">Loading...</div>
    <div v-else-if="error" class="text-red-500 mb-2">{{ error }}</div>
    <div v-else-if="totalUsers !== null" class="text-lg font-semibold text-indigo-700 mb-2">
      Total Users: {{ totalUsers }}
    </div>
    <p class="text-gray-500 mb-4 text-center max-w-md">
      Fitur manajemen user akan segera hadir di sini.<br />
      Anda akan dapat melihat, menambah, mengedit, dan menghapus user dari sistem.
    </p>
    <button @click="fetchTotalUsers" class="mt-2 px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition">Refresh</button>
  </div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'Users',
  data() {
    return {
      totalUsers: null,
      loading: false,
      error: null
    }
  },
  methods: {
    async fetchTotalUsers() {
      this.loading = true
      this.error = null
      try {
        const token = localStorage.getItem('token')
        const res = await axios.get('http://localhost:8081/api/v1/users/count', {
          headers: { Authorization: `Bearer ${token}` }
        })
        this.totalUsers = res.data.count
      } catch (err) {
        this.error = err.response?.data?.error || err.message
      } finally {
        this.loading = false
      }
    }
  },
  mounted() {
    this.fetchTotalUsers()
  }
}
</script>

<style scoped>
</style> 