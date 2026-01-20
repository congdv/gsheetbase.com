import { Card, Descriptions, Tag, Button, Space, Modal, Form, Input, Switch, Popconfirm, List, Divider, Typography, Tooltip } from 'antd'
import { LinkOutlined, CopyOutlined, CheckCircleOutlined, LockOutlined } from '@ant-design/icons'
import { useState } from 'react'
import { useAuth } from '../../context/AuthContext'
import { ScopeConsentPrompt, ScopeInfo } from '../../components/ScopeConsentPrompt'

const { Text } = Typography

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

interface OverviewTabProps {
  sheet: Sheet
  onCopy: (text: string) => void
  onNavigateToApiSettings?: () => void
  onPublish?: (values: { default_range?: string; use_first_row_as_header: boolean }) => void
  onUnpublish?: () => void
  isPublishing?: boolean
  isUnpublishing?: boolean
}

export function OverviewTab({ sheet, onCopy, onNavigateToApiSettings, onPublish, onUnpublish, isPublishing, isUnpublishing }: OverviewTabProps) {
  const [isPublishModalOpen, setIsPublishModalOpen] = useState(false)
  const [showScopePrompt, setShowScopePrompt] = useState(false)
  const [form] = Form.useForm()
  const { hasScope, requestScopes } = useAuth()
  const workerBaseUrl = import.meta.env.VITE_WORKER_BASE_URL || 'https://api.gsheetbase.com'
  const apiUrl = sheet.api_key ? `${workerBaseUrl}/v1/${sheet.api_key}` : null

  const canRead = hasScope('https://www.googleapis.com/auth/spreadsheets.readonly')
  const canWrite = hasScope('https://www.googleapis.com/auth/spreadsheets')

  const operations = [
    {
      name: 'Read Data',
      scope: 'https://www.googleapis.com/auth/spreadsheets.readonly',
      available: canRead,
      methods: ['GET'],
      description: 'Fetch rows from your sheet',
    },
    {
      name: 'Write Data',
      scope: 'https://www.googleapis.com/auth/spreadsheets',
      available: canWrite,
      methods: ['POST', 'PUT', 'PATCH'],
      description: 'Add, update, or modify sheet rows',
      comingSoon: true,
    },
  ]

  const handleRequestWriteAccess = async () => {
    try {
      await requestScopes(['https://www.googleapis.com/auth/spreadsheets'])
      setShowScopePrompt(false)
    } catch (error) {
      console.error('Failed to request write scope:', error)
    }
  }

  const scopeInfo: ScopeInfo[] = [
    {
      scope: 'spreadsheets',
      reason: 'Write access to your Google Sheets',
      example: 'Allows adding, updating, and deleting rows via API',
    },
  ]

  const handlePublishSubmit = () => {
    form.validateFields().then((values) => {
      onPublish?.(values)
      setIsPublishModalOpen(false)
      form.resetFields()
    })
  }

  return (
    <>
      <Card 
        title="Sheet Information"
        extra={
          !sheet.is_public ? (
            <Button
              type="primary"
              size="small"
              icon={<CheckCircleOutlined />}
              onClick={() => setIsPublishModalOpen(true)}
            >
              Publish Sheet
            </Button>
          ) : (
            <Popconfirm
              title="Unpublish this sheet?"
              description="This will revoke the API key and make the sheet private."
              onConfirm={onUnpublish}
              okText="Yes"
              cancelText="No"
            >
              <Button
                type="default"
                size="small"
                danger
                loading={isUnpublishing}
              >
                Unpublish
              </Button>
            </Popconfirm>
          )
        }
      >
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

      <Divider>Available API Operations</Divider>

      <List
        dataSource={operations}
        renderItem={(op) => (
          <List.Item
            actions={[
              op.available ? (
                <Tag color="green">Active</Tag>
              ) : op.comingSoon ? (
                <Tag>Coming Soon</Tag>
              ) : (
                <Button 
                  size="small" 
                  onClick={() => setShowScopePrompt(true)}
                  icon={<LockOutlined />}
                >
                  Enable
                </Button>
              )
            ]}
          >
            <List.Item.Meta
              title={<strong>{op.name}</strong>}
              description={
                <Space direction="vertical" size="small" style={{ width: '100%' }}>
                  <Text type="secondary">{op.description}</Text>
                  <div>
                    {op.methods.map(m => (
                      <code key={m} style={{ 
                        marginRight: 8,
                        padding: '2px 6px',
                        background: '#f5f5f5',
                        borderRadius: 4,
                        fontSize: 12
                      }}>
                        {m}
                      </code>
                    ))}
                  </div>
                  {!op.available && !op.comingSoon && (
                    <Text type="warning" style={{ fontSize: 12 }}>
                      Requires additional Google permission
                    </Text>
                  )}
                </Space>
              }
            />
          </List.Item>
        )}
      />
    </Card>

      <Modal
        title="Publish Sheet"
        open={isPublishModalOpen}
        onOk={handlePublishSubmit}
        onCancel={() => {
          setIsPublishModalOpen(false)
          form.resetFields()
        }}
        okText="Publish"
        confirmLoading={isPublishing}
      >
        <Form
          form={form}
          layout="vertical"
          initialValues={{
            default_range: 'Sheet1',
            use_first_row_as_header: true,
          }}
        >
          <Form.Item
            label="Default Range"
            name="default_range"
            tooltip="Specify the default sheet/range to fetch (e.g., Sheet1, Sheet1!A1:Z100)"
          >
            <Input placeholder="e.g., Sheet1 or Sheet1!A1:Z100" />
          </Form.Item>

          <Form.Item
            label="Transform to JSON Objects"
            name="use_first_row_as_header"
            valuePropName="checked"
            tooltip="If enabled, the first row will be used as keys for JSON objects"
          >
            <Switch />
          </Form.Item>
        </Form>
      </Modal>

      <ScopeConsentPrompt
        open={showScopePrompt}
        onConsent={handleRequestWriteAccess}
        onCancel={() => setShowScopePrompt(false)}
        scopes={scopeInfo}
      />
    </>
  )
}
