<template>
  <div class="relative min-h-screen bg-gray-100 p-4 sm:p-8 overflow-hidden">
    <!-- Background effect -->
    <div class="pointer-events-none select-none absolute inset-0 -z-10">
      <div class="absolute top-[-100px] left-[-100px] w-[400px] h-[400px] bg-gradient-to-br from-indigo-400 via-blue-300 to-purple-300 opacity-40 rounded-full blur-3xl animate-blob1"></div>
      <div class="absolute bottom-[-120px] right-[-80px] w-[350px] h-[350px] bg-gradient-to-tr from-pink-300 via-indigo-200 to-blue-200 opacity-30 rounded-full blur-2xl animate-blob2"></div>
      <div class="absolute top-1/2 left-1/2 w-[250px] h-[250px] bg-gradient-to-br from-green-200 via-blue-100 to-indigo-200 opacity-20 rounded-full blur-2xl animate-blob3"></div>
    </div>
    <div class="max-w-7xl mx-auto">
      <!-- Responsive Header -->
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-4 gap-2">
        <h2 class="text-2xl font-bold text-gray-900">Daftar Buku</h2>
        <button
          @click="logout"
          class="w-full sm:w-auto px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700 transition"
        >
          Logout
        </button>
      </div>
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-4 gap-2">
        <div class="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
          <input v-model="search" @input="onSearch" type="text" placeholder="Cari judul/penulis..." class="w-full sm:w-64 px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-indigo-500 transition" />
          <select v-model="selectedCategory" @change="onCategoryChange" class="w-full sm:w-48 px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-indigo-500 transition">
            <option value="">Semua Kategori</option>
            <option value="Fiksi Ilmiah">Fiksi Ilmiah (Science Fiction)</option>
            <option value="Fantasi">Fantasi (Fantasy)</option>
            <option value="Misteri">Misteri (Mystery)</option>
            <option value="Thriller">Thriller</option>
            <option value="Romantis">Romantis (Romance)</option>
            <option value="Horor">Horor (Horror)</option>
            <option value="Sejarah">Sejarah (Historical)</option>
          </select>
        </div>
      </div>
      <div v-if="loading" class="flex justify-center items-center h-40">
        <svg class="animate-spin -ml-1 mr-3 h-8 w-8 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span>Loading...</span>
      </div>
      <div v-else>
        <div v-if="books.length === 0" class="text-center text-gray-500 py-8">
          Tidak ada buku ditemukan
        </div>
        <div class="grid grid-cols-1 xs:grid-cols-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
          <div v-for="book in books" :key="book.id" class="bg-white rounded-xl shadow-lg p-4 flex flex-col transition hover:shadow-2xl hover:-translate-y-1 duration-200 group">
            <div class="relative">
              <img v-if="book.cover" :src="getCoverUrl(book.cover)" alt="Cover" class="h-44 w-full object-cover rounded-lg mb-3 cursor-pointer transition group-hover:scale-105 duration-200" @click="showBookDetail(book)" />
              <div v-else class="h-44 w-full bg-gray-200 rounded-lg mb-3 flex items-center justify-center text-gray-400 text-4xl cursor-pointer" @click="showBookDetail(book)">
                <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
              </div>
              <span v-if="book.stock === 0" class="absolute top-2 right-2 bg-red-500 text-white text-xs px-2 py-1 rounded shadow">Out of Stock</span>
              <span v-if="isNewBook(book)" class="absolute top-2 left-2 bg-green-500 text-white text-xs px-2 py-1 rounded shadow animate-pulse">Baru</span>
            </div>
            <div class="flex-1">
              <h2 class="text-lg font-semibold text-gray-900 mb-1 cursor-pointer transition group-hover:text-indigo-700" @click="showBookDetail(book)">{{ book.title }}</h2>
              <p class="text-xs text-gray-500 mb-1">Author: <span class="font-medium text-gray-700">{{ book.author }}</span></p>
              <p class="text-xs text-gray-500 mb-1">Category: <span class="font-medium text-gray-700">{{ book.category }}</span></p>
              <p class="text-xs text-gray-500 mb-1">ISBN: <span class="font-medium text-gray-700">{{ book.isbn }}</span></p>
              <p class="text-xs text-gray-500 mb-2">Stock: <span :class="book.stock > 0 ? 'text-green-600' : 'text-red-600'">{{ book.stock > 0 ? book.stock : 'Out of Stock' }}</span></p>
            </div>
            <button 
              class="mt-2 w-full bg-gradient-to-r from-indigo-500 to-blue-500 text-white py-2 rounded-lg font-semibold shadow hover:from-indigo-600 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:ring-offset-2 transition disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="book.stock <= 0 || borrowingLoading[book.id]"
              @click="borrowBook(book)"
            >
              <span v-if="borrowingLoading[book.id]">Meminjam...</span>
              <span v-else>Pinjam</span>
            </button>
          </div>
        </div>
        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex justify-center mt-8">
          <nav class="inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
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
    <!-- Modal Detail Buku -->
    <transition name="fade">
      <div v-if="showDetailModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-40">
        <div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 relative mx-2 animate-fadeIn">
          <button class="absolute top-2 right-2 text-gray-400 hover:text-gray-600" @click="showDetailModal = false">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
          </button>
          <div v-if="selectedBook">
            <img v-if="selectedBook.cover" :src="getCoverUrl(selectedBook.cover)" alt="Cover" class="h-48 w-full object-cover rounded mb-4" />
            <h2 class="text-2xl font-bold mb-2 text-indigo-700">{{ selectedBook.title }}</h2>
            <p class="mb-1"><span class="font-semibold">Penulis:</span> {{ selectedBook.author }}</p>
            <p class="mb-1"><span class="font-semibold">Kategori:</span> {{ selectedBook.category }}</p>
            <p class="mb-1"><span class="font-semibold">ISBN:</span> {{ selectedBook.isbn }}</p>
            <p class="mb-1"><span class="font-semibold">Stok:</span> <span :class="selectedBook.stock > 0 ? 'text-green-600' : 'text-red-600'">{{ selectedBook.stock > 0 ? selectedBook.stock : 'Out of Stock' }}</span></p>
            <p class="mb-3"><span class="font-semibold">Deskripsi:</span> <span class="block max-h-32 overflow-y-auto break-words pr-2">{{ selectedBook.description || '-' }}</span></p>
            <button 
              class="w-full bg-gradient-to-r from-indigo-500 to-blue-500 text-white py-2 rounded-lg font-semibold shadow hover:from-indigo-600 hover:to-blue-600 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:ring-offset-2 transition disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="selectedBook.stock <= 0 || borrowingLoading[selectedBook.id]"
              @click="borrowBook(selectedBook, true)"
            >
              <span v-if="borrowingLoading[selectedBook.id]">Meminjam...</span>
              <span v-else>Pinjam Buku Ini</span>
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted, reactive } from 'vue'
import { useStore } from 'vuex'
import { useToast } from 'vue-toastification'
import { useRouter } from 'vue-router'
import { setupSessionTimeout } from '../utils/sessionTimeout'

