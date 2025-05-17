<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Navbar Admin Area Responsive Dropdown -->
    <nav class="bg-white shadow-sm w-full py-2 px-2 sm:py-4 sm:px-4">
      <div class="max-w-7xl mx-auto flex justify-between items-center h-12 sm:h-16">
        <div class="flex items-center space-x-2 sm:space-x-8">
          <span class="text-lg sm:text-xl font-bold text-indigo-600">Admin Area</span>
        </div>
        <!-- Desktop Menu -->
        <div class="hidden sm:flex items-center space-x-4">
          <router-link to="/" class="text-gray-700 hover:text-indigo-600 text-sm font-medium">Dashboard</router-link>
          <router-link to="/books" class="text-gray-700 hover:text-indigo-600 text-sm font-medium">Books</router-link>
          <router-link to="/borrowings" class="text-gray-700 hover:text-indigo-600 text-sm font-medium">Borrowings</router-link>
          <router-link to="/users" class="text-gray-700 hover:text-indigo-600 text-sm font-medium">Users</router-link>
          <button @click="logout" class="ml-4 px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 text-sm font-medium">Logout</button>
        </div>
        <!-- Mobile Hamburger -->
        <div class="sm:hidden flex items-center">
          <button @click="showDropdown = !showDropdown" class="inline-flex items-center justify-center p-2 rounded-md text-gray-700 hover:text-indigo-600 focus:outline-none">
            <svg v-if="!showDropdown" class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
            <svg v-else class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
      <!-- Mobile Dropdown -->
      <transition name="fade">
        <div v-if="showDropdown" class="sm:hidden bg-white shadow rounded-b-lg absolute left-0 right-0 z-50">
          <div class="flex flex-col py-2 px-4 space-y-2">
            <router-link @click="showDropdown = false" to="/" class="text-gray-700 hover:text-indigo-600 text-base font-medium">Dashboard</router-link>
            <router-link @click="showDropdown = false" to="/books" class="text-gray-700 hover:text-indigo-600 text-base font-medium">Books</router-link>
            <router-link @click="showDropdown = false" to="/borrowings" class="text-gray-700 hover:text-indigo-600 text-base font-medium">Borrowings</router-link>
            <router-link @click="showDropdown = false" to="/users" class="text-gray-700 hover:text-indigo-600 text-base font-medium">Users</router-link>
            <button @click="logout" class="w-full text-left px-2 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 text-base font-medium">Logout</button>
          </div>
        </div>
      </transition>
    </nav>
    <div class="container mx-auto max-w-full px-2 sm:px-4 py-4 sm:py-6">
      <div class="mb-6 sm:mb-8">
        <h1 class="text-xl sm:text-3xl font-bold text-gray-900">Dashboard</h1>
        <p class="text-gray-600 text-xs sm:text-base">Welcome to the Library Management System</p>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="flex justify-center items-center h-64">
        <svg class="animate-spin -ml-1 mr-3 h-8 w-8 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span>Memuat data...</span>
      </div>

      <div v-else>
        <div v-if="errorMsg" class="text-center text-red-500 py-4">{{ errorMsg }}</div>
        <!-- Statistics Cards -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div class="bg-gradient-to-r from-indigo-100 to-indigo-200 rounded-xl shadow p-5 flex items-center transition hover:shadow-lg hover:-translate-y-1 duration-200">
            <div class="bg-indigo-500 text-white rounded-full p-4 mr-4">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a4 4 0 00-3-3.87M9 20H4v-2a4 4 0 013-3.87m9-4V7a4 4 0 00-8 0v2m8 4a4 4 0 01-8 0" /></svg>
            </div>
            <div>
              <p class="text-gray-600 text-xs">Total Members</p>
              <p class="text-2xl font-bold">{{ totalMembers }}</p>
            </div>
          </div>
          <div class="bg-gradient-to-r from-purple-100 to-purple-200 rounded-xl shadow p-5 flex items-center transition hover:shadow-lg hover:-translate-y-1 duration-200">
            <div class="bg-purple-500 text-white rounded-full p-4 mr-4">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"/></svg>
            </div>
            <div>
              <p class="text-gray-600 text-xs">Total Books</p>
              <p class="text-2xl font-bold">{{ totalBooks }}</p>
            </div>
          </div>
          <div class="bg-gradient-to-r from-green-100 to-green-200 rounded-xl shadow p-5 flex items-center transition hover:shadow-lg hover:-translate-y-1 duration-200">
            <div class="bg-green-500 text-white rounded-full p-4 mr-4">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            </div>
            <div>
              <p class="text-gray-600 text-xs">Active Borrowings</p>
              <p class="text-2xl font-bold">{{ activeBorrowings }}</p>
            </div>
          </div>
          <div class="bg-gradient-to-r from-red-100 to-red-200 rounded-xl shadow p-5 flex items-center transition hover:shadow-lg hover:-translate-y-1 duration-200">
            <div class="bg-red-500 text-white rounded-full p-4 mr-4">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
            </div>
            <div>
              <p class="text-gray-600 text-xs">Overdue Books</p>
              <p class="text-2xl font-bold">{{ overdueBorrowings.length }}</p>
            </div>
          </div>
        </div>
        <!-- Most Borrowed Books -->
        <div class="bg-white rounded-xl shadow mb-8 overflow-x-auto">
          <div class="border-b pb-2 mb-4 flex items-center">
            <h2 class="text-lg font-semibold text-gray-800">Most Borrowed Books</h2>
          </div>
          <div class="p-2 sm:p-4">
            <div v-if="Array.isArray(mostBorrowedBooks) && mostBorrowedBooks.length === 0" class="text-center text-gray-500 py-4">
              Belum ada data buku yang sering dipinjam.
            </div>
            <div v-else>
              <table class="min-w-full divide-y divide-gray-200 rounded-xl overflow-hidden">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Title</th>
                    <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Author</th>
                    <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Borrow Count</th>
                    <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="book in mostBorrowedBooks" :key="book.id" class="odd:bg-gray-50">
                    <td class="px-6 py-4 font-medium">{{ book.title }}</td>
                    <td class="px-6 py-4">{{ book.author }}</td>
                    <td class="px-6 py-4">{{ book.borrow_count }}</td>
                    <td class="px-6 py-4 font-medium">
                      <router-link :to="`/books/${book.id}`" class="text-indigo-600 hover:text-indigo-900 transition">View</router-link>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        <!-- Overdue Borrowings -->
        <div class="bg-white rounded-xl shadow overflow-x-auto">
          <div class="border-b pb-2 mb-4 flex items-center">
            <h2 class="text-lg font-semibold text-gray-800">Overdue Books</h2>
          </div>
          <div class="p-2 sm:p-4">
            <div v-if="Array.isArray(overdueBorrowings) && overdueBorrowings.length === 0" class="text-center text-gray-500 py-4">
              Tidak ada buku yang terlambat dikembalikan.
            </div>
            <div v-else>
              <table class="min-w-full divide-y divide-gray-200 rounded-xl overflow-hidden">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Book</th>
                    <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Due Date</th>
                    <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Days Overdue</th>
                    <th class="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="borrowing in overdueBorrowings" :key="borrowing.id" class="odd:bg-gray-50">
                    <td class="px-6 py-4 font-medium">{{ borrowing.book_title }}</td>
                    <td class="px-6 py-4">{{ formatDate(borrowing.due_date) }}</td>
                    <td class="px-6 py-4 text-red-500 font-semibold">{{ getDaysOverdue(borrowing.due_date) }} days</td>
                    <td class="px-6 py-4 font-medium">
                      <button @click="returnBook(borrowing.id)" class="text-indigo-600 hover:text-indigo-900 transition">Return</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { setupSessionTimeout } from '../utils/sessionTimeout'

