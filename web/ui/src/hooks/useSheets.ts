import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { message } from 'antd'
import api from '../lib/axios'

export interface Sheet {
  id: string
  user_id: string
  sheet_id: string
  sheet_name: string
  description: string
  is_public: boolean
  api_key?: string
  default_range?: string
  use_first_row_as_header: boolean
  created_at: string
}

interface RegisterSheetParams {
  sheet_id: string
  sheet_name?: string
  description?: string
}

interface PublishSheetParams {
  sheetId: string
  values: {
    default_range?: string
    use_first_row_as_header: boolean
  }
}

export const useSheets = () => {
  const queryClient = useQueryClient()

  // Fetch all registered sheets
  const {
    data: sheets = [],
    isLoading,
    error,
  } = useQuery({
    queryKey: ['sheets'],
    queryFn: async () => {
      const response = await api.get('/sheets/registered')
      return (response.data.sheets || []) as Sheet[]
    },
  })

  // Register new sheet
  const registerMutation = useMutation({
    mutationFn: async (values: RegisterSheetParams) => {
      const response = await api.post('/sheets/register', values)
      return response.data
    },
    onSuccess: () => {
      message.success('Sheet registered successfully!')
      queryClient.invalidateQueries({ queryKey: ['sheets'] })
    },
    onError: (error: any) => {
      message.error(error?.response?.data?.error || 'Failed to register sheet')
    },
  })

  // Delete sheet
  const deleteMutation = useMutation({
    mutationFn: async (sheetId: string) => {
      await api.delete(`/sheets/registered/${sheetId}`)
    },
    onSuccess: () => {
      message.success('Sheet removed successfully!')
      queryClient.invalidateQueries({ queryKey: ['sheets'] })
    },
    onError: (error: any) => {
      message.error(error?.response?.data?.error || 'Failed to remove sheet')
    },
  })

  // Publish sheet
  const publishMutation = useMutation({
    mutationFn: async ({ sheetId, values }: PublishSheetParams) => {
      const response = await api.post(`/sheets/${sheetId}/publish`, values)
      return response.data
    },
    onSuccess: () => {
      message.success('Sheet published successfully!')
      queryClient.invalidateQueries({ queryKey: ['sheets'] })
    },
    onError: (error: any) => {
      message.error(error?.response?.data?.error || 'Failed to publish sheet')
    },
  })

  // Unpublish sheet
  const unpublishMutation = useMutation({
    mutationFn: async (sheetId: string) => {
      await api.delete(`/sheets/${sheetId}/unpublish`)
    },
    onSuccess: () => {
      message.success('Sheet unpublished successfully!')
      queryClient.invalidateQueries({ queryKey: ['sheets'] })
    },
    onError: (error: any) => {
      message.error(error?.response?.data?.error || 'Failed to unpublish sheet')
    },
  })

  return {
    sheets,
    isLoading,
    error,
    registerSheet: registerMutation.mutate,
    isRegistering: registerMutation.isPending,
    deleteSheet: deleteMutation.mutate,
    isDeleting: deleteMutation.isPending,
    publishSheet: publishMutation.mutate,
    isPublishing: publishMutation.isPending,
    unpublishSheet: unpublishMutation.mutate,
    isUnpublishing: unpublishMutation.isPending,
  }
}
