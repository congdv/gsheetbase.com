import { Modal, Form, Input, Button, Card, Typography, Space } from 'antd'
import { useEffect, useState } from 'react'
import { useCreateSheetMutation } from '../../hooks/useSheets'

interface RegisterSheetModalProps {
  open: boolean
  isRegistering: boolean
  onCancel: () => void
  onRegister: (values: { sheet_id: string; sheet_name?: string; description?: string }) => void
}


const { Text } = Typography;

const extractSheetId = (url: string): string => {
  // Extract sheet ID from Google Sheets URL
  const match = url.match(/\/d\/([a-zA-Z0-9-_]+)/)
  return match ? match[1] : url
}

// Phase 1 templates as plain text
const TEMPLATES = [
  {
    key: 'link-in-bio',
    name: 'Link-in-Bio',
    schema: `id, title, url, icon, order, active, description`
  },
  {
    key: 'directory',
    name: 'Directory',
    schema: `id, name, description, category, website_url, logo_url, featured, tags, created_at`
  },
  {
    key: 'job-board',
    name: 'Job Board',
    schema: `id, job_title, company, location, salary, job_type, description, apply_url, status, posted_date`
  },
  {
    key: 'lead-gen',
    name: 'Lead Gen',
    schema: `timestamp, full_name, email, interest_level, source_url, status, notes`
  },
  {
    key: 'services-page',
    name: 'Services Page',
    schema: `id, service_name, description, price, billing_type, package_name, turnaround, testimonial, client_name, booking_url`
  },
]

export const RegisterSheetModal = ({
  open,
  isRegistering,
  onCancel,
  onRegister,
}: RegisterSheetModalProps) => {
  const [form] = Form.useForm()
  const [showTemplate, setShowTemplate] = useState(false)
  const [selectedTemplate, setSelectedTemplate] = useState<string | null>(null)
  const [creatingSheet, setCreatingSheet] = useState(false)
  const createSheetMutation = useCreateSheetMutation()

  useEffect(() => {
    if (!open) {
      form.resetFields()
      setShowTemplate(false)
      setSelectedTemplate(null)
    }
  }, [open, form])

  const handleOk = () => {
    form.validateFields().then((values) => {
      onRegister(values)
      form.resetFields()
    })
  }

  const handleCancel = () => {
    form.resetFields()
    setShowTemplate(false)
    setSelectedTemplate(null)
    onCancel()
  }

  // Call backend to create sheet
  const handleCreateSheet = async () => {
    if (!selectedTemplate) return
    setCreatingSheet(true)
    createSheetMutation.mutate(
      { template: selectedTemplate },
      {
        onSuccess: (data) => {
          form.setFieldsValue({ sheet_id: data.sheet_id })
          setCreatingSheet(false)
          setShowTemplate(false)
        },
        onError: () => {
          setCreatingSheet(false)
        },
      }
    )
  }

  return (
    <Modal
      title="Register New Sheet"
      open={open}
      onCancel={handleCancel}
      onOk={handleOk}
      confirmLoading={isRegistering}
      footer={null}
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
            addonAfter={
              <Button size="small" onClick={() => setShowTemplate(true)}>
                Create New
              </Button>
            }
          />
        </Form.Item>
        <Form.Item label="Sheet Name (optional)" name="sheet_name">
          <Input placeholder="My awesome sheet" />
        </Form.Item>
        <Form.Item label="Description (optional)" name="description">
          <Input.TextArea placeholder="What this sheet is used for..." rows={3} />
        </Form.Item>
        <Form.Item>
          <Button type="primary" onClick={handleOk} loading={isRegistering} block>
            Register
          </Button>
        </Form.Item>
      </Form>

      {/* Template selection modal */}
      <Modal
        title="Create New Sheet"
        open={showTemplate}
        onCancel={() => setShowTemplate(false)}
        footer={null}
        width={600}
      >
        <Text strong>Select a template to create a new Google Sheet:</Text>
        <Space direction="vertical" style={{ width: '100%' , marginTop: '10px'}} size="large">
          <div
            style={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: 24,
              justifyContent: 'flex-start',
              alignItems: 'stretch',
            }}
          >
            {TEMPLATES.map((tpl) => (
              <Card
                key={tpl.key}
                title={tpl.name}
                bordered
                style={{
                  width: 240,
                  minHeight: 140,
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'space-between',
                  cursor: 'pointer',
                  borderColor: selectedTemplate === tpl.key ? '#1677ff' : '#f0f0f0',
                  boxShadow: selectedTemplate === tpl.key ? '0 0 0 2px #1677ff' : 'none',
                  background: selectedTemplate === tpl.key ? '#e6f4ff' : '#fff',
                  transition: 'border-color 0.2s, box-shadow 0.2s, background 0.2s',
                }}
                onClick={() => setSelectedTemplate(tpl.key)}
                hoverable
              >
                <div style={{ flex: 1, display: 'flex', alignItems: 'center' }}>
                  <Text type="secondary" style={{ fontSize: 12 }}>
                    {tpl.schema}
                  </Text>
                </div>
              </Card>
            ))}
          </div>
          <Button
            type="primary"
            disabled={!selectedTemplate}
            loading={creatingSheet}
            onClick={handleCreateSheet}
            block
          >
            Create Sheet
          </Button>
          <Text type="secondary" italic>
            Custom schema: <span style={{ color: '#aaa' }}>Coming soon</span>
          </Text>
        </Space>
      </Modal>
    </Modal>
  )
}
