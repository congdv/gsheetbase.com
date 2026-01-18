import { useState, useEffect } from 'react'
import { Typography, Card, Button, Table, Modal, Form, Input, message, Popconfirm, Space, Switch, Tag, Tooltip } from 'antd'
import { PlusOutlined, DeleteOutlined, LinkOutlined, ApiOutlined, CopyOutlined } from '@ant-design/icons'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import styled from 'styled-components'
import api from '../../lib/axios'

const { Title, Paragraph, Text } = Typography

const PageContainer = styled.div`
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
`

const PageHeader = styled.div`
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
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

export default function HomePage() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [isApiModalOpen, setIsApiModalOpen] = useState(false)
  const [selectedSheet, setSelectedSheet] = useState<Sheet | null>(null)
  const [form] = Form.useForm()
  const [apiForm] = Form.useForm()
  const queryClient = useQueryClient()

  // Fetch registered sheets
  const { data: sheetsData, isLoading } = useQuery({
    queryKey: ['sheets'],
    queryFn: async () => {
      const response = await api.get('/sheets/registered')
      return response.data.sheets || []
    },
  })

  // Register new sheet
  const registerMutation = useMutation({
    mutationFn: async (values: { sheet_id: string; sheet_name?: string; description?: string }) => {
      const response = await api.post('/sheets/register', values)
      return response.data
    },
    onSuccess: () => {
      message.success('Sheet registered successfully!')
      queryClient.invalidateQueries({ queryKey: ['sheets'] })
      setIsModalOpen(false)
      form.resetFields()
    },
    onError: (error: any) => {
      message.error(error?.response?.data?.error || 'Failed to register sheet')
    },
  })

  // Delete sheet
  const deleteMutation = useMutation({
    mutationFn: async (sheetId: string) => {
      await api.delete(`/sheets/registered/${sheetId}`)
    },
    onSuccess: () => {
      message.success('Sheet removed successfully!')
      queryClient.invalidateQueries({ queryKey: ['sheets'] })
    },
    onError: (error: any) => {
      message.error(error?.response?.data?.error || 'Failed to remove sheet')
    },
  })

  // Publish sheet
  const publishMutation = useMutation({
    mutationFn: async ({ sheetId, values }: { sheetId: string; values: any }) => {
      const response = await api.post(`/sheets/${sheetId}/publish`, values)
      return response.data
    },
    onSuccess: (data) => {
      message.success('Sheet published successfully!')
      queryClient.invalidateQueries({ queryKey: ['sheets'] })
      setIsApiModalOpen(false)
      apiForm.resetFields()
    },
    onError: (error: any) => {
      message.error(error?.response?.data?.error || 'Failed to publish sheet')
    },
  })

  // Unpublish sheet
  const unpublishMutation = useMutation({
    mutationFn: async (sheetId: string) => {
      await api.delete(`/sheets/${sheetId}/unpublish`)
    },
    onSuccess: () => {
      message.success('Sheet unpublished successfully!')
      queryClient.invalidateQueries({ queryKey: ['sheets'] })
    },
    onError: (error: any) => {
      message.error(error?.response?.data?.error || 'Failed to unpublish sheet')
    },
  })

  const handleRegister = () => {
    form.validateFields().then((values) => {
      registerMutation.mutate(values)
    })
  }

  const handlePublish = () => {
    if (!selectedSheet) return
    apiForm.validateFields().then((values) => {
      publishMutation.mutate({ sheetId: selectedSheet.id, values })
    })
  }

  const handleOpenApiSettings = (sheet: Sheet) => {
    setSelectedSheet(sheet)
    apiForm.setFieldsValue({
      default_range: sheet.default_range || 'Sheet1',
      use_first_row_as_header: sheet.use_first_row_as_header !== false,
    })
    setIsApiModalOpen(true)
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    message.success('Copied to clipboard!')
  }

  const extractSheetId = (url: string): string => {
    // Extract sheet ID from Google Sheets URL
    const match = url.match(/\/d\/([a-zA-Z0-9-_]+)/)
    return match ? match[1] : url
  }

  const columns = [
    {
      title: 'Sheet Name',
      dataIndex: 'sheet_name',
      key: 'sheet_name',
      render: (text: string, record: Sheet) => text || record.sheet_id,
    },
    {
      title: 'Description',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: 'API Status',
      key: 'api_status',
      render: (_: any, record: Sheet) => (
        record.is_public ? (
          <Tag color="green">Public</Tag>
        ) : (
          <Tag>Private</Tag>
        )
      ),
    },
    {
      title: 'Sheet ID',
      dataIndex: 'sheet_id',
      key: 'sheet_id',
      render: (text: string) => (
        <Text copyable={{ text }}>{text.substring(0, 20)}...</Text>
      ),
    },
    {
      title: 'Registered',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text: string) => new Date(text).toLocaleDateString(),
    },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: any, record: Sheet) => (
        <Space>
          <Button
            type="link"
            icon={<LinkOutlined />}
            href={`https://docs.google.com/spreadsheets/d/${record.sheet_id}`}
            target="_blank"
          >
            Open
          </Button>
          {record.is_public ? (
            <>
              <Button
                type="link"
                icon={<ApiOutlined />}
                onClick={() => handleOpenApiSettings(record)}
              >
                API Settings
              </Button>
              <Popconfirm
                title="Unpublish this sheet?"
                description="This will revoke the API key and make the sheet private."
                onConfirm={() => unpublishMutation.mutate(record.id)}
                okText="Yes"
                cancelText="No"
              >
                <Button type="link" danger>
                  Unpublish
                </Button>
              </Popconfirm>
            </>
          ) : (
            <Button
              type="link"
              icon={<ApiOutlined />}
              onClick={() => handleOpenApiSettings(record)}
            >
              Publish
            </Button>
          )}
          <Popconfirm
            title="Remove this sheet?"
            description="This will remove the sheet from your allowed list."
            onConfirm={() => deleteMutation.mutate(record.sheet_id)}
            okText="Yes"
            cancelText="No"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              Remove
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <PageContainer>
      <PageHeader>
        <div>
          <Title level={2}>My Sheets</Title>
          <Paragraph>Manage your Google Sheets connections</Paragraph>
        </div>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={() => setIsModalOpen(true)}
          size="large"
        >
          Register Sheet
        </Button>
      </PageHeader>

      <Card>
        <Table
          columns={columns}
          dataSource={sheetsData}
          loading={isLoading}
          rowKey="id"
          pagination={{ pageSize: 10 }}
        />
      </Card>

      <Modal
        title="Register New Sheet"
        open={isModalOpen}
        onCancel={() => {
          setIsModalOpen(false)
          form.resetFields()
        }}
        onOk={handleRegister}
        confirmLoading={registerMutation.isPending}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            label="Sheet URL or ID"
            name="sheet_id"
            rules={[{ required: true, message: 'Please enter the Google Sheets URL or ID' }]}
            getValueFromEvent={(e) => {
              const value = e.target.value.trim()
              return extractSheetId(value)
            }}
          >
            <Input
              placeholder="https://docs.google.com/spreadsheets/d/YOUR_SHEET_ID/edit"
              onChange={(e) => {
                const extracted = extractSheetId(e.target.value.trim())
                form.setFieldValue('sheet_id', extracted)
              }}
            />
          </Form.Item>
          <Form.Item
            label="Sheet Name (optional)"
            name="sheet_name"
          >
            <Input placeholder="My awesome sheet" />
          </Form.Item>
          <Form.Item
            label="Description (optional)"
            name="description"
          >
            <Input.TextArea placeholder="What this sheet is used for..." rows={3} />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title={selectedSheet?.is_public ? 'API Settings' : 'Publish Sheet as API'}
        open={isApiModalOpen}
        onCancel={() => {
          setIsApiModalOpen(false)
          setSelectedSheet(null)
          apiForm.resetFields()
        }}
        onOk={selectedSheet?.is_public ? undefined : handlePublish}
        confirmLoading={publishMutation.isPending}
        footer={selectedSheet?.is_public ? [
          <Button key="close" onClick={() => setIsApiModalOpen(false)}>
            Close
          </Button>
        ] : undefined}
        width={700}
      >
        {selectedSheet?.is_public && selectedSheet.api_key ? (
          <div style={{ marginBottom: 24 }}>
            <Typography.Title level={5}>API Endpoint</Typography.Title>
            <Input
              value={`https://api.gsheetbase.com/v1/${selectedSheet.api_key}`}
              readOnly
              addonAfter={
                <Tooltip title="Copy">
                  <CopyOutlined 
                    onClick={() => copyToClipboard(`https://api.gsheetbase.com/v1/${selectedSheet.api_key}`)}
                    style={{ cursor: 'pointer' }}
                  />
                </Tooltip>
              }
            />
            <Paragraph style={{ marginTop: 8, fontSize: 12, color: '#666' }}>
              Use this endpoint to access your sheet data. Add <code>?range=A1:Z100</code> to customize the range.
            </Paragraph>
            
            <Typography.Title level={5} style={{ marginTop: 16 }}>Example Usage</Typography.Title>
            <pre style={{ background: '#f5f5f5', padding: 12, borderRadius: 4, overflow: 'auto' }}>
{`// JavaScript
fetch('https://api.gsheetbase.com/v1/${selectedSheet.api_key}')
  .then(res => res.json())
  .then(data => console.log(data))

// cURL
curl https://api.gsheetbase.com/v1/${selectedSheet.api_key}`}
            </pre>

            <Typography.Title level={5} style={{ marginTop: 16 }}>Settings</Typography.Title>
            <p><strong>Default Range:</strong> {selectedSheet.default_range || 'Sheet1'}</p>
            <p><strong>Transform to JSON Objects:</strong> {selectedSheet.use_first_row_as_header ? 'Yes' : 'No'}</p>
          </div>
        ) : (
          <Form form={apiForm} layout="vertical" initialValues={{ use_first_row_as_header: true, default_range: 'Sheet1' }}>
            <Typography.Paragraph>
              Publishing your sheet will generate an API key and make it publicly accessible via a REST endpoint.
            </Typography.Paragraph>
            
            <Form.Item
              label="Default Range"
              name="default_range"
              tooltip="The default sheet range to fetch (e.g., Sheet1, Sheet1!A1:Z100)"
            >
              <Input placeholder="Sheet1" />
            </Form.Item>
            
            <Form.Item
              label="Use First Row as Headers"
              name="use_first_row_as_header"
              valuePropName="checked"
              tooltip="Convert the first row into JSON object keys"
            >
              <Switch />
            </Form.Item>
            
            <Typography.Paragraph type="secondary" style={{ fontSize: 12 }}>
              When enabled, the API will transform rows like <code>[["Name", "Age"], ["Alice", 30]]</code> 
              into JSON objects: <code>[{"{"}Name: "Alice", Age: 30{"}"}]</code>
            </Typography.Paragraph>
          </Form>
        )}
      </Modal>
    </PageContainer>
  )
}
