import { Card, Descriptions, Tag, Button, Space } from 'antd'
import { LinkOutlined, CopyOutlined } from '@ant-design/icons'

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

interface OverviewTabProps {
  sheet: Sheet
  onCopy: (text: string) => void
  onNavigateToApiSettings?: () => void
}

export function OverviewTab({ sheet, onCopy, onNavigateToApiSettings }: OverviewTabProps) {
  const workerBaseUrl = import.meta.env.VITE_WORKER_BASE_URL || 'https://api.gsheetbase.com'
  const apiUrl = sheet.api_key ? `${workerBaseUrl}/v1/${sheet.api_key}` : null

  return (
    <Card title="Sheet Information">
      <Descriptions bordered column={1}>
        <Descriptions.Item label="Sheet Name">
          {sheet.sheet_name || '—'}
        </Descriptions.Item>
        <Descriptions.Item label="Description">
          {sheet.description || '—'}
        </Descriptions.Item>
        <Descriptions.Item label="Sheet ID">
          <Space>
            <span>{sheet.sheet_id}</span>
            <Button
              size="small"
              icon={<CopyOutlined />}
              onClick={() => onCopy(sheet.sheet_id)}
            >
              Copy
            </Button>
          </Space>
        </Descriptions.Item>
        <Descriptions.Item label="API URL">
          <Space>
            <span
              style={{
                fontFamily: 'monospace',
                fontSize: '12px',
                cursor: apiUrl && onNavigateToApiSettings ? 'pointer' : 'default',
                color: apiUrl && onNavigateToApiSettings ? '#1890ff' : 'inherit',
                textDecoration: apiUrl && onNavigateToApiSettings ? 'underline' : 'none',
              }}
              onClick={() => apiUrl && onNavigateToApiSettings?.()}
            >
              {apiUrl || '—'}
            </span>
            {apiUrl && (
              <Button
                size="small"
                icon={<CopyOutlined />}
                onClick={() => onCopy(apiUrl)}
              >
                Copy
              </Button>
            )}
          </Space>
        </Descriptions.Item>
        <Descriptions.Item label="Status">
          {sheet.is_public ? (
            <Tag color="green">Public</Tag>
          ) : (
            <Tag>Private</Tag>
          )}
        </Descriptions.Item>
        <Descriptions.Item label="Registered">
          {new Date(sheet.created_at).toLocaleString()}
        </Descriptions.Item>
        <Descriptions.Item label="Google Sheets">
          <Button
            type="link"
            icon={<LinkOutlined />}
            href={`https://docs.google.com/spreadsheets/d/${sheet.sheet_id}`}
            target="_blank"
            style={{ padding: 0 }}
          >
            Open in Google Sheets
          </Button>
        </Descriptions.Item>
      </Descriptions>
    </Card>
  )
}
