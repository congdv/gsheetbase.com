import { Card, Space, Typography, Input, Tooltip, Descriptions, Result, Button, Select, Spin, Alert, Tag, Row, Col, Switch, message } from 'antd'
import { CopyOutlined, RocketOutlined, LockOutlined } from '@ant-design/icons'
import { useState } from 'react'
import api from '../../lib/axios'
import { queryClient } from '../../lib/queryClient'
import { useConfig } from '../../context/ConfigContext'
import { useAuth } from '../../context/AuthContext'
import { ScopeConsentPrompt, ScopeInfo } from '../../components/ScopeConsentPrompt'
import { AuthManagementCard } from '../../components/sheets/AuthManagementCard'
import { ApiTesterCard } from '../../components/sheets/ApiTesterCard'
import { GOOGLE_SCOPE } from '@/constants/common'

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
  auth_type?: string
  auth_bearer_token?: string
  auth_basic_username?: string
  created_at: string
}

interface ApiSettingsTabProps {
  sheet: Sheet
  onCopy: (text: string) => void
  onPublish: () => void
}

export function ApiSettingsTab({ sheet, onCopy, onPublish }: ApiSettingsTabProps) {
  const config = useConfig()
  const apiUrl = `${config.workerBaseUrl}/v1/${sheet.api_key}`
  const { hasScope, requestScopes } = useAuth()

  const [showScopePrompt, setShowScopePrompt] = useState(false)
  const [pendingMethod, setPendingMethod] = useState<{ method: string; checked: boolean } | null>(null)
  const [updatingMethods, setUpdatingMethods] = useState<Record<string, boolean>>({})
  const [response, setResponse] = useState<{
    status: number
    data: any
    headers?: Record<string, string>
    error?: string
  } | null>(null)

  const canReadAndWrite = hasScope(GOOGLE_SCOPE.READ_WRITE_SCOPE)
  const canReadOnly = hasScope(GOOGLE_SCOPE.READ_ONLY_SCOPE)

  const handleRequestWriteScope = async (selectedScopes?: string[]) => {
    try {
      const scopes = selectedScopes ?? [GOOGLE_SCOPE.READ_WRITE_SCOPE]
      await requestScopes(scopes)
      setShowScopePrompt(false)
      if (pendingMethod) {
        const { method, checked } = pendingMethod
        setPendingMethod(null)
        // Proceed with the pending toggle after consent is granted
        handleToggleMethod(method, checked)
      }
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
      // Invalidate sheets query so parent page refetches updated sheet data
      await queryClient.invalidateQueries({ queryKey: ['sheets'] })
    } catch (error: any) {
      message.error(`Failed to update method settings: ${error.response?.data?.error || error.message}`)
    } finally {
      setUpdatingMethods(prev => ({ ...prev, [method]: false }))
    }
  }

  const scopeInfo: ScopeInfo[] = [
    {
      scope: GOOGLE_SCOPE.READ_WRITE_SCOPE,
      reason: 'Write access to your Google Sheets',
      example: 'Allows adding, updating, and deleting rows via API',
      optional: false,
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
            Use this endpoint to access your sheet data.
          </Paragraph>
        </div>

        <Descriptions bordered column={1} size="small">
          <Descriptions.Item label="Default Collection">
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
                {canReadOnly ? (
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
                  {!canReadAndWrite && (
                    <Tooltip title="Requires write permission">
                      <LockOutlined style={{ color: '#999' }} />
                    </Tooltip>
                  )}
                  <Switch
                    checked={isMethodEnabled('POST')}
                    onChange={(checked) => {
                      if (checked && !canReadAndWrite) {
                        setPendingMethod({ method: 'POST', checked })
                        setShowScopePrompt(true)
                      } else {
                        handleToggleMethod('POST', checked)
                      }
                    }}
                    loading={updatingMethods['POST']}
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
                  {!canReadAndWrite && (
                    <Tooltip title="Requires write permission">
                      <LockOutlined style={{ color: '#999' }} />
                    </Tooltip>
                  )}
                  <Switch
                    checked={isMethodEnabled('PUT')}
                    onChange={(checked) => handleToggleMethod('PUT', checked)}
                    loading={updatingMethods['PUT']}
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
                  {!canReadAndWrite && (
                    <Tooltip title="Requires write permission">
                      <LockOutlined style={{ color: '#999' }} />
                    </Tooltip>
                  )}
                  <Switch
                    checked={isMethodEnabled('PATCH')}
                    onChange={(checked) => handleToggleMethod('PATCH', checked)}
                    loading={updatingMethods['PATCH']}
                  />
                </Space>
              </Col>
            </Row>

            <Row justify="space-between" align="middle">
              <Col>
                <Text><code>DELETE</code> Delete rows</Text>
              </Col>
              <Col>
                <Space>
                  {!canReadAndWrite && (
                    <Tooltip title="Requires write permission">
                      <LockOutlined style={{ color: '#999' }} />
                    </Tooltip>
                  )}
                  <Switch
                    checked={isMethodEnabled('DELETE')}
                    onChange={(checked) => handleToggleMethod('DELETE', checked)}
                    loading={updatingMethods['DELETE']}
                  />
                </Space>
              </Col>
            </Row>
          </Space>
        </Card>

        <AuthManagementCard sheetId={sheet.id} currentAuthType={sheet.auth_type || 'none'} sheet={sheet} />

        <ApiTesterCard
          apiKey={sheet.api_key!}
          isMethodEnabled={isMethodEnabled}
          authType={sheet.auth_type ?? 'none'}
          basicUsername={sheet.auth_basic_username}
        />

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
