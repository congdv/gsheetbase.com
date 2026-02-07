import axios from 'axios'

// Create axios instance with empty baseURL initially
// It will be configured in main.tsx after fetching runtime config
const api = axios.create({
  withCredentials: true,
})

export default api

// Function to configure axios with runtime config
export function configureAxios(baseURL: string) {
  api.defaults.baseURL = baseURL
}
