import { useParams, useNavigate } from 'react-router-dom'
import { Typography, Button, Spin, Result, Tabs } from 'antd'
import { ArrowLeftOutlined, InfoCircleOutlined, ApiOutlined, TableOutlined, BarChartOutlined } from '@ant-design/icons'
import { useQuery } from '@tanstack/react-query'
import { useState, useEffect } from 'react'
import styled from 'styled-components'
import api from '../../lib/axios'
import { ROUTES } from '../../constants/routes'
import { message } from 'antd'
import { OverviewTab } from './OverviewTab'
import { ApiSettingsTab } from './ApiSettingsTab'
import { DataPreviewTab } from './DataPreviewTab'
import { AnalyticsTab } from './AnalyticsTab'

const { Title, Paragraph } = Typography

const PageContainer = styled.div`
  padding: 24px;
  margin: 0 auto;
  max-width: 1200px;
`

const PageHeader = styled.div`
  margin-bottom: 24px;
`

const BackButton = styled(Button)`
  margin-bottom: 16px;
`

interface Sheet {
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

export default function SheetDetailPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [activeTab, setActiveTab] = useState<string>('overview')

  // Sync tab with URL hash
  useEffect(() => {
    const hash = window.location.hash.slice(1) // Remove '#'
    if (hash && ['overview', 'api-settings', 'data-preview', 'analytics'].includes(hash)) {
      setActiveTab(hash)
    }
  }, [])

  const handleTabChange = (key: string) => {
    setActiveTab(key)
    window.location.hash = key
  }

  // For now, fetch all sheets and filter by ID
  // TODO: Create dedicated backend endpoint GET /sheets/:id
  const { data: sheets, isLoading, error } = useQuery({
    queryKey: ['sheets'],
    queryFn: async () => {
      const response = await api.get('/sheets/registered')
      return (response.data.sheets || []) as Sheet[]
    },
  })

  const sheet = sheets?.find((s) => s.id === id)

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    message.success('Copied to clipboard!')
  }

  const handlePublish = () => {
    // Navigate to home where publish modal can be opened
    navigate(ROUTES.HOME)
    message.info('Please use the Publish action from the sheets table')
  }

  if (isLoading) {
    return (
      <PageContainer>
        <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '400px' }}>
          <Spin size="large" />
        </div>
      </PageContainer>
    )
  }

  if (error || !sheet) {
    return (
      <PageContainer>
        <Result
          status="404"
          title="Sheet Not Found"
          subTitle="The sheet you're looking for doesn't exist or has been removed."
          extra={
            <Button type="primary" onClick={() => navigate(ROUTES.HOME)}>
              Back to Home
            </Button>
          }
        />
      </PageContainer>
    )
  }

  return (
    <PageContainer>
      <PageHeader>
        <BackButton
          icon={<ArrowLeftOutlined />}
          onClick={() => navigate(ROUTES.HOME)}
        >
          Back to My Sheets
        </BackButton>
        <Title level={2}>{sheet.sheet_name || sheet.sheet_id}</Title>
        {sheet.description && <Paragraph>{sheet.description}</Paragraph>}
      </PageHeader>

      <Tabs
        activeKey={activeTab}
        onChange={handleTabChange}
        items={[
          {
            key: 'overview',
            label: (
              <span>
                <InfoCircleOutlined />
                Overview
              </span>
            ),
            children: <OverviewTab sheet={sheet} onCopy={copyToClipboard} />,
          },
          {
            key: 'api-settings',
            label: (
              <span>
                <ApiOutlined />
                API Settings
              </span>
            ),
            children: <ApiSettingsTab sheet={sheet} onCopy={copyToClipboard} onPublish={handlePublish} />,
          },
          {
            key: 'data-preview',
            label: (
              <span>
                <TableOutlined />
                Data Preview
              </span>
            ),
            children: <DataPreviewTab />,
          },
          {
            key: 'analytics',
            label: (
              <span>
                <BarChartOutlined />
                Analytics
              </span>
            ),
            children: <AnalyticsTab />,
          },
        ]}
      />
    </PageContainer>
  )
}
