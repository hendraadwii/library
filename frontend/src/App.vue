<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Main content -->
    <main>
      <router-view />
    </main>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import { useToast } from 'vue-toastification'

export default {
  name: 'App',
  setup() {
    const store = useStore()
    const router = useRouter()
    const toast = useToast()

    // Computed properties
    const isAuthenticated = computed(() => store.getters['auth/isAuthenticated'])
    const currentUser = computed(() => store.getters['auth/currentUser'])

    // Methods
    const logout = () => {
      store.dispatch('auth/logout')
      router.push('/login')
      toast.success('Logged out successfully')
    }

    return {
      isAuthenticated,
      currentUser,
      logout
    }
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Poppins:wght@300;400;500;600;700&display=swap');

* {
  font-family: 'Poppins', sans-serif;
}
</style> 