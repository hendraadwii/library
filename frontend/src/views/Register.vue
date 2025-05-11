<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-indigo-100 via-purple-50 to-pink-100 py-12 px-4 sm:px-6 lg:px-8 animate-gradient">
    <div class="max-w-md w-full space-y-8 bg-white p-8 rounded-2xl shadow-xl transform transition-all duration-300 hover:shadow-2xl animate-fade-in">
      <!-- App Logo and Title -->
      <div class="text-center animate-slide-down">
        <div class="mx-auto h-16 w-16 bg-indigo-600 rounded-full flex items-center justify-center mb-4 animate-bounce-slow">
          <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
          </svg>
        </div>
        <h2 class="text-3xl font-extrabold text-gray-900 tracking-tight animate-fade-in">
          Create Account
        </h2>
        <p class="mt-2 text-sm text-gray-600 animate-fade-in-delay">
          Join our library community
        </p>
      </div>
      
      <!-- Registration Form -->
      <form class="mt-8 space-y-6 animate-slide-up" @submit.prevent="register">
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
        
        <!-- Full Name Input -->
        <div class="space-y-1">
          <label for="full_name" class="block text-sm font-medium text-gray-700">Full Name</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd" />
              </svg>
            </div>
            <input 
              id="full_name" 
              name="full_name" 
              type="text" 
              required 
              v-model="form.full_name"
              class="appearance-none block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition duration-150 ease-in-out"
              :class="{ 'border-red-500': errors.full_name }"
              placeholder="Enter your full name"
            />
          </div>
          <p v-if="errors.full_name" class="mt-1 text-sm text-red-600">{{ errors.full_name }}</p>
        </div>
        
        <!-- Email Input -->
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
              v-model="form.email"
              class="appearance-none block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition duration-150 ease-in-out"
              :class="{ 'border-red-500': errors.email }"
              placeholder="Enter your email"
            />
          </div>
          <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
        </div>
        
        <!-- Password Input -->
        <div class="space-y-1">
          <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
              </svg>
            </div>
            <input 
              id="password" 
              name="password" 
              :type="showPassword ? 'text' : 'password'"
              autocomplete="new-password" 
              required 
              v-model="form.password"
              class="appearance-none block w-full pl-10 pr-10 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition duration-150 ease-in-out"
              :class="{ 'border-red-500': errors.password }"
              placeholder="Enter your password"
              @input="validatePassword"
            />
            <div class="absolute inset-y-0 right-0 pr-3 flex items-center">
              <button 
                type="button"
                @click="showPassword = !showPassword"
                class="text-gray-400 hover:text-gray-500 focus:outline-none"
              >
                <svg v-if="showPassword" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clip-rule="evenodd" />
                  <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.065 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                </svg>
              </button>
            </div>
          </div>
          <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
          
          <!-- Password Strength Indicator -->
          <div class="mt-2">
            <div class="flex items-center space-x-2">
              <div class="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
                <div 
                  class="h-full transition-all duration-300 ease-in-out password-strength-bar"
                  :class="{
                    'w-1/4 bg-red-500': passwordStrength === 'weak',
                    'w-2/4 bg-yellow-500': passwordStrength === 'medium',
                    'w-3/4 bg-blue-500': passwordStrength === 'strong',
                    'w-full bg-green-500': passwordStrength === 'very-strong'
                  }"
                ></div>
              </div>
              <span class="text-xs font-medium" :class="{
                'text-red-500': passwordStrength === 'weak',
                'text-yellow-500': passwordStrength === 'medium',
                'text-blue-500': passwordStrength === 'strong',
                'text-green-500': passwordStrength === 'very-strong'
              }">{{ passwordStrength }}</span>
            </div>
            <p class="mt-1 text-xs text-gray-500">
              Password must contain at least 8 characters, including uppercase, lowercase, numbers, and special characters
            </p>
          </div>
        </div>
        
        <!-- Confirm Password Input -->
        <div class="space-y-1">
          <label for="confirm_password" class="block text-sm font-medium text-gray-700">Confirm Password</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clip-rule="evenodd" />
              </svg>
            </div>
            <input 
              id="confirm_password" 
              name="confirm_password" 
              :type="showConfirmPassword ? 'text' : 'password'"
              autocomplete="new-password" 
              required 
              v-model="form.confirm_password"
              class="appearance-none block w-full pl-10 pr-10 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition duration-150 ease-in-out"
              :class="{ 'border-red-500': errors.confirm_password }"
              placeholder="Confirm your password"
            />
            <div class="absolute inset-y-0 right-0 pr-3 flex items-center">
              <button 
                type="button"
                @click="showConfirmPassword = !showConfirmPassword"
                class="text-gray-400 hover:text-gray-500 focus:outline-none"
              >
                <svg v-if="showConfirmPassword" class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
                <svg v-else class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clip-rule="evenodd" />
                  <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.065 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                </svg>
              </button>
            </div>
          </div>
          <p v-if="errors.confirm_password" class="mt-1 text-sm text-red-600">{{ errors.confirm_password }}</p>
        </div>
        
        <!-- Submit Button and Login Link -->
        <div>
          <button 
            type="submit"
            :disabled="loading || !isPasswordValid"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition duration-150 ease-in-out transform hover:scale-[1.02]"
            :class="{ 'opacity-75 cursor-not-allowed': loading || !isPasswordValid }"
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
            <span v-else>Create Account</span>
          </button>
          
          <p class="mt-4 text-center text-sm text-gray-600">
            Already have an account?
            <router-link to="/login" class="font-medium text-indigo-600 hover:text-indigo-500 transition duration-150 ease-in-out">
              Sign in
            </router-link>
          </p>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'

