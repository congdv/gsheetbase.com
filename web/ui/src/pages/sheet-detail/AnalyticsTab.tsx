import { Card, Alert, Spin, Empty, Select, Row, Col, Statistic } from 'antd'
import { Line } from '@ant-design/charts'
import { useParams } from 'react-router-dom'
import { useState } from 'react'
import { useSheetAnalytics } from '../../hooks/useSheetAnalytics'
import styled from 'styled-components'

const ChartContainer = styled.div`
  margin-top: 24px;
`

const ControlsContainer = styled.div`
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
`

const StatsContainer = styled.div`
  margin-bottom: 24px;
`

export function AnalyticsTab() {
  const { id } = useParams<{ id: string }>()
  const [days, setDays] = useState<number>(7)
  const { data, isLoading, error } = useSheetAnalytics(id!, days)

  if (isLoading) {
    return (
      <Card>
        <div style={{ display: 'flex', justifyContent: 'center', padding: '40px' }}>
          <Spin size="large" />
        </div>
      </Card>
    )
  }

  if (error) {
    return (
      <Card>
        <Alert
          message="Failed to load analytics"
          description="Unable to fetch analytics data. Please try again later."
          type="error"
          showIcon
        />
      </Card>
    )
  }

  if (!data || data.daily_usage.length === 0) {
    return (
      <Card>
        <Empty
          description="No analytics data available yet"
          style={{ padding: '40px' }}
        />
      </Card>
    )
  }

  // Build a contiguous series for the selected `days` range (fill zeros for missing dates)
  const formatDate = (d: Date) => {
    const yyyy = d.getFullYear()
    const mm = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    return `${yyyy}-${mm}-${dd}`
  }

  const end = new Date()
  end.setHours(0, 0, 0, 0)
  const start = new Date(end)
  start.setDate(end.getDate() - (days - 1))

  const dateMap = new Map<string, any>(data.daily_usage.map((d: any) => [d.date, d]))

  const filledDailyUsage = Array.from({ length: days }).map((_, i) => {
    const dt = new Date(start)
    dt.setDate(start.getDate() + i)
    const dateStr = formatDate(dt)
    const src = dateMap.get(dateStr)
    return src
      ? src
      : {
          date: dateStr,
          total_count: 0,
          get_count: 0,
          post_count: 0,
          put_count: 0,
          patch_count: 0,
          delete_count: 0,
        }
  })

  // Calculate total/avg/peak from the filled series
  const totalRequests = filledDailyUsage.reduce((sum: number, day: any) => sum + (day.total_count || 0), 0)
  const avgDailyRequests = Math.round(totalRequests / days)

  const peakDay = filledDailyUsage.reduce((max: any, day: any) => (day.total_count > (max.total_count || 0) ? day : max), filledDailyUsage[0])

  // Prepare chart data from filled series
  const chartData = filledDailyUsage.map((day: any) => ({
    date: day.date,
    total_count: day.total_count || 0,
    GET: day.get_count || 0,
    POST: day.post_count || 0,
    PUT: day.put_count || 0,
    PATCH: day.patch_count || 0,
    DELETE: day.delete_count || 0,
  }))


  const config = {
    data: chartData,
    xField: 'date',
    yField: 'total_count',
    smooth: true,
    tooltip: {
      items: [
        (datum: any) => {

          return {
            name: 'Total',
            value: datum.total_count,
          }
        },
        (datum: any) => {
          return {
            name: 'GET',
            value: datum.GET,
          }
        },
        (datum: any) => {
          return {
            name: 'POST',
            value: datum.POST,
          }
        },
        (datum: any) => {
          return {
            name: 'PATCH',
            value: datum.PATCH,
          }
        },
        (datum: any) => {
          return {
            name: 'PUT',
            value: datum.PUT,
          }
        },
        (datum: any) => {
          return {
            name: 'DELETE',
            value: datum.DELETE,
          }
        },
      ],
    },
    xAxis: {
      label: {
        autoRotate: true,
        autoHide: true,
      },
    },
    yAxis: {
      label: {
        formatter: (v: string) => {
          const num = parseInt(v)
          if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`
          if (num >= 1000) return `${(num / 1000).toFixed(1)}K`
          return v
        },
      },
    },
  }

  return (
    <Card>
      <ControlsContainer>
        <h3 style={{ margin: 0 }}>API Usage Over Time</h3>
        <Select
          value={days}
          onChange={(value: number) => setDays(value)}
          options={[
            { label: 'Last 7 Days', value: 7 },
            { label: 'Last 30 Days', value: 30 },
            { label: 'Last 90 Days', value: 90 },
          ]}
          style={{ width: 160 }}
        />
      </ControlsContainer>

      <StatsContainer>
        <Row gutter={16}>
          <Col xs={24} sm={8}>
            <Statistic
              title="Total Requests"
              value={totalRequests}
              precision={0}
            />
          </Col>
          <Col xs={24} sm={8}>
            <Statistic
              title="Avg Daily Requests"
              value={avgDailyRequests}
              precision={0}
            />
          </Col>
          <Col xs={24} sm={8}>
            <Statistic
              title="Peak Day"
              value={peakDay.total_count}
              precision={0}
              suffix={`(${peakDay.date})`}
            />
          </Col>
        </Row>
      </StatsContainer>

      <ChartContainer>
        <Line {...config} />
      </ChartContainer>
    </Card>
  )
}
