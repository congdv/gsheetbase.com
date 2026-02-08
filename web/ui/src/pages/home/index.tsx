import { useState } from 'react'
import { useAuth } from '../../context/AuthContext'
import { Typography, Card, Button } from 'antd'
import { message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import styled from 'styled-components'
import { useSheets } from '../../hooks/useSheets'
import { SheetsTable } from '../../components/sheets/SheetsTable'
import { RegisterSheetModal } from '../../components/sheets/RegisterSheetModal'
import { GOOGLE_SCOPE } from '@/constants/common'
import { ScopeConsentPrompt, ScopeInfo } from '../../components/ScopeConsentPrompt'

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
  const { hasScope, requestScopes, isLoading: authLoading } = useAuth()
  const {
    sheets,
    isLoading,
    registerSheet,
    isRegistering,
    deleteSheet,
  } = useSheets()

  const hasSheetScope = hasScope(GOOGLE_SCOPE.READ_ONLY_SCOPE)
  const [showScopePrompt, setShowScopePrompt] = useState(false)

  const scopeList: ScopeInfo[] = [
    {
      scope: GOOGLE_SCOPE.READ_ONLY_SCOPE,
      reason: 'Read access to your Google Sheets',
      example: 'Allows viewing and registering sheets',
      optional: false,
    },
    {
      scope: GOOGLE_SCOPE.READ_WRITE_SCOPE,
      reason: 'Write access to your Google Sheets',
      example: 'Allows adding, updating, and deleting rows via API',
      optional: true,
    },
  ]

  const handleRegister = (values: { sheet_id: string; sheet_name?: string; description?: string }) => {
    registerSheet(values, {
      onSuccess: () => {
        setIsRegisterModalOpen(false)
      },
    })
  }

  const handleGrantAccess = async () => {
    setShowScopePrompt(true)
  }

  const handleScopeConsent = async (selectedScopes: string[]) => {
    try {
      await requestScopes(selectedScopes)
      setShowScopePrompt(false)
      message.success('Permission request started. Complete consent in the popup.')
    } catch (error: any) {
      console.error('requestScopes failed', error)
      message.error('Failed to start permission request')
    }
  }

  return (
    <PageContainer>
      {!authLoading && !hasSheetScope ? (
        <Card style={{ maxWidth: 480, margin: '48px auto', textAlign: 'center' }}>
          <Title level={3}>Google Sheets Access Required</Title>
          <Paragraph>
            To use Gsheetbase, please grant access to your Google Sheets. This is required to register and manage your APIs.
          </Paragraph>
          <Button type="primary" size="large" onClick={handleGrantAccess}>
            Grant Google Sheets Access
          </Button>
          <ScopeConsentPrompt
            open={showScopePrompt}
            onConsent={handleScopeConsent}
            onCancel={() => setShowScopePrompt(false)}
            scopes={scopeList}
          />
        </Card>
      ) : (
        <>
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
              disabled={!hasSheetScope}
            >
              Register Sheet
            </Button>
          </PageHeader>

          <Card>
            <SheetsTable
              sheets={sheets}
              isLoading={isLoading}
              onDelete={hasSheetScope ? deleteSheet : () => { }}
            />
          </Card>

          <RegisterSheetModal
            open={isRegisterModalOpen}
            isRegistering={isRegistering}
            onCancel={() => setIsRegisterModalOpen(false)}
            onRegister={hasSheetScope ? handleRegister : () => { }}
          />

        </>
      )}
    </PageContainer>
  )
}
