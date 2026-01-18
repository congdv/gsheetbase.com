import { Modal, Form, Input, Switch, Typography, Button, Tooltip } from 'antd'
import { CopyOutlined } from '@ant-design/icons'
import { message } from 'antd'
import { useEffect } from 'react'
import { Sheet } from '../../hooks/useSheets'

const { Paragraph } = Typography

interface ApiSettingsModalProps {
  open: boolean
  sheet: Sheet | null
  isPublishing: boolean
  onCancel: () => void
  onPublish: (values: { default_range?: string; use_first_row_as_header: boolean }) => void
}

export const ApiSettingsModal = ({
  open,
  sheet,
  isPublishing,
  onCancel,
  onPublish,
}: ApiSettingsModalProps) => {
  const [form] = Form.useForm()

  useEffect(() => {
    if (open && sheet) {
      form.setFieldsValue({
        default_range: sheet.default_range || 'Sheet1',
        use_first_row_as_header: sheet.use_first_row_as_header !== false,
      })
    }
  }, [open, sheet, form])

  const handleOk = () => {
    form.validateFields().then((values) => {
      onPublish(values)
    })
  }

  const handleCancel = () => {
    form.resetFields()
    onCancel()
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    message.success('Copied to clipboard!')
  }

  const isPublishedSheet = sheet?.is_public && sheet?.api_key

  return (
    <Modal
      title={isPublishedSheet ? 'API Settings' : 'Publish Sheet as API'}
      open={open}
      onCancel={handleCancel}
      onOk={isPublishedSheet ? undefined : handleOk}
      confirmLoading={isPublishing}
      footer={
        isPublishedSheet
          ? [
              <Button key="close" onClick={handleCancel}>
                Close
              </Button>,
            ]
          : undefined
      }
      width={700}
    >
      {isPublishedSheet ? (
        <div style={{ marginBottom: 24 }}>
          <Typography.Title level={5}>API Endpoint</Typography.Title>
          <Input
            value={`https://api.gsheetbase.com/v1/${sheet.api_key}`}
            readOnly
            addonAfter={
              <Tooltip title="Copy">
                <CopyOutlined
                  onClick={() =>
                    copyToClipboard(`https://api.gsheetbase.com/v1/${sheet.api_key}`)
                  }
                  style={{ cursor: 'pointer' }}
                />
              </Tooltip>
            }
          />
          <Paragraph style={{ marginTop: 8, fontSize: 12, color: '#666' }}>
            Use this endpoint to access your sheet data. Add <code>?range=A1:Z100</code> to
            customize the range.
          </Paragraph>

          <Typography.Title level={5} style={{ marginTop: 16 }}>
            Example Usage
          </Typography.Title>
          <pre
            style={{ background: '#f5f5f5', padding: 12, borderRadius: 4, overflow: 'auto' }}
          >
            {`// JavaScript
fetch('https://api.gsheetbase.com/v1/${sheet.api_key}')
  .then(res => res.json())
  .then(data => console.log(data))

// cURL
curl https://api.gsheetbase.com/v1/${sheet.api_key}`}
          </pre>

          <Typography.Title level={5} style={{ marginTop: 16 }}>
            Settings
          </Typography.Title>
          <p>
            <strong>Default Range:</strong> {sheet.default_range || 'Sheet1'}
          </p>
          <p>
            <strong>Transform to JSON Objects:</strong>{' '}
            {sheet.use_first_row_as_header ? 'Yes' : 'No'}
          </p>
        </div>
      ) : (
        <Form
          form={form}
          layout="vertical"
          initialValues={{ use_first_row_as_header: true, default_range: 'Sheet1' }}
        >
          <Typography.Paragraph>
            Publishing your sheet will generate an API key and make it publicly accessible via a
            REST endpoint.
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
            When enabled, the API will transform rows like{' '}
            <code>[["Name", "Age"], ["Alice", 30]]</code> into JSON objects:{' '}
            <code>[{'{'}Name: "Alice", Age: 30{'}'}]</code>
          </Typography.Paragraph>
        </Form>
      )}
    </Modal>
  )
}
