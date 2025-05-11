import axios from 'axios'

const state = {
  totalMembers: null,
  loading: false,
  error: null
}

const mutations = {
  SET_TOTAL_MEMBERS(state, count) {
    state.totalMembers = count
  },
  SET_USERS_LOADING(state, loading) {
    state.loading = loading
  },
  SET_USERS_ERROR(state, error) {
    state.error = error
  }
}

const actions = {
  async fetchTotalMembers({ commit, rootState }) {
    commit('SET_USERS_LOADING', true)
    commit('SET_USERS_ERROR', null)
    try {
      const token = rootState.auth.accessToken
      const res = await axios.get('http://localhost:8081/api/v1/users/count', {
        headers: { Authorization: `Bearer ${token}` }
      })
      commit('SET_TOTAL_MEMBERS', res.data.count)
    } catch (err) {
      commit('SET_USERS_ERROR', err.response?.data?.error || err.message)
    } finally {
      commit('SET_USERS_LOADING', false)
    }
  }
}

const getters = {
  totalMembers: state => state.totalMembers,
  usersLoading: state => state.loading,
  usersError: state => state.error
}

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
} 