export default {
  name: 'Register',
  setup() {
    const store = useStore()
    const router = useRouter()
    const toast = useToast()
    
    // Form data
    const form = reactive({
      full_name: '',
      email: '',
      password: '',
      confirm_password: ''
    })
    
    // Form validation
    const errors = reactive({
      full_name: '',
      email: '',
      password: '',
      confirm_password: ''
    })
    
    // States
    const loading = ref(false)
    const errorMessage = ref('')
    const showPassword = ref(false)
    const showConfirmPassword = ref(false)
    const passwordStrength = ref('weak')
    
    // Computed properties
    const isAuthenticated = computed(() => store.getters['auth/isAuthenticated'])
    const isPasswordValid = computed(() => {
      return passwordStrength.value === 'strong' || passwordStrength.value === 'very-strong'
    })
    
    // If already authenticated, redirect to home
    if (isAuthenticated.value) {
      router.push('/')
    }
    
    // Methods
    const validatePassword = () => {
      const password = form.password
      let strength = 0
      
      // Length check
      if (password.length >= 8) strength++
      
      // Uppercase check
      if (/[A-Z]/.test(password)) strength++
      
      // Lowercase check
      if (/[a-z]/.test(password)) strength++
      
      // Number check
      if (/[0-9]/.test(password)) strength++
      
      // Special character check
      if (/[^A-Za-z0-9]/.test(password)) strength++
      
      // Set password strength
      if (strength <= 2) passwordStrength.value = 'weak'
      else if (strength === 3) passwordStrength.value = 'medium'
      else if (strength === 4) passwordStrength.value = 'strong'
      else passwordStrength.value = 'very-strong'
    }
    
    const validateForm = () => {
      let isValid = true
      
      // Reset errors
      errors.full_name = ''
      errors.email = ''
      errors.password = ''
      errors.confirm_password = ''
      errorMessage.value = ''
      
      // Validate full name
      if (!form.full_name) {
        errors.full_name = 'Full name is required'
        isValid = false
      }
      
      // Validate email
      if (!form.email) {
        errors.email = 'Email is required'
        isValid = false
      } else if (!/^[\w-]+(\.[\w-]+)*@([\w-]+\.)+[a-zA-Z]{2,7}$/.test(form.email)) {
        errors.email = 'Email is invalid'
        isValid = false
      }
      
      // Validate password
      if (!form.password) {
        errors.password = 'Password is required'
        isValid = false
      } else if (!isPasswordValid.value) {
        errors.password = 'Password is too weak'
        isValid = false
      }
      
      // Validate password confirmation
      if (!form.confirm_password) {
        errors.confirm_password = 'Please confirm your password'
        isValid = false
      } else if (form.password !== form.confirm_password) {
        errors.confirm_password = 'Passwords do not match'
        isValid = false
      }
      
      return isValid
    }
    
    const register = async () => {
      if (!validateForm()) return
      
      loading.value = true
      
      try {
        // Call register action
        await store.dispatch('auth/register', {
          email: form.email,
          password: form.password,
          full_name: form.full_name
        })
        
        // Logout user after successful registration
        await store.dispatch('auth/logout')
        
        // Show success message
        toast.success('Account created successfully! Please login to continue.')
        
        // Redirect to login page
        router.push('/login')
      } catch (error) {
        // Set error message
        if (error.response && error.response.data && error.response.data.error) {
          errorMessage.value = error.response.data.error
        } else {
          errorMessage.value = 'An error occurred. Please try again.'
        }
        
        // Show error message
        toast.error(errorMessage.value)
      } finally {
        loading.value = false
      }
    }
    
    return {
      form,
      errors,
      loading,
      errorMessage,
      showPassword,
      showConfirmPassword,
      passwordStrength,
      isPasswordValid,
      validatePassword,
      register
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

/* New animations */
@keyframes gradient {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.animate-gradient {
  background-size: 200% 200%;
  animation: gradient 15s ease infinite;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.animate-fade-in {
  animation: fadeIn 0.5s ease-out;
}

.animate-fade-in-delay {
  animation: fadeIn 0.5s ease-out 0.2s both;
}

@keyframes slideDown {
  from {
    transform: translateY(-20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.animate-slide-down {
  animation: slideDown 0.5s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.animate-slide-up {
  animation: slideUp 0.5s ease-out;
}

@keyframes bounceSlow {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.animate-bounce-slow {
  animation: bounceSlow 2s ease-in-out infinite;
}

/* Input field animations */
input:focus {
  transform: scale(1.01);
  transition: transform 0.2s ease-in-out;
}

/* Button hover animation */
button:not(:disabled):hover {
  transform: translateY(-1px);
  transition: transform 0.2s ease-in-out;
}

/* Password strength indicator animation */
.password-strength-bar {
  transition: width 0.3s ease-in-out, background-color 0.3s ease-in-out;
}

/* Error message animation */
.error-message {
  animation: shake 0.5s ease-in-out;
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  25% {
    transform: translateX(-5px);
  }
  75% {
    transform: translateX(5px);
  }
}
</style> 