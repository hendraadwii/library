<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Books</h1>
      <button 
        @click="showAddBookModal = true" 
        class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
      >
        Add New Book
      </button>
    </div>

    <!-- Search and Filter Controls -->
    <div class="bg-white rounded-lg shadow mb-6 p-4">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label for="title" class="block text-sm font-medium text-gray-700">Title</label>
          <input 
            type="text"
            id="title"
            v-model="filters.title"
            placeholder="Search by title"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
        </div>
        <div>
          <label for="author" class="block text-sm font-medium text-gray-700">Author</label>
          <input 
            type="text"
            id="author"
            v-model="filters.author"
            placeholder="Search by author"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
        </div>
        <div>
          <label for="category" class="block text-sm font-medium text-gray-700">Category</label>
          <input 
            type="text"
            id="category"
            v-model="filters.category"
            placeholder="Filter by category"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
        </div>
        <div class="flex items-end">
          <button 
            @click="searchBooks" 
            class="bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 w-full"
          >
            Search
          </button>
        </div>
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
      <!-- Books Table -->
      <div class="bg-white rounded-lg shadow overflow-hidden">
        <div v-if="Array.isArray(books) && books.length === 0" class="text-center text-gray-500 py-8">
          Tidak ada buku ditemukan
        </div>
        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Title</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Author</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Category</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stock</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-if="loading" class="hover:bg-gray-50">
                <td colspan="5" class="px-6 py-4 text-center">
                  <div class="flex justify-center items-center">
                    <svg class="animate-spin h-5 w-5 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <span class="ml-2">Loading...</span>
                  </div>
                </td>
              </tr>
              <tr v-else-if="books.length === 0" class="hover:bg-gray-50">
                <td colspan="5" class="px-6 py-4 text-center text-gray-500">
                  No books found
                </td>
              </tr>
              <tr v-for="book in books" :key="book.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-16 w-12">
                      <img 
                        v-if="book.cover" 
                        :src="getCoverUrl(book.cover)" 
                        :alt="book.title"
                        class="h-16 w-12 object-cover rounded"
                      />
                      <div 
                        v-else 
                        class="h-16 w-12 bg-gray-200 rounded flex items-center justify-center"
                      >
                        <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
                        </svg>
                      </div>
                    </div>
                    <div class="ml-4">
                      <div class="text-sm font-medium text-gray-900">{{ book.title }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm text-gray-500">{{ book.author }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm text-gray-500">{{ book.category }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm text-gray-500">{{ book.stock }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span 
                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                    :class="book.stock > 0 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  >
                    {{ book.stock > 0 ? 'Available' : 'Out of Stock' }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <router-link
                    v-if="book.id"
                    :to="`/books/${book.id}`"
                    class="text-indigo-600 hover:text-indigo-900 mr-3"
                  >View</router-link>
                  <span v-else class="text-gray-400 mr-3">No ID</span>
                  <button @click="editBook(book)" class="text-blue-600 hover:text-blue-900 mr-3">Edit</button>
                  <button @click="confirmDeleteBook(book)" class="text-red-600 hover:text-red-900">Delete</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="px-4 py-3 bg-white border-t border-gray-200 sm:px-6">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-700">
                Showing <span class="font-medium">{{ (currentPage - 1) * perPage + 1 }}</span> to 
                <span class="font-medium">{{ Math.min(currentPage * perPage, totalBooks) }}</span> of 
                <span class="font-medium">{{ totalBooks }}</span> results
              </p>
            </div>
            <div>
              <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                <button
                  @click="changePage(currentPage - 1)"
                  :disabled="currentPage === 1"
                  class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
                  :class="{ 'opacity-50 cursor-not-allowed': currentPage === 1 }"
                >
                  <span class="sr-only">Previous</span>
                  &laquo;
                </button>
                <button
                  v-for="page in pageNumbers"
                  :key="page"
                  @click="changePage(page)"
                  class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium"
                  :class="page === currentPage ? 'bg-indigo-50 text-indigo-600 z-10' : 'text-gray-500 hover:bg-gray-50'"
                >
                  {{ page }}
                </button>
                <button
                  @click="changePage(currentPage + 1)"
                  :disabled="currentPage === totalPages"
                  class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50"
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

    <!-- Add/Edit Book Modal -->
    <div v-if="showAddBookModal || showEditBookModal" class="fixed z-20 inset-0 flex items-center justify-center bg-black bg-opacity-40 transition-all duration-200">
      <div class="w-full h-full flex items-center justify-center">
        <div class="w-full max-w-xs sm:max-w-md md:max-w-lg mx-auto bg-white rounded-xl shadow-2xl border border-gray-200 overflow-hidden transform transition-all duration-200 max-h-[90vh] flex flex-col">
          <div class="px-4 py-6 sm:px-6 overflow-y-auto flex-1" style="max-height: 80vh;">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-lg font-bold text-gray-900">
                {{ showEditBookModal ? 'Edit Book' : 'Add New Book' }}
              </h3>
              <button @click="cancelBookModal" class="text-gray-500 hover:text-gray-700">
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
              </button>
            </div>
            <form @submit.prevent="saveBook">
              <div class="grid grid-cols-1 gap-3 sm:gap-4">
                <div>
                  <label for="book-title" class="block text-xs font-medium text-gray-700 mb-1">Title</label>
                  <input 
                    type="text"
                    id="book-title"
                    v-model="bookForm.title"
                    class="block w-full rounded border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 text-xs sm:text-sm px-3 py-2 mb-1"
                    :class="{ 'border-red-500': formErrors.title }"
                    placeholder="Book Title"
                  />
                  <p v-if="formErrors.title" class="text-xs text-red-600">{{ formErrors.title }}</p>
                </div>
                <div>
                  <label for="book-author" class="block text-xs font-medium text-gray-700 mb-1">Author</label>
                  <input 
                    type="text"
                    id="book-author"
                    v-model="bookForm.author"
                    @input="onAuthorInput"
                    class="block w-full rounded border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 text-xs sm:text-sm px-3 py-2 mb-1"
                    :class="{ 'border-red-500': formErrors.author }"
                    placeholder="Author Name"
                  />
                  <p v-if="formErrors.author" class="text-xs text-red-600">{{ formErrors.author }}</p>
                </div>
                <div>
                  <label for="book-isbn" class="block text-xs font-medium text-gray-700 mb-1">ISBN</label>
                  <input 
                    type="text"
                    id="book-isbn"
                    v-model="bookForm.isbn"
                    @input="onIsbnInput"
                    maxlength="10"
                    class="block w-full rounded border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 text-xs sm:text-sm px-3 py-2 mb-1"
                    :class="{ 'border-red-500': formErrors.isbn }"
                    placeholder="10 digit ISBN"
                  />
                  <p v-if="formErrors.isbn" class="text-xs text-red-600">{{ formErrors.isbn }}</p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Kategori</label>
                  <select v-model="bookForm.category" class="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-indigo-500" required>
                    <option value="Fiksi Ilmiah">Fiksi Ilmiah (Science Fiction)</option>
                    <option value="Fantasi">Fantasi (Fantasy)</option>
                    <option value="Misteri">Misteri (Mystery)</option>
                    <option value="Thriller">Thriller</option>
                    <option value="Romantis">Romantis (Romance)</option>
                    <option value="Horor">Horor (Horror)</option>
                    <option value="Sejarah">Sejarah (Historical)</option>
                  </select>
                  <div v-if="formErrors.category" class="text-xs text-red-500 mt-1">{{ formErrors.category }}</div>
                </div>
                <div>
                  <label for="book-stock" class="block text-xs font-medium text-gray-700 mb-1">Stock</label>
                  <input 
                    type="number"
                    id="book-stock"
                    v-model.number="bookForm.stock"
                    min="0"
                    class="block w-full rounded border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 text-xs sm:text-sm px-3 py-2 mb-1"
                    :class="{ 'border-red-500': formErrors.stock }"
                    placeholder="Stock"
                  />
                  <p v-if="formErrors.stock" class="text-xs text-red-600">{{ formErrors.stock }}</p>
                </div>
                <div>
                  <label for="book-cover" class="block text-xs font-medium text-gray-700 mb-1">Cover</label>
                  <input type="file" id="book-cover" @change="onCoverChange" accept="image/*" class="block w-full text-xs" />
                </div>
                <div>
                  <label for="book-description" class="block text-xs font-medium text-gray-700 mb-1">Description</label>
                  <textarea 
                    id="book-description"
                    v-model="bookForm.description"
                    rows="2"
                    class="block w-full rounded border border-gray-300 focus:border-indigo-500 focus:ring-indigo-500 text-xs sm:text-sm px-3 py-2 mb-1"
                    :class="{ 'border-red-500': formErrors.description }"
                    placeholder="Description"
                  ></textarea>
                  <p v-if="formErrors.description" class="text-xs text-red-600">{{ formErrors.description }}</p>
                </div>
              </div>
              <div class="mt-6 flex flex-col gap-2">
                <button 
                  type="submit"
                  class="w-full rounded-full px-4 py-2 bg-indigo-600 text-base sm:text-sm font-semibold text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition"
                >
                  {{ showEditBookModal ? 'Update' : 'Create' }}
                </button>
                <button 
                  type="button"
                  @click="cancelBookModal" 
                  class="w-full rounded-full px-4 py-2 bg-white border border-gray-300 text-base sm:text-sm font-semibold text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed z-10 inset-0 overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 transition-opacity" aria-hidden="true">
          <div class="absolute inset-0 bg-gray-500 opacity-75"></div>
        </div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
                <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-gray-900">
                  Delete Book
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">
                    Are you sure you want to delete "{{ bookToDelete ? bookToDelete.title : '' }}"? This action cannot be undone.
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button 
              @click="deleteBook" 
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Delete
            </button>
            <button 
              @click="showDeleteModal = false" 
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useToast } from 'vue-toastification'

export default {
  name: 'Books',
  setup() {
    const store = useStore()
    const toast = useToast()
    
    // State
    const loading = ref(true)
    const showAddBookModal = ref(false)
    const showEditBookModal = ref(false)
    const showDeleteModal = ref(false)
    const bookToDelete = ref(null)
    const coverBaseUrl = 'http://localhost:8081/static/cover/'
    const getCoverUrl = (cover) => cover ? `${coverBaseUrl}${cover}` : null
    const filters = reactive({
      title: '',
      author: '',
      category: ''
    })
    const bookForm = reactive({
      id: null,
      title: '',
      author: '',
      isbn: '',
      category: '',
      description: '',
      stock: 1
    })
    const formErrors = reactive({
      title: '',
      author: '',
      isbn: '',
      category: '',
      description: '',
      stock: ''
    })
    const coverFile = ref(null)
    
    // Computed
    const books = computed(() => store.getters['books/allBooks'])
    const totalBooks = computed(() => store.getters['books/totalBooks'])
    const currentPage = computed(() => store.getters['books/currentPage'])
    const perPage = computed(() => store.getters['books/perPage'])
    const totalPages = computed(() => store.getters['books/totalPages'])
    
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
    
    // Methods
    const resetForm = () => {
      bookForm.id = null
      bookForm.title = ''
      bookForm.author = ''
      bookForm.isbn = ''
      bookForm.category = ''
      bookForm.description = ''
      bookForm.stock = 1
      
      // Clear errors
      Object.keys(formErrors).forEach(key => {
        formErrors[key] = ''
      })
    }
    
    const validateForm = () => {
      let isValid = true
      
      // Clear errors
      Object.keys(formErrors).forEach(key => {
        formErrors[key] = ''
      })
      
      // Validate required fields
      if (!bookForm.title) {
        formErrors.title = 'Title is required'
        isValid = false
      }
      
      // Author hanya huruf
      if (!bookForm.author) {
        formErrors.author = 'Author is required'
        isValid = false
      } else if (!/^[A-Za-z\s]+$/.test(bookForm.author)) {
        formErrors.author = 'Author must only contain letters'
        isValid = false
      }
      
      // ISBN hanya 10 digit angka
      if (!bookForm.isbn) {
        formErrors.isbn = 'ISBN is required'
        isValid = false
      } else if (!/^\d{10}$/.test(bookForm.isbn)) {
        formErrors.isbn = 'ISBN must be exactly 10 digits'
        isValid = false
      }
      
      if (!bookForm.category) {
        formErrors.category = 'Category is required'
        isValid = false
      }
      
      if (bookForm.stock < 0) {
        formErrors.stock = 'Stock quantity cannot be negative'
        isValid = false
      }
      
      return isValid
    }
    
    const searchBooks = async () => {
      loading.value = true
      try {
        await store.dispatch('books/getBooks', {
          page: 1, // Reset to first page when searching
          title: filters.title,
          author: filters.author,
          category: filters.category
        })
      } catch (error) {
        toast.error('Failed to search books: ' + (error.response?.data?.error || error.message))
      } finally {
        loading.value = false
      }
    }
    
    const changePage = async (page) => {
      if (page < 1 || page > totalPages.value) return
      
      loading.value = true
      try {
        await store.dispatch('books/getBooks', {
          page,
          title: filters.title,
          author: filters.author,
          category: filters.category
        })
      } catch (error) {
        toast.error('Failed to change page: ' + (error.response?.data?.error || error.message))
      } finally {
        loading.value = false
      }
    }
    
    const editBook = (book) => {
      resetForm()
      
      // Populate form with book data
      bookForm.id = book.id
      bookForm.title = book.title
      bookForm.author = book.author
      bookForm.isbn = book.isbn
      bookForm.category = book.category
      bookForm.description = book.description
      bookForm.stock = book.stock
      
      showEditBookModal.value = true
    }
    
    const onCoverChange = (e) => {
      coverFile.value = e.target.files[0]
    }
    
    const onAuthorInput = (e) => {
      // Hanya huruf dan spasi
      e.target.value = e.target.value.replace(/[^A-Za-z\s]/g, '')
      bookForm.author = e.target.value
    }
    
    const onIsbnInput = (e) => {
      // Hanya angka, maksimal 10 digit
      e.target.value = e.target.value.replace(/[^0-9]/g, '').slice(0, 10)
      bookForm.isbn = e.target.value
    }
    
    const saveBook = async () => {
      if (!validateForm()) return
      try {
        const formData = new FormData()
        formData.append('title', bookForm.title)
        formData.append('author', bookForm.author)
        formData.append('isbn', bookForm.isbn)
        formData.append('category', bookForm.category)
        formData.append('description', bookForm.description)
        formData.append('stock', bookForm.stock)
        
        // Append cover file if it exists
        if (coverFile.value) {
          formData.append('cover', coverFile.value)
        }

        if (showEditBookModal.value) {
          await store.dispatch('books/updateBook', {
            bookId: bookForm.id,
            bookData: formData
          })
          toast.success('Book updated successfully')
        } else {
          await store.dispatch('books/createBook', formData)
          toast.success('Book created successfully')
        }
        
        // Reset form and close modal
        cancelBookModal()
        // Refresh the book list
        await store.dispatch('books/getBooks', { 
          page: currentPage.value,
          title: filters.title,
          author: filters.author,
          category: filters.category
        })
      } catch (error) {
        toast.error((error.response?.data?.error || error.message))
      }
    }
    
    const cancelBookModal = () => {
      showAddBookModal.value = false
      showEditBookModal.value = false
      resetForm()
    }
    
    const confirmDeleteBook = (book) => {
      bookToDelete.value = book
      showDeleteModal.value = true
    }
    
    const deleteBook = async () => {
      if (!bookToDelete.value) return
      
      try {
        await store.dispatch('books/deleteBook', bookToDelete.value.id)
        toast.success('Book deleted successfully')
        showDeleteModal.value = false
        bookToDelete.value = null
        
        // Refresh books after deletion
        await store.dispatch('books/getBooks', { page: currentPage.value })
      } catch (error) {
        toast.error('Failed to delete book: ' + (error.response?.data?.error || error.message))
      }
    }
    
    // Load books on component mount
    onMounted(async () => {
      loading.value = true
      try {
        await store.dispatch('books/getBooks')
      } catch (error) {
        toast.error('Failed to load books: ' + (error.response?.data?.error || error.message))
      } finally {
        loading.value = false
      }
    })
    
    return {
      // State
      loading,
      books,
      totalBooks,
      currentPage,
      perPage,
      totalPages,
      filters,
      showAddBookModal,
      showEditBookModal,
      showDeleteModal,
      bookForm,
      formErrors,
      bookToDelete,
      pageNumbers,
      coverFile,
      coverBaseUrl,
      getCoverUrl,
      
      // Methods
      searchBooks,
      changePage,
      editBook,
      saveBook,
      cancelBookModal,
      confirmDeleteBook,
      deleteBook,
      onCoverChange,
      onAuthorInput,
      onIsbnInput
    }
  }
}
</script> 