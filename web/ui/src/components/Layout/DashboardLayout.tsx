import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout, Button, Avatar, Dropdown, Menu } from 'antd'
import { UserOutlined, LogoutOutlined, HomeOutlined, CreditCardOutlined } from '@ant-design/icons'
import styled from 'styled-components'
import { useAuth } from '../../context/AuthContext'
import { ROUTES } from '../../constants/routes'
import logo from '../../assets/gsheetbase.svg'

const { Header, Content, Sider } = Layout

const StyledLayout = styled(Layout)`
  min-height: 100vh;
  background: #fafafa;
`;

const StyledHeader = styled(Header)`
  && {
    background: #ffffff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 24px;
    position: sticky;
    top: 0;
    z-index: 100;
  }
`

const LogoContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 3px;
  cursor: pointer;
  
  img {
    height: 32px;
    width: auto;
  }
`

const LogoText = styled.div`
  font-size: 20px;
  font-weight: 600;
  color: #333;
`

const UserSection = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`

const UserEmail = styled.div`
  font-size: 14px;
  color: #666;
  
  @media (max-width: 768px) {
    display: none;
  }
`

const StyledSider = styled(Sider)`
  && {
    background: #ffffff;
    border-right: 1px solid #f0f0f0;
    
    .ant-layout-sider-children {
      display: flex;
      flex-direction: column;
    }
    
    .ant-menu {
      border-right: none;
    }
  }
`

const StyledContent = styled(Content)`
  background: #fafafa;
  padding: 24px;
  min-height: calc(100vh - 64px);
`

const MainLayout = styled(Layout)`
  min-height: calc(100vh - 64px);
`

export default function DashboardLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout } = useAuth()

  const handleLogout = async () => {
    await logout()
    navigate(ROUTES.LOGIN)
  }

  const userMenuItems = [
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: 'Logout',
      onClick: handleLogout
    }
  ]

  const menuItems = [
    {
      key: ROUTES.HOME,
      icon: <HomeOutlined />,
      label: 'Home',
      onClick: () => navigate(ROUTES.HOME)
    },
    {
      key: ROUTES.BILLING,
      icon: <CreditCardOutlined />,
      label: 'Billing',
      onClick: () => navigate(ROUTES.BILLING)
    }
  ]

  // Determine selected key based on current path
  const getSelectedKey = () => {
    if (location.pathname === ROUTES.HOME) return ROUTES.HOME
    if (location.pathname.startsWith(ROUTES.BILLING)) return ROUTES.BILLING
    return ROUTES.HOME
  }

  return (
    <StyledLayout>
      <StyledHeader>
        <LogoContainer onClick={() => navigate(ROUTES.HOME)}>
          <img src={logo} alt="gsheetbase Logo" />
          <LogoText>gsheetbase</LogoText>
        </LogoContainer>
        <UserSection>
          <UserEmail>{user?.email}</UserEmail>
          <Dropdown menu={{ items: userMenuItems }} trigger={['click']} placement="bottomRight">
            <Avatar icon={<UserOutlined />} style={{ cursor: 'pointer' }} />
          </Dropdown>
        </UserSection>
      </StyledHeader>

      <MainLayout>
        <StyledSider width={240} collapsedWidth={80} breakpoint="lg">
          <Menu
            mode="inline"
            selectedKeys={[getSelectedKey()]}
            items={menuItems}
            style={{ paddingTop: '16px' }}
          />
        </StyledSider>

        <StyledContent>
          <Outlet />
        </StyledContent>
      </MainLayout>
    </StyledLayout>
  )
}
