<template>
  <div class="container mx-auto px-2 sm:px-4 py-4 sm:py-8">
    <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-6 gap-4">
      <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">My Borrowings</h1>
      <div class="flex flex-row flex-wrap gap-2 sm:gap-3">
        <button 
          @click="getActiveBorrowings" 
          class="px-4 py-2 rounded-md font-semibold text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
          :class="{ 'bg-indigo-900': filter === 'active' }"
        >
          Active
        </button>
        <button 
          @click="getOverdueBorrowings" 
          class="px-4 py-2 rounded-md font-semibold text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
          :class="{ 'bg-red-900': filter === 'overdue' }"
        >
          Overdue
        </button>
        <button 
          @click="getReturnedBorrowings" 
          class="px-4 py-2 rounded-md font-semibold text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2"
          :class="{ 'bg-green-900': filter === 'returned' }"
        >
          Returned
        </button>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center h-64">
      <svg class="animate-spin -ml-1 mr-3 h-8 w-8 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <span>Loading...</span>
    </div>

    <div v-else>
      <!-- Empty state -->
      <div v-if="filteredBorrowings.length === 0" class="bg-white rounded-lg shadow p-8 text-center">
        <div class="flex flex-col items-center">
          <svg class="h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"></path>
          </svg>
          <h3 class="mt-2 text-lg font-medium text-gray-900">No borrowings found</h3>
          <p class="mt-1 text-sm text-gray-500">
            {{ 
              filter === 'active' ? "You don't have any active borrowings." :
              filter === 'overdue' ? "You don't have any overdue books." :
              filter === 'returned' ? "You haven't returned any books yet." :
              "You haven't borrowed any books yet."
            }}
          </p>
          <div class="mt-6">
            <router-link to="/books" class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
              Browse Books
            </router-link>
          </div>
        </div>
      </div>

      <!-- Borrowings Table -->
      <div v-else class="bg-white rounded-lg shadow overflow-hidden">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200 text-xs sm:text-sm">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 uppercase tracking-wider">Book</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 uppercase tracking-wider">Borrow Date</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 uppercase tracking-wider">Due Date</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 uppercase tracking-wider">Return Date</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="borrowing in filteredBorrowings" :key="borrowing.id">
                <!-- Book Title and Author -->
                <td class="px-4 py-3 whitespace-nowrap">
                  <div class="flex items-center">
                    <div>
                      <div class="text-sm font-medium text-gray-900">{{ borrowing.book_title }}</div>
                      <div class="text-xs text-gray-500">{{ borrowing.book_author }}</div>
                    </div>
                  </div>
                </td>
                <!-- Borrow Date -->
                <td class="px-4 py-3 whitespace-nowrap">
                  <div class="text-sm text-gray-900">{{ formatDate(borrowing.borrow_date) }}</div>
                </td>
                <!-- Due Date -->
                <td class="px-4 py-3 whitespace-nowrap">
                  <div class="text-sm" :class="isOverdue(borrowing) ? 'text-red-600 font-medium' : 'text-gray-900'">
                    {{ formatDate(borrowing.due_date) }}
                    <span v-if="isOverdue(borrowing)" class="text-xs text-red-600">
                      ({{ getDaysOverdue(borrowing) }} days overdue)
                    </span>
                  </div>
                </td>
                <!-- Status -->
                <td class="px-4 py-3 whitespace-nowrap">
                  <span 
                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                    :class="getStatusClass(borrowing)"
                  >
                    {{ borrowing.status === 'borrowed' ? (isOverdue(borrowing) ? 'Overdue' : 'Borrowed') : 'Returned' }}
                  </span>
                </td>
                <!-- Return Date -->
                <td class="px-4 py-3 whitespace-nowrap">
                  <div class="text-sm text-gray-900">
                    {{ borrowing.return_date ? formatDate(borrowing.return_date) : '-' }}
                  </div>
                </td>
                <!-- Actions -->
                <td class="px-4 py-3 whitespace-nowrap text-sm font-medium">
                  <router-link
                    v-if="borrowing.book_id"
                    :to="`/books/${borrowing.book_id}`"
                    class="text-indigo-600 hover:text-indigo-900 mr-3"
                  >
                    View Book
                  </router-link>
                  <span v-else class="text-gray-400 mr-3">No Book</span>
                  <button 
                    v-if="borrowing.status === 'borrowed' && !borrowing.return_date && borrowing.id" 
                    @click="returnBook(borrowing.id)" 
                    class="text-green-600 hover:text-green-900"
                  >
                    Return
                  </button>
                  <span v-else-if="borrowing.status === 'borrowed' && !borrowing.return_date && !borrowing.id" class="text-gray-400">No ID</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <!-- Pagination -->
        <div v-if="totalPages > 1" class="px-2 py-3 bg-white border-t border-gray-200 sm:px-6">
          <div class="flex flex-col sm:flex-row items-center justify-between gap-2">
            <div>
              <p class="text-xs sm:text-sm text-gray-700">
                Showing <span class="font-medium">{{ (currentPage - 1) * perPage + 1 }}</span> to 
                <span class="font-medium">{{ Math.min(currentPage * perPage, totalBorrowings) }}</span> of 
                <span class="font-medium">{{ totalBorrowings }}</span> results
              </p>
            </div>
            <div>
              <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                <button
                  @click="changePage(currentPage - 1)"
                  :disabled="currentPage === 1"
                  class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50"
                  :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
                >
                  <span class="sr-only">Previous</span>
                  &laquo;
                </button>
                <button
                  v-for="page in pageNumbers"
                  :key="page"
                  @click="changePage(page)"
                  class="relative inline-flex items-center px-3 sm:px-4 py-2 border border-gray-300 bg-white text-xs sm:text-sm font-medium"
                  :class="page === currentPage ? 'bg-indigo-50 text-indigo-600 z-10' : 'text-gray-500 hover:bg-gray-50'"
                >
                  {{ page }}
                </button>
                <button
                  @click="changePage(currentPage + 1)"
                  :disabled="currentPage === totalPages"
                  class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50"
                  :class="{ 'opacity-50 cursor-not-allowed': currentPage === totalPages }"
                >
                  <span class="sr-only">Next</span>
                  &raquo;
                </button>
              </nav>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useToast } from 'vue-toastification'

