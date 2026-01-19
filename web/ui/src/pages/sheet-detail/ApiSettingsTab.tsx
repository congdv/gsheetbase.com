import { Card, Space, Typography, Input, Tooltip, Descriptions, Result, Button, Select, Spin, Alert } from 'antd'
import { CopyOutlined, RocketOutlined, SendOutlined } from '@ant-design/icons'
import { useState } from 'react'
import axios from 'axios'

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
  const workerBaseUrl = import.meta.env.VITE_WORKER_BASE_URL || 'https://api.gsheetbase.com'
  const apiUrl = `${workerBaseUrl}/v1/${sheet.api_key}`

  const [testUrl, setTestUrl] = useState(apiUrl)
  const [httpMethod, setHttpMethod] = useState<'GET' | 'POST'>('GET')
  const [loading, setLoading] = useState(false)
  const [response, setResponse] = useState<{
    status: number
    data: any
    headers?: Record<string, string>
    error?: string
  } | null>(null)

  const handleSendRequest = async () => {
    setLoading(true)
    setResponse(null)

    try {
      const startTime = Date.now()
      const res = await axios({
        method: httpMethod,
        url: testUrl,
        validateStatus: () => true, // Don't throw on any status
      })
      const endTime = Date.now()

      setResponse({
        status: res.status,
        data: res.data,
        headers: {
          'content-type': res.headers['content-type'] || '',
          'x-response-time': `${endTime - startTime}ms`,
        },
      })
    } catch (error: any) {
      setResponse({
        status: 0,
        data: null,
        error: error.message || 'Request failed',
      })
    } finally {
      setLoading(false)
    }
  }

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
            value={apiUrl}
            readOnly
            addonAfter={
              <Tooltip title="Copy">
                <CopyOutlined
                  onClick={() => onCopy(apiUrl)}
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

        <Card title="API Tester" size="small">
          <Space direction="vertical" size="middle" style={{ width: '100%' }}>
            <div>
              <Text strong style={{ display: 'block', marginBottom: 8 }}>Test Your API</Text>
              <Space.Compact style={{ width: '100%' }}>
                <Select
                  value={httpMethod}
                  onChange={setHttpMethod}
                  style={{ width: 100 }}
                  options={[
                    { label: 'GET', value: 'GET' },
                    { label: 'POST', value: 'POST' },
                  ]}
                />
                <Input
                  value={testUrl}
                  onChange={(e) => setTestUrl(e.target.value)}
                  placeholder="API URL"
                />
                <Button
                  type="primary"
                  icon={<SendOutlined />}
                  onClick={handleSendRequest}
                  loading={loading}
                >
                  Send
                </Button>
              </Space.Compact>
            </div>

            {response && (
              <div>
                <Space style={{ marginBottom: 8 }}>
                  <Text strong>Response:</Text>
                  <Text
                    type={response.status >= 200 && response.status < 300 ? 'success' : 'danger'}
                  >
                    Status {response.status}
                  </Text>
                  {response.headers?.['x-response-time'] && (
                    <Text type="secondary">{response.headers['x-response-time']}</Text>
                  )}
                </Space>

                {response.error ? (
                  <Alert
                    message="Request Failed"
                    description={response.error}
                    type="error"
                    showIcon
                  />
                ) : (
                  <pre
                    style={{
                      background: '#f5f5f5',
                      padding: 12,
                      borderRadius: 4,
                      overflow: 'auto',
                      maxHeight: 400,
                      border: '1px solid #d9d9d9',
                    }}
                  >
                    {JSON.stringify(response.data, null, 2)}
                  </pre>
                )}
              </div>
            )}
          </Space>
        </Card>

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
fetch('${apiUrl}')
  .then(res => res.json())
  .then(data => console.log(data))

// cURL
curl ${apiUrl}`}
          </pre>
        </div>
      </Space>
    </Card>
  )
}
