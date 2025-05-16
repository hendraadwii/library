import { createRouter, createWebHistory } from 'vue-router'
import store from '../store'

// Lazy-loaded components
const Home = () => import('../views/Home.vue')
const Login = () => import('../views/Login.vue')
const Register = () => import('../views/Register.vue')
const ForgotPassword = () => import('../views/ForgotPassword.vue')
const Books = () => import('../views/Books.vue')
const BookDetails = () => import('../views/BookDetails.vue')
const Borrowings = () => import('../views/Borrowings.vue')
const NotFound = () => import('../views/NotFound.vue')
const MemberDashboard = () => import('../views/MemberDashboard.vue')
const Users = () => import('../views/Users.vue')

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { guestOnly: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { guestOnly: true }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: ForgotPassword,
    meta: { guestOnly: true }
  },
  {
    path: '/books',
    name: 'Books',
    component: Books,
    meta: { requiresAuth: true }
  },
  {
    path: '/books/:id',
    name: 'BookDetails',
    component: BookDetails,
    meta: { requiresAuth: true },
    props: true
  },
  {
    path: '/borrowings',
    name: 'Borrowings',
    component: Borrowings,
    meta: { requiresAuth: true }
  },
  {
    path: '/member-dashboard',
    name: 'MemberDashboard',
    component: MemberDashboard,
    meta: { requiresAuth: true, requiresMember: true }
  },
  {
    path: '/users',
    name: 'Users',
    component: Users,
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  // Selalu coba rehydrate auth dari localStorage jika belum authenticated
  if (!store.getters['auth/isAuthenticated']) {
    await store.dispatch('auth/initAuth')
  }
  const isAuthenticated = store.getters['auth/isAuthenticated']
  const currentUser = store.getters['auth/currentUser']

  // Route requires authentication but user is not authenticated
  if (to.matched.some(record => record.meta.requiresAuth) && !isAuthenticated) {
    next('/login')
    return
  }

  // Route is for guests only but user is authenticated
  if (to.matched.some(record => record.meta.guestOnly) && isAuthenticated) {
    // Redirect ke dashboard sesuai role
    if (currentUser && currentUser.role === 'admin') {
      next('/')
    } else if (currentUser && currentUser.role === 'member') {
      next('/member-dashboard')
    } else {
      next('/')
    }
    return
  }

  // Admin only
  if (to.meta.requiresAdmin) {
    if (!isAuthenticated || !currentUser || currentUser.role !== 'admin') {
      // Jika member mencoba akses Home, redirect ke member-dashboard
      if (currentUser && currentUser.role === 'member') {
        return next({ name: 'MemberDashboard' })
      }
      return next({ name: 'Login' })
    }
  }
  // Member only
  if (to.meta.requiresMember) {
    if (!isAuthenticated || !currentUser || currentUser.role !== 'member') {
      // Jika admin mencoba akses member-dashboard, redirect ke Home
      if (currentUser && currentUser.role === 'admin') {
        return next({ name: 'Home' })
      }
      return next({ name: 'Login' })
    }
  }

  next()
})

export default router 