export default {
  name: 'Borrowings',
  setup() {
    const store = useStore()
    const toast = useToast()
    
    // State
    const loading = ref(true)
    const filter = ref('active') // active, overdue, returned, all
    const showAddBorrowingModal = ref(false)
    
    // Computed
    const borrowings = computed(() => store.getters['borrowings/allBorrowings'])
    const totalBorrowings = computed(() => store.getters['borrowings/totalBorrowings'])
    const currentPage = computed(() => store.getters['borrowings/currentPage'])
    const perPage = computed(() => store.getters['borrowings/perPage'])
    const totalPages = computed(() => store.getters['borrowings/totalPages'])
    
    // Computed for pagination display
    const pageNumbers = computed(() => {
      const pages = []
      const maxVisible = 5
      let startPage = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
      const endPage = Math.min(totalPages.value, startPage + maxVisible - 1)
      
      // Adjust start page if we're near the end
      if (endPage - startPage + 1 < maxVisible) {
        startPage = Math.max(1, endPage - maxVisible + 1)
      }
      
      for (let i = startPage; i <= endPage; i++) {
        pages.push(i)
      }
      
      return pages
    })
    
    // Filtered borrowings based on current filter
    const filteredBorrowings = computed(() => {
      if (filter.value === 'active') {
        return borrowings.value.filter(b => b.status === 'borrowed' && !isOverdue(b))
      }
      if (filter.value === 'overdue') {
        return borrowings.value.filter(b => b.status === 'borrowed' && isOverdue(b))
      }
      if (filter.value === 'returned') {
        return borrowings.value.filter(b => b.status === 'returned' || b.return_date)
      }
      return borrowings.value
    })
    
    // Format date
    const formatDate = (dateString) => {
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })
    }
    
    // Check if a borrowing is overdue
    const isOverdue = (borrowing) => {
      if (borrowing.return_date) return false
      
      const today = new Date()
      const dueDate = new Date(borrowing.due_date)
      return today > dueDate
    }
    
    // Get days overdue
    const getDaysOverdue = (borrowing) => {
      if (borrowing.return_date || !isOverdue(borrowing)) return 0
      
      const today = new Date()
      const dueDate = new Date(borrowing.due_date)
      const diffTime = Math.abs(today - dueDate)
      return Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    }
    
    // Get status class for styling
    const getStatusClass = (borrowing) => {
      if (borrowing.return_date || borrowing.status === 'returned') {
        return 'bg-green-100 text-green-800'
      }
      
      if (isOverdue(borrowing)) {
        return 'bg-red-100 text-red-800'
      }
      
      return 'bg-yellow-100 text-yellow-800'
    }
    
    // Load borrowings
    const loadBorrowings = async (options = {}) => {
      loading.value = true
      
      try {
        await store.dispatch('borrowings/getBorrowings', options)
      } catch (error) {
        toast.error('Failed to load borrowings: ' + (error.response?.data?.error || error.message))
      } finally {
        loading.value = false
      }
    }
    
    // Filter functions
    const getActiveBorrowings = () => {
      filter.value = 'active'
      loadBorrowings({ status: 'borrowed', overdue: false })
    }
    
    const getOverdueBorrowings = () => {
      filter.value = 'overdue'
      loadBorrowings({ overdue: true })
    }
    
    const getReturnedBorrowings = () => {
      filter.value = 'returned'
      loadBorrowings({ status: 'returned' })
    }
    
    const getAllBorrowings = () => {
      filter.value = 'all'
      loadBorrowings()
    }
    
    // Change page
    const changePage = async (page) => {
      if (page < 1 || page > totalPages.value) return
      
      loading.value = true
      try {
        let options = { page }
        
        if (filter.value === 'active') {
          options.status = 'borrowed'
          options.overdue = false
        } else if (filter.value === 'overdue') {
          options.overdue = true
        } else if (filter.value === 'returned') {
          options.status = 'returned'
        }
        
        await store.dispatch('borrowings/getBorrowings', options)
      } catch (error) {
        toast.error('Failed to change page: ' + (error.response?.data?.error || error.message))
      } finally {
        loading.value = false
      }
    }
    
    // Return a book
    const returnBook = async (borrowingId) => {
      if (!borrowingId) {
        toast.error('Invalid borrowing ID');
        return;
      }
      try {
        await store.dispatch('borrowings/returnBook', borrowingId)
        toast.success('Book returned successfully')
        // Reload borrowings to reflect the updated status
        await loadBorrowings({ 
          status: filter.value === 'active' ? 'borrowed' : '',
          overdue: filter.value === 'overdue'
        })
      } catch (error) {
        toast.error('Failed to return book: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Load active borrowings on component mount
    onMounted(() => {
      getActiveBorrowings()
    })
    
    return {
      loading,
      filter,
      borrowings,
      filteredBorrowings,
      totalBorrowings,
      currentPage,
      perPage,
      totalPages,
      pageNumbers,
      formatDate,
      isOverdue,
      getDaysOverdue,
      getStatusClass,
      getActiveBorrowings,
      getOverdueBorrowings,
      getReturnedBorrowings,
      getAllBorrowings,
      changePage,
      returnBook,
      showAddBorrowingModal
    }
  }
}
</script> 