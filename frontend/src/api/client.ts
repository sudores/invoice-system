import axios from 'axios'
import router from '../router'
import { useUserStore } from '../stores/user'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
})

export async function refreshToken() {
  const refreshToken = localStorage.getItem('refreshToken')
  if (!refreshToken) return null

  try {
    const resp = await api.post('/login/refresh', { refresh_token: refreshToken })
    const data = resp.data

    if (data?.jwt) {
      localStorage.setItem('token', data.jwt)
      return data.jwt
    }
  } catch (err) {
    console.error('Refresh token failed', err)
  }
  return null
}



api.interceptors.request.use((config) => {
  const user = useUserStore()
  if (user.token) {
    config.headers.Authorization = `Bearer ${user.token}`
  }
  return config
})

api.interceptors.response.use(
  res => res,
  async err => {
    const originalRequest = err.config

    if (err.response && err.response.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      const newToken = await refreshToken()
      if (newToken) {
        originalRequest.headers['Authorization'] = `Bearer ${newToken}`
        return api(originalRequest) // retry the original request
      }

      // If refresh fails, redirect to login
      localStorage.removeItem('token')
      localStorage.removeItem('refresh_token')
      router.push('/login')
    }

    return Promise.reject(err)
  }
)


export default api

