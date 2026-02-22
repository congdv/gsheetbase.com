import { Card, Progress, Alert, Spin, Typography, Tag, Row, Col, Tooltip } from 'antd'
import { ThunderboltOutlined } from '@ant-design/icons'
import { useSheetAnalytics } from '../../hooks/useSheetAnalytics'
import { useSubscriptionPlan } from '../../hooks/useSubscriptionPlan'
import styled from 'styled-components'
import dayjs from 'dayjs'

const { Text, Title } = Typography

const SectionTitle = styled(Text)`
  display: block;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.45);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 16px;
`

const QuotaRow = styled.div`
  margin-bottom: 24px;
`

const QuotaLabel = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
`

const PlanHeader = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
  padding-bottom: 20px;
  border-bottom: 1px solid #f0f0f0;
`

const PLAN_COLORS: Record<string, string> = {
    free: 'default',
    starter: 'blue',
    pro: 'purple',
    enterprise: 'gold',
}

function getProgressColor(percent: number): string {
    if (percent >= 95) return '#ff4d4f'
    if (percent >= 80) return '#faad14'
    return '#52c41a'
}

function formatQuota(value: number): string {
    if (value === -1 || value === 0) return '∞'
    if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`
    if (value >= 1_000) return `${(value / 1_000).toFixed(0)}K`
    return `${value}`
}

interface QuotaBarProps {
    label: string
    used: number
    limit: number
    tooltip?: string
}

function QuotaBar({ label, used, limit, tooltip }: QuotaBarProps) {
    const isUnlimited = limit <= 0
    const percent = isUnlimited ? 0 : Math.min(100, Math.round((used / limit) * 100))
    const color = getProgressColor(percent)

    return (
        <QuotaRow>
            <QuotaLabel>
                <Tooltip title={tooltip}>
                    <Text strong>{label}</Text>
                </Tooltip>
                <Text type="secondary">
                    {used.toLocaleString()} / {isUnlimited ? 'Unlimited' : limit.toLocaleString()}
                </Text>
            </QuotaLabel>
            {isUnlimited ? (
                <Progress percent={100} strokeColor="#52c41a" format={() => 'Unlimited'} />
            ) : (
                <Progress
                    percent={percent}
                    strokeColor={color}
                    format={(p) => `${p}%`}
                    status={percent >= 95 ? 'exception' : 'normal'}
                />
            )}
        </QuotaRow>
    )
}

interface UsageTabProps {
    sheetId: string
}

export function UsageTab({ sheetId }: UsageTabProps) {
    const { data: analytics, isLoading: analyticsLoading, error: analyticsError } = useSheetAnalytics(sheetId, 30)
    const { data: plan, isLoading: planLoading, error: planError } = useSubscriptionPlan()

    if (analyticsLoading || planLoading) {
        return (
            <Card>
                <div style={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
                    <Spin size="large" />
                </div>
            </Card>
        )
    }

    if (analyticsError || planError) {
        return (
            <Card>
                <Alert
                    message="Failed to load usage data"
                    description="Unable to fetch usage or plan data. Please try again later."
                    type="error"
                    showIcon
                />
            </Card>
        )
    }

    // Derive per-sheet metrics from 30-day analytics
    const today = dayjs().format('YYYY-MM-DD')
    const todayEntry = analytics?.daily_usage.find((d) => d.date === today)
    const todayUpdates = todayEntry
        ? (todayEntry.post_count || 0) + (todayEntry.put_count || 0) + (todayEntry.patch_count || 0)
        : 0

    const monthlyGets = (analytics?.daily_usage || []).reduce((sum, d) => sum + (d.get_count || 0), 0)
    const monthlyUpdates = (analytics?.daily_usage || []).reduce(
        (sum, d) => sum + (d.post_count || 0) + (d.put_count || 0) + (d.patch_count || 0),
        0
    )

    const planName = plan?.plan ?? 'free'
    const planTagColor = PLAN_COLORS[planName] ?? 'default'

    const showUpgradeAlert =
        plan &&
        plan.daily_update_quota > 0 &&
        plan.monthly_get_quota > 0 &&
        plan.monthly_update_quota > 0 &&
        (todayUpdates / plan.daily_update_quota >= 0.8 ||
            monthlyGets / plan.monthly_get_quota >= 0.8 ||
            monthlyUpdates / plan.monthly_update_quota >= 0.8)

    return (
        <Card>
            <PlanHeader>
                <ThunderboltOutlined style={{ fontSize: 20, color: '#1677ff' }} />
                <div>
                    <Title level={5} style={{ margin: 0 }}>
                        Usage vs Plan Quota
                    </Title>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        Usage for this sheet over the last 30 days
                    </Text>
                </div>
                <Tag color={planTagColor} style={{ marginLeft: 'auto', textTransform: 'capitalize' }}>
                    {planName}
                </Tag>
            </PlanHeader>

            {showUpgradeAlert && (
                <Alert
                    message="Approaching quota limit"
                    description={
                        <span>
                            You're using 80%+ of one or more quotas on your <strong>{planName}</strong> plan.{' '}
                            <a href="/billing">Upgrade your plan</a> to avoid disruptions.
                        </span>
                    }
                    type="warning"
                    showIcon
                    closable
                    style={{ marginBottom: 24 }}
                />
            )}

            <Row gutter={[32, 0]}>
                <Col xs={24} md={12}>
                    <SectionTitle>Daily</SectionTitle>
                    <QuotaBar
                        label="Today's UPDATE requests"
                        used={todayUpdates}
                        limit={plan?.daily_update_quota ?? 0}
                        tooltip="POST, PUT, PATCH requests made to this sheet's API endpoint today"
                    />
                </Col>
                <Col xs={24} md={12}>
                    <SectionTitle>Monthly (last 30 days)</SectionTitle>
                    <QuotaBar
                        label="GET requests"
                        used={monthlyGets}
                        limit={plan?.monthly_get_quota ?? 0}
                        tooltip="Read requests to this sheet's API endpoint over the last 30 days"
                    />
                    <QuotaBar
                        label="UPDATE requests"
                        used={monthlyUpdates}
                        limit={plan?.monthly_update_quota ?? 0}
                        tooltip="Write requests (POST, PUT, PATCH) to this sheet's API endpoint over the last 30 days"
                    />
                </Col>
            </Row>

            <div style={{ marginTop: 8, paddingTop: 16, borderTop: '1px solid #f0f0f0' }}>
                <Text type="secondary" style={{ fontSize: 12 }}>
                    Rate limits: {formatQuota(plan?.get_rate_limit ?? 0)} GET req/min · {formatQuota(plan?.update_rate_limit ?? 0)} UPDATE req/min
                </Text>
            </div>
        </Card>
    )
}
