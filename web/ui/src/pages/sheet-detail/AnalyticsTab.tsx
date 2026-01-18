import { Card, Result } from 'antd'
import { BarChartOutlined } from '@ant-design/icons'

export function AnalyticsTab() {
  return (
    <Card>
      <Result
        icon={<BarChartOutlined style={{ color: '#52c41a' }} />}
        title="Analytics Coming Soon"
        subTitle="Track API usage, request counts, response times, and geographic distribution of your API calls."
      />
    </Card>
  )
}
