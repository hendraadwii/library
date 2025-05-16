import { createStore } from 'vuex'
import auth from './modules/auth'
import books from './modules/books'
import borrowings from './modules/borrowings'
import users from './modules/users'

const store = createStore({
  state: {
    loading: false,
    error: null
  },
  mutations: {
    SET_LOADING(state, loading) {
      state.loading = loading
    },
    SET_ERROR(state, error) {
      state.error = error
    },
    CLEAR_ERROR(state) {
      state.error = null
    }
  },
  actions: {
    setLoading({ commit }, loading) {
      commit('SET_LOADING', loading)
    },
    setError({ commit }, error) {
      commit('SET_ERROR', error)
    },
    clearError({ commit }) {
      commit('CLEAR_ERROR')
    }
  },
  getters: {
    loading: state => state.loading,
    error: state => state.error
  },
  modules: {
    auth,
    books,
    borrowings,
    users
  }
})

export default store 