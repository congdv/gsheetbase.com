import { Card, Space, Typography, Input, Tooltip, Descriptions, Result, Button, Select, Spin, Alert, Tag, Row, Col, Switch, message } from 'antd'
import { CopyOutlined, RocketOutlined, SendOutlined, LockOutlined } from '@ant-design/icons'
import { useState } from 'react'
import axios from 'axios'
import api from '../../lib/axios'
import { useAuth } from '../../context/AuthContext'
import { ScopeConsentPrompt, ScopeInfo } from '../../components/ScopeConsentPrompt'

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
  allow_write?: boolean
  allowed_methods?: string[]
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
  const { hasScope, requestScopes } = useAuth()

  const [testUrl, setTestUrl] = useState(apiUrl)
  const [httpMethod, setHttpMethod] = useState<'GET' | 'POST' | 'PUT' | 'PATCH'>('GET')
  const [requestBody, setRequestBody] = useState('{\n  \n}')
  const [loading, setLoading] = useState(false)
  const [showScopePrompt, setShowScopePrompt] = useState(false)
  const [updatingMethods, setUpdatingMethods] = useState<Record<string, boolean>>({})
  const [response, setResponse] = useState<{
    status: number
    data: any
    headers?: Record<string, string>
    error?: string
  } | null>(null)

  const canRead = hasScope('https://www.googleapis.com/auth/spreadsheets.readonly')
  const canWrite = hasScope('https://www.googleapis.com/auth/spreadsheets')

  const handleSendRequest = async () => {
    setLoading(true)
    setResponse(null)

    try {
      const startTime = Date.now()
      const config: any = {
        method: httpMethod,
        url: testUrl,
        validateStatus: () => true,
      }
      
      // Add request body for write methods
      if (['POST', 'PUT', 'PATCH'].includes(httpMethod)) {
        try {
          config.data = JSON.parse(requestBody)
        } catch (e) {
          setResponse({
            status: 0,
            data: null,
            error: 'Invalid JSON in request body',
          })
          setLoading(false)
          return
        }
      }
      
      const res = await axios(config)
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

  const handleRequestWriteScope = async () => {
    try {
      await requestScopes(['https://www.googleapis.com/auth/spreadsheets'])
      setShowScopePrompt(false)
    } catch (error) {
      console.error('Failed to request write scope:', error)
    }
  }

  const isMethodEnabled = (method: string): boolean => {
    return sheet.allowed_methods?.includes(method) || false
  }

  const handleToggleMethod = async (method: string, checked: boolean) => {
    setUpdatingMethods(prev => ({ ...prev, [method]: true }))
    try {
      const currentMethods = sheet.allowed_methods || []
      const updatedMethods = checked
        ? [...currentMethods, method]
        : currentMethods.filter(m => m !== method)
      
      await api.patch(
        `/sheets/${sheet.id}/write-settings`,
        { allowed_methods: updatedMethods }
      )
      message.success(`${method} ${checked ? 'enabled' : 'disabled'} successfully`)
      // Trigger a refresh by updating parent component (if needed)
      window.location.reload() // Simple refresh for now
    } catch (error: any) {
      message.error(`Failed to update method settings: ${error.response?.data?.error || error.message}`)
    } finally {
      setUpdatingMethods(prev => ({ ...prev, [method]: false }))
    }
  }

  const scopeInfo: ScopeInfo[] = [
    {
      scope: 'spreadsheets',
      reason: 'Write access to your Google Sheets',
      example: 'Allows adding, updating, and deleting rows via API',
    },
  ]

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

        <Card title="API Methods & Permissions" size="small">
          <Space direction="vertical" style={{ width: '100%' }}>
            <Row justify="space-between" align="middle">
              <Col>
                <Text><code>GET</code> Fetch sheet data</Text>
              </Col>
              <Col>
                {canRead ? (
                  <Tag color="success">Always Available</Tag>
                ) : (
                  <Tooltip title="Grant read permission to enable">
                    <Tag>Locked</Tag>
                  </Tooltip>
                )}
              </Col>
            </Row>
            
            <Row justify="space-between" align="middle">
              <Col>
                <Text><code>POST</code> Add new rows</Text>
              </Col>
              <Col>
                <Space>
                  {!canWrite && (
                    <Tooltip title="Requires write permission">
                      <LockOutlined style={{ color: '#999' }} />
                    </Tooltip>
                  )}
                  <Switch 
                    checked={isMethodEnabled('POST')} 
                    onChange={(checked) => handleToggleMethod('POST', checked)} 
                    loading={updatingMethods['POST']}
                    disabled={!canWrite}
                  />
                </Space>
              </Col>
            </Row>

            <Row justify="space-between" align="middle">
              <Col>
                <Text><code>PUT</code> Update rows</Text>
              </Col>
              <Col>
                <Space>
                  {!canWrite && (
                    <Tooltip title="Requires write permission">
                      <LockOutlined style={{ color: '#999' }} />
                    </Tooltip>
                  )}
                  <Switch 
                    checked={isMethodEnabled('PUT')} 
                    onChange={(checked) => handleToggleMethod('PUT', checked)} 
                    loading={updatingMethods['PUT']}
                    disabled={!canWrite}
                  />
                </Space>
              </Col>
            </Row>

            <Row justify="space-between" align="middle">
              <Col>
                <Text><code>PATCH</code> Partially update rows</Text>
              </Col>
              <Col>
                <Space>
                  {!canWrite && (
                    <Tooltip title="Requires write permission">
                      <LockOutlined style={{ color: '#999' }} />
                    </Tooltip>
                  )}
                  <Switch 
                    checked={isMethodEnabled('PATCH')} 
                    onChange={(checked) => handleToggleMethod('PATCH', checked)} 
                    loading={updatingMethods['PATCH']}
                    disabled={!canWrite}
                  />
                </Space>
              </Col>
            </Row>
          </Space>
        </Card>

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
                    { label: 'POST', value: 'POST', disabled: !isMethodEnabled('POST') },
                    { label: 'PUT', value: 'PUT', disabled: !isMethodEnabled('PUT') },
                    { label: 'PATCH', value: 'PATCH', disabled: !isMethodEnabled('PATCH') },
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

            {['POST', 'PUT', 'PATCH'].includes(httpMethod) && (
              <div>
                <Space direction="horizontal" style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                  <Text strong>Request Body</Text>
                  <Button 
                    size="small" 
                    onClick={() => {
                      const example = httpMethod === 'POST' 
                        ? '{\n  "data": [\n    ["John", "Doe", "30"],\n    ["Jane", "Smith", "25"]\n  ]\n}'
                        : httpMethod === 'PUT'
                        ? '{\n  "data": [\n    ["Updated", "Row", "1"],\n    ["Updated", "Row", "2"]\n  ],\n  "range": "Sheet1!A2:C3"\n}'
                        : '{\n  "data": [\n    ["Patched", "Value"]\n  ],\n  "range": "Sheet1!A2:B2"\n}'
                      setRequestBody(example)
                    }}
                  >
                    Load Example
                  </Button>
                </Space>
                <Input.TextArea
                  value={requestBody}
                  onChange={(e) => setRequestBody(e.target.value)}
                  placeholder="JSON request body"
                  rows={6}
                  style={{ fontFamily: 'monospace' }}
                />
                <Text type="secondary" style={{ fontSize: 12, display: 'block', marginTop: 4 }}>
                  {httpMethod === 'POST' && 'POST appends new rows to the sheet. Range is optional.'}
                  {httpMethod === 'PUT' && 'PUT updates existing rows. Range is required.'}
                  {httpMethod === 'PATCH' && 'PATCH partially updates rows. Range is required.'}
                </Text>
              </div>
            )}

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
            {httpMethod === 'GET' ? `// JavaScript - Fetch data
fetch('${apiUrl}')
  .then(res => res.json())
  .then(data => console.log(data))

// cURL - Fetch data
curl ${apiUrl}` : 
httpMethod === 'POST' ? `// JavaScript - Add new rows
fetch('${apiUrl}', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    data: [
      ["John", "Doe", "30"],
      ["Jane", "Smith", "25"]
    ]
  })
})

// cURL - Add new rows
curl -X POST ${apiUrl} \\
  -H 'Content-Type: application/json' \\
  -d '{"data": [["John", "Doe", "30"]]}'` :
httpMethod === 'PUT' ? `// JavaScript - Update rows
fetch('${apiUrl}', {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    data: [["Updated", "Row"]],
    range: "Sheet1!A2:B2"
  })
})

// cURL - Update rows
curl -X PUT ${apiUrl} \\
  -H 'Content-Type: application/json' \\
  -d '{"data": [["Updated", "Row"]], "range": "Sheet1!A2:B2"}'` :
`// JavaScript - Partial update
fetch('${apiUrl}', {
  method: 'PATCH',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    data: [["Patched"]],
    range: "Sheet1!A2"
  })
})

// cURL - Partial update
curl -X PATCH ${apiUrl} \\
  -H 'Content-Type: application/json' \\
  -d '{"data": [["Patched"]], "range": "Sheet1!A2"}'`}
          </pre>
        </div>

        <ScopeConsentPrompt
          open={showScopePrompt}
          onConsent={handleRequestWriteScope}
          onCancel={() => setShowScopePrompt(false)}
          scopes={scopeInfo}
        />
      </Space>
    </Card>
  )
}
