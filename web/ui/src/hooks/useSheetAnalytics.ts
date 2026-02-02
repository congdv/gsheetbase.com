import { useQuery } from '@tanstack/react-query'
import api from '../lib/axios'

export interface DailyUsageSummary {
  date: string
  total_count: number
  get_count: number
  post_count: number
  put_count: number
  patch_count: number
  delete_count: number
}

export interface SheetAnalyticsResponse {
  sheet_id: string
  sheet_name: string
  api_key: string
  period_days: number
  start_date: string
  end_date: string
  daily_usage: DailyUsageSummary[]
}

export const useSheetAnalytics = (sheetId: string, days: number = 30) => {
  return useQuery({
    queryKey: ['sheet-analytics', sheetId, days],
    queryFn: async () => {
      const response = await api.get<SheetAnalyticsResponse>(
        `/sheets/${sheetId}/analytics`,
        { params: { days } }
      )
      return response.data
    },
    enabled: !!sheetId,
  })
}
