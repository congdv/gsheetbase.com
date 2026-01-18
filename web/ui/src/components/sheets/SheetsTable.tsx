import { Table, Button, Space, Tag, Popconfirm, Typography } from 'antd'
import { LinkOutlined, ApiOutlined, DeleteOutlined } from '@ant-design/icons'
import { useNavigate } from 'react-router-dom'
import type { ColumnsType } from 'antd/es/table'
import { Sheet } from '../../hooks/useSheets'
import { ROUTES } from '../../constants/routes'
import styled from 'styled-components'

const { Text } = Typography

const ClickableRow = styled.tr`
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #f5f5f5;
  }
`

interface SheetsTableProps {
  sheets: Sheet[]
  isLoading: boolean
  onOpenApiSettings: (sheet: Sheet) => void
  onUnpublish: (sheetId: string) => void
  onDelete: (sheetId: string) => void
}

export const SheetsTable = ({
  sheets,
  isLoading,
  onOpenApiSettings,
  onUnpublish,
  onDelete,
}: SheetsTableProps) => {
  const navigate = useNavigate()

  const handleRowClick = (record: Sheet) => {
    navigate(ROUTES.SHEET_DETAIL.replace(':id', record.id))
  }

  const columns: ColumnsType<Sheet> = [
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
      render: (_: any, record: Sheet) =>
        record.is_public ? <Tag color="green">Public</Tag> : <Tag>Private</Tag>,
    },
    {
      title: 'Sheet ID',
      dataIndex: 'sheet_id',
      key: 'sheet_id',
      render: (text: string) => <Text copyable={{ text }}>{text.substring(0, 20)}...</Text>,
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
        <Space onClick={(e) => e.stopPropagation()}>
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
                onClick={() => onOpenApiSettings(record)}
              >
                API Settings
              </Button>
              <Popconfirm
                title="Unpublish this sheet?"
                description="This will revoke the API key and make the sheet private."
                onConfirm={() => onUnpublish(record.id)}
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
              onClick={() => onOpenApiSettings(record)}
            >
              Publish
            </Button>
          )}
          <Popconfirm
            title="Remove this sheet?"
            description="This will remove the sheet from your allowed list."
            onConfirm={() => onDelete(record.sheet_id)}
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
    <Table
      columns={columns}
      dataSource={sheets}
      loading={isLoading}
      rowKey="id"
      pagination={{ pageSize: 10 }}
      onRow={(record) => ({
        onClick: () => handleRowClick(record),
      })}
      components={{
        body: {
          row: ClickableRow,
        },
      }}
    />
  )
}
