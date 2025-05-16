<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-100 via-purple-50 to-pink-100 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-2xl shadow-xl transform transition-all duration-300 hover:shadow-2xl">
      <!-- App Logo and Title -->
      <div class="text-center">
        <div class="mx-auto h-16 w-16 bg-indigo-600 rounded-full flex items-center justify-center mb-4">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
        </div>
        <h2 class="text-3xl font-extrabold text-gray-900 tracking-tight">
          Forgot Password
        </h2>
        <p class="mt-2 text-sm text-gray-600">
          <span v-if="step === 1">Masukkan email untuk reset password</span>
          <span v-else>Masukkan PIN yang dikirim ke email & password baru</span>
        </p>
      </div>
      
      <!-- Forgot Password Form -->
      <form class="mt-8 space-y-6" @submit.prevent="handleSubmit">
        <!-- Success Message -->
        <div v-if="successMessage" class="bg-green-50 border-l-4 border-green-400 p-4 rounded-md" role="alert">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-green-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm text-green-700">{{ successMessage }}</p>
            </div>
          </div>
        </div>

        <!-- Error Alert -->
        <div v-if="errorMessage" class="bg-red-50 border-l-4 border-red-400 p-4 rounded-md" role="alert">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
              </svg>
            </div>
            <div class="ml-3">
              <p class="text-sm text-red-700">{{ errorMessage }}</p>
            </div>
          </div>
        </div>
        
        <!-- Email Input -->
        <div v-if="step === 1">
          <div class="space-y-1">
            <label for="email" class="block text-sm font-medium text-gray-700">Email address</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M2.003 5.884L10 9.882l7.997-3.998A2 2 0 0016 4H4a2 2 0 00-1.997 1.884z" />
                  <path d="M18 8.118l-8 4-8-4V14a2 2 0 002 2h12a2 2 0 002-2V8.118z" />
                </svg>
              </div>
              <input 
                id="email" 
                name="email" 
                type="email" 
                autocomplete="email" 
                required 
                v-model="email"
                class="appearance-none block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition duration-150 ease-in-out"
                :class="{ 'border-red-500': errors.email }"
                placeholder="Enter your email"
              />
            </div>
            <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
          </div>
        </div>
        <div v-else>
          <div class="space-y-1">
            <label for="pin" class="block text-sm font-medium text-gray-700">PIN (6 digit)</label>
            <input id="pin" name="pin" type="text" maxlength="6" v-model="pin" class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" :class="{ 'border-red-500': errors.pin }" placeholder="Masukkan PIN" />
            <p v-if="errors.pin" class="mt-1 text-sm text-red-600">{{ errors.pin }}</p>
          </div>
          <div class="space-y-1">
            <label for="newPassword" class="block text-sm font-medium text-gray-700">Password Baru</label>
            <input id="newPassword" name="newPassword" type="password" v-model="newPassword" class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" :class="{ 'border-red-500': errors.newPassword }" placeholder="Password baru" />
            <p v-if="errors.newPassword" class="mt-1 text-sm text-red-600">{{ errors.newPassword }}</p>
          </div>
          <div class="space-y-1">
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Konfirmasi Password Baru</label>
            <input id="confirmPassword" name="confirmPassword" type="password" v-model="confirmPassword" class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" :class="{ 'border-red-500': errors.confirmPassword }" placeholder="Konfirmasi password baru" />
            <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">{{ errors.confirmPassword }}</p>
          </div>
        </div>
        
        <!-- Submit Button -->
        <div>
          <button 
            type="submit"
            :disabled="loading"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition duration-150 ease-in-out transform hover:scale-[1.02]"
            :class="{ 'opacity-75 cursor-not-allowed': loading }"
          >
            <span class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg v-if="!loading" class="h-5 w-5 text-indigo-500 group-hover:text-indigo-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
              </svg>
              <svg v-else class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
            </span>
            <span v-if="loading">Processing...</span>
            <span v-else>{{ step === 1 ? 'Kirim PIN' : 'Reset Password' }}</span>
          </button>
          
          <p class="mt-4 text-center text-sm text-gray-600">
            Ingat password?
            <router-link to="/login" class="font-medium text-indigo-600 hover:text-indigo-500 transition duration-150 ease-in-out">
              Kembali ke login
            </router-link>
          </p>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { useStore } from 'vuex'
import { useToast } from 'vue-toastification'

export default {
  name: 'ForgotPassword',
  setup() {
    const store = useStore()
    const toast = useToast()
    
    // Form data
    const email = ref('')
    const pin = ref('')
    const newPassword = ref('')
    const confirmPassword = ref('')
    const errors = reactive({
      email: '',
      pin: '',
      newPassword: '',
      confirmPassword: ''
    })
    
    // States
    const loading = ref(false)
    const errorMessage = ref('')
    const successMessage = ref('')
    const step = ref(1) // 1: email, 2: pin+password
    
    // Methods
    const validateEmail = () => {
      errors.email = ''
      if (!email.value) {
        errors.email = 'Email is required'
        return false
      } else if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(email.value)) {
        errors.email = 'Email is invalid'
        return false
      }
      return true
    }
    const validatePinPassword = () => {
      errors.pin = ''
      errors.newPassword = ''
      errors.confirmPassword = ''
      let valid = true
      if (!pin.value) {
        errors.pin = 'PIN wajib diisi'
        valid = false
      }
      if (!newPassword.value) {
        errors.newPassword = 'Password baru wajib diisi'
        valid = false
      } else if (newPassword.value.length < 6) {
        errors.newPassword = 'Password minimal 6 karakter'
        valid = false
      }
      if (newPassword.value !== confirmPassword.value) {
        errors.confirmPassword = 'Konfirmasi password tidak sama'
        valid = false
      }
      return valid
    }
    const handleSubmit = async () => {
      if (step.value === 1) {
        if (!validateEmail()) return
        loading.value = true
        try {
          await store.dispatch('auth/forgotPassword', { email: email.value })
          successMessage.value = 'PIN reset password telah dikirim, hubungi admin untuk meminta PIN.'
          toast.success(successMessage.value)
          step.value = 2
        } catch (error) {
          errorMessage.value = error.response?.data?.error || 'Terjadi kesalahan.'
          toast.error(errorMessage.value)
        } finally {
          loading.value = false
        }
      } else if (step.value === 2) {
        if (!validatePinPassword()) return
        loading.value = true
        try {
          // Verifikasi PIN
          await store.dispatch('auth/verifyPin', { email: email.value, pin: pin.value })
          // Reset password
          await store.dispatch('auth/resetPassword', { email: email.value, pin: pin.value, newPassword: newPassword.value })
          successMessage.value = 'Password berhasil direset. Silakan login dengan password baru.'
          toast.success(successMessage.value)
          // Reset form
          pin.value = ''
          newPassword.value = ''
          confirmPassword.value = ''
          setTimeout(() => {
            window.location.href = '/login'
          }, 2000)
        } catch (error) {
          errorMessage.value = error.response?.data?.error || 'PIN salah atau sudah expired.'
          toast.error(errorMessage.value)
        } finally {
          loading.value = false
        }
      }
    }
    
    return {
      email,
      pin,
      newPassword,
      confirmPassword,
      errors,
      loading,
      errorMessage,
      successMessage,
      step,
      handleSubmit
    }
  }
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style> 