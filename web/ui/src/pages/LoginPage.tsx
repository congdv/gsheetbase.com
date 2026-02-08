import { GoogleOutlined } from '@ant-design/icons'
import logo from '../assets/gsheetbase.svg'
import { useConfig } from '../context/ConfigContext'
import SEO from '../components/SEO'
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
      <SEO
        title="Login"
        description="Sign in to Gsheetbase to turn your Google Sheets into REST APIs. Fast, secure, and easy to use."
        keywords="google sheets api login, sign in, authentication"
      />
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