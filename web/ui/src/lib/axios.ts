import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'

// Create axios instance with empty baseURL initially
// It will be configured in main.tsx after fetching runtime config
const api = axios.create({
  withCredentials: true,
})

// Track refresh request to prevent concurrent refresh attempts
let isRefreshing = false
let refreshSubscribers: Array<(token?: string) => void> = []

// Subscribe to refresh completion
const subscribeRefresh = (cb: (token?: string) => void) => {
  refreshSubscribers.push(cb)
}

// Notify all subscribers that refresh is done
const onRefreshed = () => {
  refreshSubscribers.forEach(cb => cb())
  refreshSubscribers = []
}

// Add response interceptor to handle 401 and refresh session
api.interceptors.response.use(
  response => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retried?: boolean }

    // Only handle 401 errors for non-refresh endpoints
    if (error.response?.status === 401 && originalRequest && !originalRequest._retried && originalRequest.url !== '/auth/refresh-session') {
      // Mark request as retried to prevent infinite loops
      originalRequest._retried = true

      if (!isRefreshing) {
        isRefreshing = true

        try {
          // Attempt to refresh session
          const response = await axios.post(
            `${api.defaults.baseURL || ''}/api/auth/refresh-session`,
            {},
            { withCredentials: true }
          )

          isRefreshing = false
          onRefreshed()

          // Retry original request with new token
          return api(originalRequest)
        } catch (refreshError) {
          // Refresh failed - session is truly expired or invalid
          isRefreshing = false
          refreshSubscribers = []

          // Let the original error propagate
          return Promise.reject(error)
        }
      } else {
        // Refresh is in progress, queue this request
        return new Promise(resolve => {
          subscribeRefresh(() => {
            resolve(api(originalRequest))
          })
        })
      }
    }

    // For other errors or already-retried requests, reject normally
    return Promise.reject(error)
  }
)

export default api

// Function to configure axios with runtime config
export function configureAxios(baseURL: string) {
  api.defaults.baseURL = baseURL
}
