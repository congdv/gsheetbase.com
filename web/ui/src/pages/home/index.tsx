import { useState } from 'react'
import { Typography, Card, Button } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import styled from 'styled-components'
import { useSheets, Sheet } from '../../hooks/useSheets'
import { SheetsTable } from '../../components/sheets/SheetsTable'
import { RegisterSheetModal } from '../../components/sheets/RegisterSheetModal'
import { ApiSettingsModal } from '../../components/sheets/ApiSettingsModal'

const { Title, Paragraph } = Typography

const PageContainer = styled.div`
  padding: 24px;
  margin: 0 auto;
`

const PageHeader = styled.div`
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
`

export default function HomePage() {
  const [isRegisterModalOpen, setIsRegisterModalOpen] = useState(false)
  const [isApiModalOpen, setIsApiModalOpen] = useState(false)
  const [selectedSheet, setSelectedSheet] = useState<Sheet | null>(null)

  const {
    sheets,
    isLoading,
    registerSheet,
    isRegistering,
    deleteSheet,
    publishSheet,
    isPublishing,
    unpublishSheet,
  } = useSheets()

  const handleRegister = (values: { sheet_id: string; sheet_name?: string; description?: string }) => {
    registerSheet(values, {
      onSuccess: () => {
        setIsRegisterModalOpen(false)
      },
    })
  }

  const handlePublish = (values: { default_range?: string; use_first_row_as_header: boolean }) => {
    if (!selectedSheet) return
    publishSheet(
      { sheetId: selectedSheet.id, values },
      {
        onSuccess: () => {
          setIsApiModalOpen(false)
          setSelectedSheet(null)
        },
      }
    )
  }

  const handleOpenApiSettings = (sheet: Sheet) => {
    setSelectedSheet(sheet)
    setIsApiModalOpen(true)
  }

  const handleCloseApiModal = () => {
    setIsApiModalOpen(false)
    setSelectedSheet(null)
  }

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
          onClick={() => setIsRegisterModalOpen(true)}
          size="large"
        >
          Register Sheet
        </Button>
      </PageHeader>

      <Card>
        <SheetsTable
          sheets={sheets}
          isLoading={isLoading}
          onOpenApiSettings={handleOpenApiSettings}
          onUnpublish={unpublishSheet}
          onDelete={deleteSheet}
        />
      </Card>

      <RegisterSheetModal
        open={isRegisterModalOpen}
        isRegistering={isRegistering}
        onCancel={() => setIsRegisterModalOpen(false)}
        onRegister={handleRegister}
      />

      <ApiSettingsModal
        open={isApiModalOpen}
        sheet={selectedSheet}
        isPublishing={isPublishing}
        onCancel={handleCloseApiModal}
        onPublish={handlePublish}
      />
    </PageContainer>
  )
}
