import { useState } from 'react'
import { Typography, Card, Button } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import styled from 'styled-components'
import { useSheets } from '../../hooks/useSheets'
import { SheetsTable } from '../../components/sheets/SheetsTable'
import { RegisterSheetModal } from '../../components/sheets/RegisterSheetModal'
import SEO from '../../components/SEO'

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

  const {
    sheets,
    isLoading,
    registerSheet,
    isRegistering,
    deleteSheet,
  } = useSheets()

  const handleRegister = (values: { sheet_id: string; sheet_name?: string; description?: string }) => {
    registerSheet(values, {
      onSuccess: () => {
        setIsRegisterModalOpen(false)
      },
    })
  }

  return (
    <PageContainer>
      <SEO
        title="Dashboard"
        description="Manage your Google Sheets API connections. View, register, and configure REST API endpoints for your spreadsheets."
        noIndex={true}
      />
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
          onDelete={deleteSheet}
        />
      </Card>

      <RegisterSheetModal
        open={isRegisterModalOpen}
        isRegistering={isRegistering}
        onCancel={() => setIsRegisterModalOpen(false)}
        onRegister={handleRegister}
      />
    </PageContainer>
  )
}