export default {
  name: 'MemberDashboard',
  setup() {
    const store = useStore()
    const toast = useToast()
    const router = useRouter()
    
    // Setup session timeout
    const { resetTimer, setupActivityListeners, cleanup } = setupSessionTimeout(3)
    
    // Buku
    const books = computed(() => store.getters['books/allBooks'])
    const totalPages = computed(() => store.getters['books/totalPages'])
    const currentPage = computed(() => store.getters['books/currentPage'])
    const pageNumbers = computed(() => {
      const pages = []
      const maxVisible = 5
      let startPage = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
      const endPage = Math.min(totalPages.value, startPage + maxVisible - 1)
      if (endPage - startPage + 1 < maxVisible) {
        startPage = Math.max(1, endPage - maxVisible + 1)
      }
      for (let i = startPage; i <= endPage; i++) {
        pages.push(i)
      }
      return pages
    })
    const loading = ref(true)
    const search = ref('')
    const selectedCategory = ref('')
    const showDetailModal = ref(false)
    const selectedBook = ref(null)
    const borrowingLoading = reactive({})
    const coverBaseUrl = 'http://localhost:8081/static/cover/';
    const getCoverUrl = (cover) => cover ? `${coverBaseUrl}${cover}` : 'default.jpg';

    // Filter kategori dan search
    const fetchBooks = async (page = 1) => {
      loading.value = true
      try {
        await store.dispatch('books/getBooks', { page, search: search.value, category: selectedCategory.value })
      } finally {
        loading.value = false
      }
    }
    const changePage = (page) => {
      if (page < 1 || page > totalPages.value) return
      fetchBooks(page)
    }
    const onSearch = () => {
      fetchBooks(1)
    }
    const onCategoryChange = () => {
      fetchBooks(1)
    }
    const showBookDetail = (book) => {
      selectedBook.value = book
      showDetailModal.value = true
    }
    const borrowBook = async (book, fromModal = false) => {
      if (book.stock <= 0) return
      borrowingLoading[book.id] = true
      try {
        await store.dispatch('borrowings/borrowBook', {
          book_id: book.id,
          due_date: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString() // default 7 hari
        })
        toast.success('Buku berhasil dipinjam!')
        await fetchBooks(currentPage.value)
        if (fromModal) showDetailModal.value = false
      } catch (error) {
        toast.error('Gagal meminjam buku: ' + (error.response?.data?.error || error.message))
      } finally {
        borrowingLoading[book.id] = false
      }
    }
    // Badge Baru: jika buku dibuat dalam 7 hari terakhir
    const isNewBook = (book) => {
      if (!book.created_at) return false
      const created = new Date(book.created_at)
      const now = new Date()
      const diff = (now - created) / (1000 * 60 * 60 * 24)
      return diff <= 7
    }
    const logout = () => {
      store.dispatch('auth/logout')
      router.push('/login')
      toast.success('Logged out successfully')
    }
    onMounted(() => {
      fetchBooks(currentPage.value)
      // Setup session timeout after data is loaded
      setupActivityListeners()
      resetTimer()
    })

    // Cleanup on component unmount
    onUnmounted(() => {
      cleanup()
    })

    return {
      books,
      totalPages,
      currentPage,
      pageNumbers,
      loading,
      search,
      selectedCategory,
      changePage,
      onSearch,
      onCategoryChange,
      showBookDetail,
      showDetailModal,
      selectedBook,
      borrowBook,
      borrowingLoading,
      isNewBook,
      coverBaseUrl,
      getCoverUrl,
      logout
    }
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
.animate-fadeIn {
  animation: fadeIn 0.25s;
}
@keyframes fadeIn {
  from { opacity: 0; transform: scale(0.98); }
  to { opacity: 1; transform: scale(1); }
}
@keyframes blob1 {
  0%, 100% { transform: scale(1) translate(0, 0); }
  33% { transform: scale(1.1) translate(30px, -20px); }
  66% { transform: scale(0.95) translate(-20px, 20px); }
}
@keyframes blob2 {
  0%, 100% { transform: scale(1) translate(0, 0); }
  50% { transform: scale(1.08) translate(-30px, 10px); }
}
@keyframes blob3 {
  0%, 100% { transform: scale(1) translate(0, 0); }
  40% { transform: scale(1.12) translate(10px, -10px); }
  80% { transform: scale(0.9) translate(-10px, 10px); }
}
.animate-blob1 { animation: blob1 12s infinite ease-in-out; }
.animate-blob2 { animation: blob2 14s infinite ease-in-out; }
.animate-blob3 { animation: blob3 16s infinite ease-in-out; }
</style> 