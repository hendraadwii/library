<template>
  <div class="container mx-auto px-2 sm:px-4 py-4 sm:py-8">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
      <h1 class="text-xl sm:text-2xl font-bold text-gray-800">User Management</h1>
      <button
        @click="openAddUserModal"
        class="w-full sm:w-auto bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 transition-colors"
      >
        Add New User
      </button>
    </div>

    <!-- Search and Filters -->
    <div class="bg-white rounded-lg shadow-md p-3 sm:p-4 mb-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <input
            v-model="search"
            type="text"
            placeholder="Search by name or email..."
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            @input="debouncedSearch"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
          <select
            v-model="roleFilter"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            @change="fetchUsers"
          >
            <option value="">All Roles</option>
            <option value="admin">Admin</option>
            <option value="member">Member</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            v-model="statusFilter"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            @change="fetchUsers"
          >
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="locked">Locked</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-8">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error Message -->
    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4">
      {{ error }}
    </div>

    <!-- Users Table -->
    <div v-else class="bg-white rounded-lg shadow-md overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created At</th>
              <th class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50">
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ user.full_name }}</div>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-500">{{ user.email }}</div>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap">
                <span
                  :class="[
                    'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                    user.role === 'admin' ? 'bg-purple-100 text-purple-800' : 'bg-green-100 text-green-800'
                  ]"
                >
                  {{ user.role }}
                </span>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap">
                <span
                  :class="[
                    'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                    user.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                  ]"
                >
                  {{ user.status }}
                </span>
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {{ formatDate(user.created_at) }}
              </td>
              <td class="px-3 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">
                <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
                  <button
                    @click="openEditUserModal(user)"
                    class="text-indigo-600 hover:text-indigo-900"
                  >
                    Edit
                  </button>
                  <button
                    @click="openResetPasswordModal(user)"
                    class="text-yellow-600 hover:text-yellow-900"
                  >
                    Reset Password
                  </button>
                  <button
                    @click="toggleUserStatus(user)"
                    :class="[
                      user.status === 'active' ? 'text-red-600 hover:text-red-900' : 'text-green-600 hover:text-green-900'
                    ]"
                  >
                    {{ user.status === 'active' ? 'Lock' : 'Unlock' }}
                  </button>
                  <button
                    @click="openDeleteModal(user)"
                    class="text-red-600 hover:text-red-900"
                  >
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="bg-white px-2 sm:px-4 py-3 flex flex-col sm:flex-row items-center justify-between border-t border-gray-200 gap-4">
        <div class="w-full sm:w-auto flex justify-between sm:hidden">
          <button
            @click="changePage(currentPage - 1)"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <button
            @click="changePage(currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              Showing
              <span class="font-medium">{{ (currentPage - 1) * perPage + 1 }}</span>
              to
              <span class="font-medium">{{ Math.min(currentPage * perPage, totalUsers) }}</span>
              of
              <span class="font-medium">{{ totalUsers }}</span>
              results
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
              <button
                v-for="page in displayedPages"
                :key="page"
                @click="changePage(page)"
                :class="[
                  'relative inline-flex items-center px-3 sm:px-4 py-2 border text-sm font-medium',
                  page === currentPage
                    ? 'z-10 bg-indigo-50 border-indigo-500 text-indigo-600'
                    : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                ]"
              >
                {{ page }}
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit User Modal -->
    <div v-if="showAddUserModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
      <div class="relative top-10 sm:top-20 mx-auto p-4 sm:p-5 border w-full sm:w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium leading-6 text-gray-900 mb-4">
            {{ editingUser ? 'Edit User' : 'Add New User' }}
          </h3>
          <form @submit.prevent="saveUser">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
              <input
                v-model="userForm.name"
                type="text"
                required
                minlength="2"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
              <input
                v-model="userForm.email"
                type="email"
                required
                pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Role</label>
              <select
                v-model="userForm.role"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="member">Member</option>
                <option value="admin">Admin</option>
              </select>
            </div>
            <div v-if="!editingUser" class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
              <input
                v-model="userForm.password"
                type="password"
                required
                minlength="6"
                pattern="(?=.*\d).{6,}"
                title="Password must be at least 6 characters long and contain at least one number"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div class="flex flex-col sm:flex-row justify-end gap-2 sm:gap-3">
              <button
                type="button"
                @click="showAddUserModal = false"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md"
              >
                {{ editingUser ? 'Update' : 'Create' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Reset Password Modal -->
    <div v-if="showResetPasswordModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
      <div class="relative top-10 sm:top-20 mx-auto p-4 sm:p-5 border w-full sm:w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium leading-6 text-gray-900 mb-4">Reset Password</h3>
          <form @submit.prevent="resetPassword">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
              <input
                v-model="resetPasswordForm.password"
                type="password"
                required
                minlength="6"
                pattern="(?=.*\d).{6,}"
                title="Password must be at least 6 characters long and contain at least one number"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
            <div class="flex flex-col sm:flex-row justify-end gap-2 sm:gap-3">
              <button
                type="button"
                @click="showResetPasswordModal = false"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-md"
              >
                Reset Password
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
      <div class="relative top-10 sm:top-20 mx-auto p-4 sm:p-5 border w-full sm:w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium leading-6 text-gray-900 mb-4">Confirm Delete</h3>
          <p class="text-sm text-gray-500 mb-4">
            Are you sure you want to delete user "{{ userToDelete?.full_name }}"? This action cannot be undone.
          </p>
          <div class="flex flex-col sm:flex-row justify-end gap-2 sm:gap-3">
            <button
              @click="showDeleteModal = false"
              class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md"
            >
              Cancel
            </button>
            <button
              @click="deleteUser"
              class="w-full sm:w-auto px-4 py-2 text-sm font-medium text-white bg-red-600 hover:bg-red-700 rounded-md"
            >
              Delete
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
import debounce from 'lodash/debounce'

export default {
  name: 'Users',
  setup() {
    const store = useStore()
    const toast = useToast()

    // State
    const search = ref('')
    const roleFilter = ref('')
    const statusFilter = ref('')
    const showAddUserModal = ref(false)
    const showResetPasswordModal = ref(false)
    const showDeleteModal = ref(false)
    const editingUser = ref(null)
    const userToDelete = ref(null)
    const userForm = reactive({
      name: '',
      email: '',
      role: 'member',
      password: ''
    })
    const resetPasswordForm = reactive({
      userId: null,
      password: ''
    })

    // Computed
    const users = computed(() => store.state.users.users)
    const totalUsers = computed(() => store.state.users.totalUsers)
    const currentPage = computed(() => store.state.users.currentPage)
    const perPage = computed(() => store.state.users.perPage)
    const totalPages = computed(() => store.state.users.totalPages)
    const loading = computed(() => store.state.users.loading)
    const error = computed(() => store.state.users.error)

    const displayedPages = computed(() => {
      const pages = []
      const maxPages = 5
      let start = Math.max(1, currentPage.value - Math.floor(maxPages / 2))
      let end = Math.min(totalPages.value, start + maxPages - 1)
      
      if (end - start + 1 < maxPages) {
        start = Math.max(1, end - maxPages + 1)
      }
      
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      
      return pages
    })

    // Methods
    const fetchUsers = () => {
      store.dispatch('users/getUsers', {
        page: currentPage.value,
        perPage: perPage.value,
        search: search.value,
        role: roleFilter.value,
        status: statusFilter.value
      })
    }

    const debouncedSearch = debounce(() => {
      fetchUsers()
    }, 300)

    const formatDate = (dateString) => {
      if (!dateString) return '-'
      const date = new Date(dateString)
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })
    }

    const openAddUserModal = () => {
      editingUser.value = null
      userForm.name = ''
      userForm.email = ''
      userForm.role = 'member'
      userForm.password = ''
      showAddUserModal.value = true
    }

    const openEditUserModal = (user) => {
      editingUser.value = user
      userForm.name = user.full_name
      userForm.email = user.email
      userForm.role = user.role
      userForm.password = ''
      showAddUserModal.value = true
    }

    const openResetPasswordModal = (user) => {
      resetPasswordForm.userId = user.id
      resetPasswordForm.password = ''
      showResetPasswordModal.value = true
    }

    const openDeleteModal = (user) => {
      userToDelete.value = user
      showDeleteModal.value = true
    }

    const saveUser = async () => {
      try {
        if (editingUser.value) {
          await store.dispatch('users/updateUser', {
            id: editingUser.value.id,
            ...userForm
          })
          toast.success('User updated successfully')
        } else {
          await store.dispatch('users/createUser', userForm)
          toast.success('User created successfully')
        }
        showAddUserModal.value = false
        fetchUsers()
      } catch (error) {
        toast.error(error.response?.data?.error || 'An error occurred')
      }
    }

    const resetPassword = async () => {
      try {
        await store.dispatch('users/resetUserPassword', {
          userId: resetPasswordForm.userId,
          password: resetPasswordForm.password
        })
        toast.success('Password reset successfully')
        showResetPasswordModal.value = false
      } catch (error) {
        toast.error(error.response?.data?.error || 'An error occurred')
      }
    }

    const deleteUser = async () => {
      try {
        await store.dispatch('users/deleteUser', userToDelete.value.id)
        toast.success('User deleted successfully')
        showDeleteModal.value = false
        fetchUsers()
      } catch (error) {
        toast.error(error.response?.data?.error || 'An error occurred')
      }
    }

    const toggleUserStatus = async (user) => {
      try {
        await store.dispatch('users/toggleUserStatus', {
          userId: user.id,
          currentStatus: user.status
        })
        toast.success(`User ${user.status === 'active' ? 'locked' : 'unlocked'} successfully`)
        fetchUsers()
      } catch (error) {
        toast.error(error.response?.data?.error || 'An error occurred')
      }
    }

    const changePage = (page) => {
      store.commit('users/SET_CURRENT_PAGE', page)
      fetchUsers()
    }

    // Lifecycle hooks
    onMounted(() => {
      fetchUsers()
    })

    return {
      // State
      search,
      roleFilter,
      statusFilter,
      showAddUserModal,
      showResetPasswordModal,
      showDeleteModal,
      editingUser,
      userToDelete,
      userForm,
      resetPasswordForm,

      // Computed
      users,
      totalUsers,
      currentPage,
      perPage,
      totalPages,
      loading,
      error,
      displayedPages,

      // Methods
      fetchUsers,
      debouncedSearch,
      formatDate,
      openAddUserModal,
      openEditUserModal,
      openResetPasswordModal,
      openDeleteModal,
      saveUser,
      resetPassword,
      deleteUser,
      toggleUserStatus,
      changePage
    }
  }
}
</script>