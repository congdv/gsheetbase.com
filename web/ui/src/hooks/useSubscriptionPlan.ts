import { useQuery } from '@tanstack/react-query'
import api from '../lib/axios'

export interface PlanInfo {
    plan: string
    billing_period?: string
    status: string
    subscription_ends_at?: string
    get_rate_limit: number
    update_rate_limit: number
    daily_update_quota: number
    monthly_get_quota: number
    monthly_update_quota: number
    monthly_price_cents: number
    annual_price_cents: number
    cache_min_ttl_seconds: number
    custom_domain: boolean
    priority_support: boolean
}

export const useSubscriptionPlan = () => {
    return useQuery({
        queryKey: ['subscription-plan'],
        queryFn: async () => {
            const response = await api.get<{ plan: PlanInfo }>('/subscription/plan')
            return response.data.plan
        },
        staleTime: 5 * 60 * 1000, // 5 minutes
    })
}
