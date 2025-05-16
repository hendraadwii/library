<template>
  <div class="container mx-auto px-2 sm:px-4 py-4 sm:py-8">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
      <h1 class="text-xl sm:text-2xl font-bold text-gray-800">Buku yang Dipinjam</h1>
      <router-link 
        to="/member-dashboard"
        class="w-full sm:w-auto bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors text-center"
      >
        Kembali ke Daftar Buku
      </router-link>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Error Message -->
    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4">
      {{ error }}
    </div>

    <!-- Borrowings Table -->
    <div v-else-if="borrowings.length > 0" class="bg-white rounded-lg shadow-md overflow-hidden mb-8">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Judul Buku</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal Pinjam</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Batas Kembali</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Denda</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Aksi</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="borrowing in borrowings.filter(b => !b.return_date)" :key="borrowing.id" class="hover:bg-gray-50">
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ borrowing.book_title }}</div>
                <div class="text-sm text-gray-500">Pengarang: {{ borrowing.book_author }}</div>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(borrowing.borrow_date) }}
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(borrowing.due_date) }}
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap">
                <span
                  :class="[
                    'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                    getStatusClass(borrowing)
                  ]"
                >
                  {{ getStatusText(borrowing) }}
                </span>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm">
                <span :class="getFineClass(borrowing)">
                  {{ formatFine(borrowing) }}
                </span>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm">
                <button
                  v-if="!borrowing.return_date"
                  @click="returnBook(borrowing.id)"
                  class="text-green-600 hover:text-green-900 font-semibold"
                >
                  Kembalikan
                </button>
                <span v-else class="text-gray-400">Sudah dikembalikan</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Riwayat Peminjaman -->
    <div v-if="borrowings.filter(b => b.return_date).length > 0" class="bg-white rounded-lg shadow-md overflow-hidden">
      <div class="px-4 py-2 border-b font-semibold text-gray-700">Riwayat Peminjaman</div>
      <div class="overflow-x-auto" style="max-height: 400px; overflow-y: auto;">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Judul Buku</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal Pinjam</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Batas Kembali</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tanggal Kembali</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="borrowing in borrowings.filter(b => b.return_date)" :key="borrowing.id">
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ borrowing.book_title }}</div>
                <div class="text-sm text-gray-500">Pengarang: {{ borrowing.book_author }}</div>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(borrowing.borrow_date) }}
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(borrowing.due_date) }}
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(borrowing.return_date) }}
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap">
                <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">Dikembalikan</span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-8 text-gray-500">
      Belum ada buku yang dipinjam
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useToast } from 'vue-toastification'

export default {
  name: 'MemberBorrowings',
  setup() {
    const store = useStore()
    const toast = useToast()

    // State
    const loading = ref(true)
    const error = ref(null)

    // Computed
    const borrowings = computed(() => store.getters['borrowings/memberBorrowings'] || [])

    // Methods
    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleDateString('id-ID', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    const getStatusClass = (borrowing) => {
      if (borrowing.return_date) return 'bg-green-100 text-green-800'
      if (isOverdue(borrowing)) return 'bg-red-100 text-red-800'
      return 'bg-blue-100 text-blue-800'
    }

    const getStatusText = (borrowing) => {
      if (borrowing.return_date) return 'Dikembalikan'
      if (isOverdue(borrowing)) return 'Terlambat'
      return 'Dipinjam'
    }

    const getFineClass = (borrowing) => {
      const fine = calculateFine(borrowing)
      if (fine > 0) return 'text-red-600 font-medium'
      return 'text-gray-500'
    }

    const formatFine = (borrowing) => {
      const fine = calculateFine(borrowing)
      if (fine > 0) return `Rp ${fine.toLocaleString('id-ID')}`
      return '-'
    }

    const isOverdue = (borrowing) => {
      if (borrowing.return_date) return false
      const dueDate = new Date(borrowing.due_date)
      const now = new Date()
      return now > dueDate
    }

    const calculateFine = (borrowing) => {
      if (borrowing.return_date) return 0
      const dueDate = new Date(borrowing.due_date)
      const now = new Date()
      if (now <= dueDate) return 0
      
      const diffTime = Math.abs(now - dueDate)
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
      return diffDays * 1000 // Rp 1.000 per hari keterlambatan
    }

    const returnBook = async (borrowingId) => {
      if (!borrowingId) {
        toast.error('ID peminjaman tidak valid');
        return;
      }
      try {
        await store.dispatch('borrowings/returnBook', borrowingId)
        toast.success('Buku berhasil dikembalikan')
        // Refresh data
        await store.dispatch('borrowings/getMemberBorrowings')
      } catch (err) {
        toast.error(err.response?.data?.error || 'Gagal mengembalikan buku')
      }
    }

    const fetchBorrowings = async () => {
      loading.value = true
      error.value = null
      console.log('fetchBorrowings called')
      try {
        await store.dispatch('borrowings/getMemberBorrowings')
        console.log('fetchBorrowings success')
      } catch (err) {
        error.value = err.response?.data?.error || 'Gagal memuat data peminjaman'
        toast.error(error.value)
        console.error('fetchBorrowings error:', error.value)
      } finally {
        loading.value = false
        console.log('fetchBorrowings finished, loading:', loading.value)
      }
    }

    onMounted(() => {
      fetchBorrowings()
    })

    return {
      loading,
      error,
      borrowings,
      formatDate,
      getStatusClass,
      getStatusText,
      getFineClass,
      formatFine,
      returnBook
    }
  }
}
</script> 