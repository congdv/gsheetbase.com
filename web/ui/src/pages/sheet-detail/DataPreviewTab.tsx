import { Card, Result } from 'antd'
import { TableOutlined } from '@ant-design/icons'

export function DataPreviewTab() {
  return (
    <Card>
      <Result
        icon={<TableOutlined style={{ color: '#1890ff' }} />}
        title="Data Preview Coming Soon"
        subTitle="This feature will allow you to preview your sheet data in a table format with customizable range selection and export options."
      />
    </Card>
  )
}
