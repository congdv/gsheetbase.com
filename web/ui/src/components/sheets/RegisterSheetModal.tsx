import { Modal, Form, Input } from 'antd'
import { useEffect } from 'react'

interface RegisterSheetModalProps {
  open: boolean
  isRegistering: boolean
  onCancel: () => void
  onRegister: (values: { sheet_id: string; sheet_name?: string; description?: string }) => void
}

const extractSheetId = (url: string): string => {
  // Extract sheet ID from Google Sheets URL
  const match = url.match(/\/d\/([a-zA-Z0-9-_]+)/)
  return match ? match[1] : url
}

export const RegisterSheetModal = ({
  open,
  isRegistering,
  onCancel,
  onRegister,
}: RegisterSheetModalProps) => {
  const [form] = Form.useForm()

  useEffect(() => {
    if (!open) {
      form.resetFields()
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
    onCancel()
  }

  return (
    <Modal
      title="Register New Sheet"
      open={open}
      onCancel={handleCancel}
      onOk={handleOk}
      confirmLoading={isRegistering}
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
        <Form.Item label="Sheet Name (optional)" name="sheet_name">
          <Input placeholder="My awesome sheet" />
        </Form.Item>
        <Form.Item label="Description (optional)" name="description">
          <Input.TextArea placeholder="What this sheet is used for..." rows={3} />
        </Form.Item>
      </Form>
    </Modal>
  )
}
