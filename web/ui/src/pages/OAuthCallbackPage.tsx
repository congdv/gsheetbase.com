import { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Spin, Typography } from 'antd'
import { LoadingOutlined } from '@ant-design/icons'
import api from '../lib/axios'
import { useAuth } from '../context/AuthContext'

const { Title } = Typography

export default function OAuthCallbackPage() {
  const navigate = useNavigate()
  const { checkSession } = useAuth()

  useEffect(() => {
    const run = async () => {
      try {
        // Backend set the access token cookie; check the session
        await checkSession()
        navigate('/home', { replace: true })
      } catch {
        navigate('/login', { replace: true })
      }
    }
    run()
  }, [navigate, checkSession])

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      minHeight: '200px',
      textAlign: 'center'
    }}>
      <Spin
        indicator={<LoadingOutlined style={{ fontSize: 48 }} spin />}
        size="large"
      />
      <Title level={3} style={{ marginTop: 16, color: '#666' }}>
        Signing you in...
      </Title>
    </div>
  )
}