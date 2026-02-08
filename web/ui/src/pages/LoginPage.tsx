import { GoogleOutlined } from '@ant-design/icons'
import logo from '../assets/gsheetbase.svg'
import { useConfig } from '../context/ConfigContext'
import {
  PageContainer,
  StyledCard,
  TitleContainer,
  StyledTitle,
  GoogleButton,
  Footer,
} from '../styles/global'

export default function LoginPage() {
  const config = useConfig()
  const googleStart = () => {
    window.location.href = `${config.apiBaseUrl}/auth/google/start`
  }

  return (
    <PageContainer>
      <StyledCard>
        <TitleContainer>
          <img src={logo} alt="Gsheetbase Logo" style={{ height: '48px', marginBottom: '16px' }} />
          <StyledTitle level={2}>
            Gsheetbase
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
          © {new Date().getFullYear()} Gsheetbase
        </Footer>
      </StyledCard>
    </PageContainer>
  )
}