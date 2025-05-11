<template>
  <div class="container mx-auto px-4 py-8">
    <!-- Back Button -->
    <div class="mb-6">
      <router-link to="/books" class="flex items-center text-indigo-600 hover:text-indigo-900">
        <svg class="h-5 w-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
        </svg>
        Back to Books
      </router-link>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center h-64">
      <svg class="animate-spin -ml-1 mr-3 h-8 w-8 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
      <span>Loading...</span>
    </div>

    <div v-else-if="!book" class="bg-white rounded-lg shadow p-6 text-center">
      <p class="text-gray-500">Book not found or has been deleted.</p>
      <router-link to="/books" class="mt-4 inline-block text-indigo-600 hover:text-indigo-900">
        View all books
      </router-link>
    </div>

    <div v-else class="bg-white rounded-lg shadow overflow-hidden">
      <!-- Book Details Header -->
      <div class="bg-gray-50 px-6 py-4 border-b border-gray-200 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-900">{{ book.title }}</h1>
        <div>
          <button 
            v-if="userCanBorrow" 
            @click="showBorrowModal = true" 
            class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="book.stock_quantity <= 0"
          >
            Borrow Book
          </button>
        </div>
      </div>

      <!-- Book Details Content -->
      <div class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Left: Cover Image -->
          <div class="md:col-span-1">
            <div class="aspect-w-2 aspect-h-3 rounded-lg overflow-hidden bg-gray-100">
              <img 
                v-if="book.cover" 
                :src="getCoverUrl(book.cover)" 
                :alt="book.title"
                class="w-full h-full object-cover"
              />
              <div 
                v-else 
                class="w-full h-full flex items-center justify-center"
              >
                <svg class="h-24 w-24 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                </svg>
              </div>
            </div>
          </div>
          
          <!-- Right: Basic Info -->
          <div class="md:col-span-2 space-y-4">
            <div>
              <h2 class="text-xl font-semibold text-gray-800">About</h2>
              <p class="mt-2 text-gray-600">{{ book.description || 'No description available.' }}</p>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <h3 class="text-sm font-medium text-gray-500">Author</h3>
                <p class="mt-1 text-gray-900">{{ book.author }}</p>
              </div>
              
              <div>
                <h3 class="text-sm font-medium text-gray-500">ISBN</h3>
                <p class="mt-1 text-gray-900">{{ book.isbn }}</p>
              </div>
              
              <div>
                <h3 class="text-sm font-medium text-gray-500">Category</h3>
                <p class="mt-1 text-gray-900">{{ book.category }}</p>
              </div>
              
              <div>
                <h3 class="text-sm font-medium text-gray-500">Stock</h3>
                <p class="mt-1">
                  <span 
                    class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full"
                    :class="book.stock_quantity > 0 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  >
                    {{ book.stock_quantity > 0 ? `${book.stock_quantity} Available` : 'Out of Stock' }}
                  </span>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Borrow Book Modal -->
    <div v-if="showBorrowModal" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity" aria-hidden="true">
          <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
        </div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-indigo-100 sm:mx-0 sm:h-10 sm:w-10">
                <svg class="h-6 w-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-gray-900">
                  Borrow Book
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">
                    You are about to borrow "{{ book ? book.title : '' }}". Please select a return date.
                  </p>
                  
                  <div class="mt-4">
                    <label for="due-date" class="block text-sm font-medium text-gray-700">Due Date</label>
                    <input 
                      type="date" 
                      id="due-date" 
                      v-model="borrowForm.dueDate" 
                      class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                      :min="minDueDate"
                    />
                    <p v-if="formError" class="mt-1 text-sm text-red-600">{{ formError }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button 
              @click="borrowBook" 
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-indigo-600 text-base font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Borrow
            </button>
            <button 
              @click="showBorrowModal = false" 
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:w-auto sm:text-sm"
            >
              Cancel
            </button>
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
  name: 'BookDetails',
  props: {
    id: {
      type: [String, Number],
      required: true
    }
  },
  setup(props) {
    const store = useStore()
    const toast = useToast()
    
    // State
    const loading = ref(true)
    const loadingBorrowings = ref(true)
    const showBorrowModal = ref(false)
    const formError = ref('')
    const borrowForm = ref({
      dueDate: ''
    })
    const coverBaseUrl = 'http://localhost:8081/static/cover/'
    const getCoverUrl = (cover) => cover ? `${coverBaseUrl}${cover}` : null
    
    // Set default due date to 2 weeks from today
    const today = new Date()
    const twoWeeksLater = new Date(today)
    twoWeeksLater.setDate(today.getDate() + 14)
    borrowForm.value.dueDate = twoWeeksLater.toISOString().substr(0, 10)
    
    // Min due date (tomorrow)
    const tomorrow = new Date()
    tomorrow.setDate(today.getDate() + 1)
    const minDueDate = tomorrow.toISOString().substr(0, 10)
    
    // Computed
    const book = computed(() => store.getters['books/currentBook'])
    const currentUser = computed(() => store.getters['auth/currentUser'])
    
    // Determine if the current user can borrow (not already borrowed)
    const userCanBorrow = computed(() => {
      if (!currentUser.value || !userBorrowings.value) return false
      
      return !userBorrowings.value.some(
        b => b.status === 'borrowed' && !b.return_date
      )
    })
    
    // Get user's active borrowing of this book if any
    const userActiveBorrowing = computed(() => {
      if (!userBorrowings.value) return null
      
      return userBorrowings.value.find(
        b => b.status === 'borrowed' && !b.return_date
      )
    })
    
    // Get user's borrowing history for this book
    const userBorrowingHistory = computed(() => {
      if (!userBorrowings.value) return []
      
      return userBorrowings.value.filter(
        b => b.return_date || b.status === 'returned'
      )
    })
    
    // All borrowings for this book by the current user
    const userBorrowings = computed(() => {
      const allBorrowings = store.getters['borrowings/allBorrowings']
      const bookId = parseInt(props.id)
      
      if (!allBorrowings || !currentUser.value) return []
      
      return allBorrowings.filter(
        b => b.book_id === bookId && b.user_id === currentUser.value.id
      )
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
    
    // Load book details and borrowings
    const loadBookDetails = async () => {
      loading.value = true
      
      try {
        await store.dispatch('books/getBook', props.id)
        if (!book.value) {
          toast.error('Book not found')
        }
      } catch (error) {
        toast.error('Failed to load book details: ' + (error.response?.data?.error || error.message))
      } finally {
        loading.value = false
      }
    }
    
    // Load book borrowings
    const loadBorrowings = async () => {
      loadingBorrowings.value = true
      
      try {
        await store.dispatch('borrowings/getBorrowings', {
          bookId: props.id
        })
      } catch (error) {
        toast.error('Failed to load borrowing history: ' + (error.response?.data?.error || error.message))
      } finally {
        loadingBorrowings.value = false
      }
    }
    
    // Borrow book
    const borrowBook = async () => {
      if (!book.value) return
      
      formError.value = ''
      
      // Validate due date
      const dueDate = new Date(borrowForm.value.dueDate)
      const today = new Date()
      
      if (dueDate <= today) {
        formError.value = 'Due date must be in the future'
        return
      }
      
      try {
        await store.dispatch('borrowings/borrowBook', {
          book_id: book.value.id,
          due_date: dueDate.toISOString()
        })
        
        toast.success('Book borrowed successfully')
        showBorrowModal.value = false
        
        // Reload book and borrowings
        await Promise.all([
          loadBookDetails(),
          loadBorrowings()
        ])
      } catch (error) {
        toast.error('Failed to borrow book: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Return book
    const returnBook = async (borrowingId) => {
      try {
        await store.dispatch('borrowings/returnBook', borrowingId)
        toast.success('Book returned successfully')
        
        // Reload book and borrowings
        await Promise.all([
          loadBookDetails(),
          loadBorrowings()
        ])
      } catch (error) {
        toast.error('Failed to return book: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Load data on component mount
    onMounted(async () => {
      await Promise.all([
        loadBookDetails(),
        loadBorrowings()
      ])
    })
    
    return {
      loading,
      loadingBorrowings,
      book,
      showBorrowModal,
      borrowForm,
      formError,
      minDueDate,
      userCanBorrow,
      userActiveBorrowing,
      userBorrowingHistory,
      formatDate,
      borrowBook,
      returnBook,
      coverBaseUrl,
      getCoverUrl
    }
  }
}
</script> 