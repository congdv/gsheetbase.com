import { GoogleOutlined } from '@ant-design/icons'
import {
  PageContainer,
  StyledCard,
  TitleContainer,
  StyledTitle,
  GoogleButton,
  Footer,
} from '../styles/global'

export default function LoginPage() {
  const googleStart = () => {
    const base = import.meta.env.VITE_API_BASE_URL
    window.location.href = `${base}/auth/google/start`
  }

  return (
    <PageContainer>
      <StyledCard>
        <TitleContainer>
          <StyledTitle level={2}>
            GSheetBase
          </StyledTitle>
        </TitleContainer>

        <GoogleButton
          icon={<GoogleOutlined />}
          onClick={googleStart}
          block
          size="large"
        >
          Continue with Google
        </GoogleButton>

        <Footer>
          © {new Date().getFullYear()} GSheetBase
        </Footer>
      </StyledCard>
    </PageContainer>
  )
}