export default {
  name: 'Home',
  setup() {
    const store = useStore()
    const router = useRouter()
    const toast = useToast()
    const loading = ref(true)
    const mostBorrowedBooks = ref([])
    const overdueBorrowings = ref([])
    const showDropdown = ref(false)
    const errorMsg = ref('')
    
    // Setup session timeout
    const { resetTimer, setupActivityListeners, cleanup } = setupSessionTimeout(3)
    
    // Total Members (role: member)
    const totalMembers = computed(() => store.getters['users/totalMembers'])
    const usersLoading = computed(() => store.getters['users/usersLoading'])
    const usersError = computed(() => store.getters['users/usersError'])
    const userRole = computed(() => store.getters['auth/userRole'])
    
    // Computed properties
    const totalBooks = computed(() => store.getters['books/totalBooks'])
    const activeBorrowings = computed(() => 
      store.getters['borrowings/borrowingsByStatus']('borrowed') 
        ? store.getters['borrowings/borrowingsByStatus']('borrowed').length 
        : 0
    )
    
    // Format date
    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })
    }
    
    // Get days overdue
    const getDaysOverdue = (dueDate) => {
      const today = new Date()
      const due = new Date(dueDate)
      const diffTime = Math.abs(today - due)
      return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    }
    
    // Return a book
    const returnBook = async (borrowingId) => {
      try {
        await store.dispatch('borrowings/returnBook', borrowingId)
        toast.success('Book returned successfully')
        await fetchOverdueBorrowings() // Refresh the list
      } catch (error) {
        toast.error('Failed to return book: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Fetch most borrowed books
    const fetchMostBorrowedBooks = async () => {
      try {
        mostBorrowedBooks.value = await store.dispatch('books/getMostBorrowedBooks', 5)
      } catch (error) {
        toast.error('Failed to load most borrowed books: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Fetch overdue borrowings
    const fetchOverdueBorrowings = async () => {
      try {
        const response = await store.dispatch('borrowings/getOverdueBorrowings')
        overdueBorrowings.value = response.borrowings || []
      } catch (error) {
        toast.error('Failed to load overdue borrowings: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Fetch books with pagination
    const fetchBooks = async () => {
      try {
        await store.dispatch('books/getBooks')
      } catch (error) {
        toast.error('Failed to load books: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Fetch borrowings
    const fetchBorrowings = async () => {
      try {
        await store.dispatch('borrowings/getBorrowings')
      } catch (error) {
        toast.error('Failed to load borrowings: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Load data on component mount
    onMounted(async () => {
      loading.value = true
      errorMsg.value = ''
      try {
        await Promise.all([
          fetchBooks(),
          fetchBorrowings(),
          fetchMostBorrowedBooks(),
          fetchOverdueBorrowings(),
          store.dispatch('users/fetchTotalMembers')
        ])
        // Setup session timeout after data is loaded
        setupActivityListeners()
        resetTimer()
      } catch (error) {
        console.error('Error loading dashboard data:', error)
        errorMsg.value = error.response?.data?.error || error.message || 'Gagal memuat data dashboard.'
        toast.error('Gagal memuat data dashboard: ' + errorMsg.value)
      } finally {
        loading.value = false
      }
    })

    // Cleanup on component unmount
    onUnmounted(() => {
      cleanup()
    })
    
    const logout = () => {
      showDropdown.value = false
      store.dispatch('auth/logout')
      router.push('/login')
      toast.success('Logged out successfully')
    }
    
    return {
      loading,
      totalBooks,
      activeBorrowings,
      mostBorrowedBooks,
      overdueBorrowings,
      formatDate,
      getDaysOverdue,
      returnBook,
      showDropdown,
      logout,
      totalMembers,
      usersLoading,
      usersError,
      userRole,
      errorMsg
    }
  }
}
</script>

<style>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style> 