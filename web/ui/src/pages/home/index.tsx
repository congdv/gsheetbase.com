import { useState, useEffect } from 'react'
import { Typography, Card, Button, Table, Modal, Form, Input, message, Popconfirm, Space } from 'antd'
import { PlusOutlined, DeleteOutlined, LinkOutlined } from '@ant-design/icons'
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
  created_at: string
}

export default function HomePage() {
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [form] = Form.useForm()
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

  const handleRegister = () => {
    form.validateFields().then((values) => {
      registerMutation.mutate(values)
    })
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
    </PageContainer>
  )
}
