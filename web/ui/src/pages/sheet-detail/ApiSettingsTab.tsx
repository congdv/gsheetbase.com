import { Card, Space, Typography, Input, Tooltip, Descriptions, Result, Button } from 'antd'
import { CopyOutlined, RocketOutlined } from '@ant-design/icons'

const { Paragraph, Text } = Typography

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

interface ApiSettingsTabProps {
  sheet: Sheet
  onCopy: (text: string) => void
  onPublish: () => void
}

export function ApiSettingsTab({ sheet, onCopy, onPublish }: ApiSettingsTabProps) {
  if (!sheet.is_public || !sheet.api_key) {
    return (
      <Card>
        <Result
          icon={<RocketOutlined />}
          title="Sheet Not Published"
          subTitle="This sheet is not published yet. Publish it to generate an API endpoint."
          extra={
            <Button type="primary" onClick={onPublish}>
              Publish Sheet
            </Button>
          }
        />
      </Card>
    )
  }

  return (
    <Card title="API Configuration">
      <Space direction="vertical" size="middle" style={{ width: '100%' }}>
        <div>
          <Text strong>API Endpoint</Text>
          <Input
            value={`https://api.gsheetbase.com/v1/${sheet.api_key}`}
            readOnly
            addonAfter={
              <Tooltip title="Copy">
                <CopyOutlined
                  onClick={() =>
                    onCopy(`https://api.gsheetbase.com/v1/${sheet.api_key}`)
                  }
                  style={{ cursor: 'pointer' }}
                />
              </Tooltip>
            }
            style={{ marginTop: 8 }}
          />
          <Paragraph style={{ marginTop: 8, fontSize: 12, color: '#666' }}>
            Use this endpoint to access your sheet data. Add <code>?range=A1:Z100</code> to customize the range.
          </Paragraph>
        </div>

        <Descriptions bordered column={1} size="small">
          <Descriptions.Item label="Default Range">
            {sheet.default_range || 'Sheet1'}
          </Descriptions.Item>
          <Descriptions.Item label="Transform to JSON Objects">
            {sheet.use_first_row_as_header ? 'Yes' : 'No'}
          </Descriptions.Item>
        </Descriptions>

        <div>
          <Text strong>Example Usage</Text>
          <pre
            style={{
              background: '#f5f5f5',
              padding: 12,
              borderRadius: 4,
              overflow: 'auto',
              marginTop: 8,
            }}
          >
            {`// JavaScript
fetch('https://api.gsheetbase.com/v1/${sheet.api_key}')
  .then(res => res.json())
  .then(data => console.log(data))

// cURL
curl https://api.gsheetbase.com/v1/${sheet.api_key}`}
          </pre>
        </div>
      </Space>
    </Card>
  )
}
