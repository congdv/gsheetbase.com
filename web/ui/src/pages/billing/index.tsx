import { Typography, Card, Empty } from 'antd'
import { CreditCardOutlined } from '@ant-design/icons'
import styled from 'styled-components'

const { Title, Paragraph } = Typography

const PageContainer = styled.div`
  margin: 0 auto;
`

const PlaceholderCard = styled(Card)`
  text-align: center;
  padding: 48px 24px;
  
  .ant-empty-description {
    font-size: 16px;
    color: #666;
  }
`

const IconWrapper = styled.div`
  font-size: 64px;
  color: #d9d9d9;
  margin-bottom: 24px;
`

export default function BillingPage() {
  return (
    <PageContainer>
      <Title level={2}>Billing</Title>
      <Paragraph type="secondary">
        Manage your subscription and billing information
      </Paragraph>

      <PlaceholderCard>
        <IconWrapper>
          <CreditCardOutlined />
        </IconWrapper>
        <Empty
          description={
            <>
              <div style={{ marginBottom: '8px', fontSize: '18px', fontWeight: 500 }}>
                Coming Soon
              </div>
              <div>
                Billing and subscription management will be available here
              </div>
            </>
          }
        />
      </PlaceholderCard>
    </PageContainer>
  )
}
