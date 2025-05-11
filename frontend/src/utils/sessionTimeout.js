import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'

export const setupSessionTimeout = (timeoutMinutes = 3) => {
  const store = useStore()
  const router = useRouter()
  const toast = useToast()
  
  let timeoutId = null
  const timeoutMs = timeoutMinutes * 60 * 1000

  const resetTimer = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
    timeoutId = setTimeout(() => {
      // Show notification
      toast.warning('Session telah berakhir', {
        position: 'top-center',
        timeout: 2000
      })
      
      // Logout and redirect
      store.dispatch('auth/logout')
      router.push('/login')
    }, timeoutMs)
  }

  // Reset timer on user activity
  const setupActivityListeners = () => {
    const events = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart']
    events.forEach(event => {
      document.addEventListener(event, resetTimer)
    })
  }

  // Clean up event listeners
  const cleanup = () => {
    if (timeoutId) {
      clearTimeout(timeoutId)
    }
    const events = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart']
    events.forEach(event => {
      document.removeEventListener(event, resetTimer)
    })
  }

  return {
    resetTimer,
    setupActivityListeners,
    cleanup
  }
} 