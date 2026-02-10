import {
    Card,
    Space,
    Typography,
    Button,
    Select,
    Input,
    Form,
    Modal,
    Alert,
    Tag,
    Row,
    Col,
    Spin,
    message,
    Popconfirm,
    Divider,
} from 'antd'
import {
    CopyOutlined,
    LockOutlined,
    UnlockOutlined,
    ReloadOutlined,
    DeleteOutlined,
} from '@ant-design/icons'
import { useState, useEffect } from 'react'
import { useConfig } from '../../context/ConfigContext'
import api from '../../lib/axios'
import { queryClient } from '../../lib/queryClient'

const { Text, Paragraph } = Typography

interface AuthManagementCardProps {
    sheetId: string
    currentAuthType: string
    sheet?: { api_key?: string }
}

interface AuthStatus {
    auth_type: string
    auth_bearer_token_set: boolean
    auth_basic_username_set: boolean
    auth_basic_password_set: boolean
}

export function AuthManagementCard({ sheetId, currentAuthType, sheet }: AuthManagementCardProps) {
    const [form] = Form.useForm()
    const [authStatus, setAuthStatus] = useState<AuthStatus | null>(null)
    const [loading, setLoading] = useState(false)
    const [saving, setSaving] = useState(false)
    const [generatedToken, setGeneratedToken] = useState<string | null>(null)
    const [showTokenModal, setShowTokenModal] = useState(false)
    const [showBasicModal, setShowBasicModal] = useState(false)

    const config = useConfig()
    const apiUrl = sheet?.api_key ? `${config.workerBaseUrl}/v1/${sheet.api_key}` : ''

    // Fetch current auth status
    const fetchAuthStatus = async () => {
        setLoading(true)
        try {
            const res = await api.get(`/sheets/${sheetId}/auth`)
            setAuthStatus(res.data.auth)
        } catch (error: any) {
            message.error('Failed to fetch auth status')
            console.error(error)
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchAuthStatus()
    }, [sheetId])

    const handleSetAuthType = async (authType: string) => {
        setSaving(true)
        try {
            await api.post(`/sheets/${sheetId}/auth/type`, { auth_type: authType })
            message.success(`Auth type changed to: ${authType}`)
            await fetchAuthStatus()
            await queryClient.invalidateQueries({ queryKey: ['sheets'] })
        } catch (error: any) {
            message.error(error.response?.data?.error || 'Failed to set auth type')
        } finally {
            setSaving(false)
        }
    }

    const handleGenerateBearerToken = async () => {
        setSaving(true)
        try {
            const res = await api.post(`/sheets/${sheetId}/auth/bearer`)
            setGeneratedToken(res.data.token.token)
            setShowTokenModal(true)
            message.success('Bearer token generated')
            await fetchAuthStatus()
            await queryClient.invalidateQueries({ queryKey: ['sheets'] })
        } catch (error: any) {
            message.error(error.response?.data?.error || 'Failed to generate token')
        } finally {
            setSaving(false)
        }
    }

    const handleRotateToken = async () => {
        setSaving(true)
        try {
            const res = await api.post(`/sheets/${sheetId}/auth/bearer/rotate`)
            setGeneratedToken(res.data.token.token)
            setShowTokenModal(true)
            message.success('Bearer token rotated successfully')
            await fetchAuthStatus()
            await queryClient.invalidateQueries({ queryKey: ['sheets'] })
        } catch (error: any) {
            message.error(error.response?.data?.error || 'Failed to rotate token')
        } finally {
            setSaving(false)
        }
    }

    const handleSetBasicAuth = async () => {
        try {
            const values = await form.validateFields()
            setSaving(true)
            await api.post(`/sheets/${sheetId}/auth/basic`, {
                username: values.username,
                password: values.password,
            })
            message.success('Basic auth credentials set successfully')
            form.resetFields()
            setShowBasicModal(false)
            await fetchAuthStatus()
            await queryClient.invalidateQueries({ queryKey: ['sheets'] })
        } catch (error: any) {
            if (error.response?.data?.error) {
                message.error(error.response.data.error)
            }
        } finally {
            setSaving(false)
        }
    }

    const handleDisableAuth = async () => {
        setSaving(true)
        try {
            await api.delete(`/sheets/${sheetId}/auth`)
            message.success('Authentication disabled')
            await fetchAuthStatus()
            await queryClient.invalidateQueries({ queryKey: ['sheets'] })
        } catch (error: any) {
            message.error(error.response?.data?.error || 'Failed to disable auth')
        } finally {
            setSaving(false)
        }
    }

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text)
        message.success('Copied to clipboard!')
    }

    if (loading) {
        return (
            <Card title="Authentication" size="small">
                <Spin />
            </Card>
        )
    }

    const authType = authStatus?.auth_type || currentAuthType || 'none'

    return (
        <Card title="API Authentication" size="small">
            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                <Alert
                    message="Secure your API endpoint with authentication"
                    description="Choose how you want to control access to this sheet's data via the public API."
                    type="info"
                    showIcon
                />

                {/* Auth Type Selection */}
                <div>
                    <Text strong style={{ display: 'block', marginBottom: 8 }}>
                        Authentication Method
                    </Text>
                    <Space direction="vertical" style={{ width: '100%' }}>
                        <Space>
                            <Select
                                value={authType}
                                onChange={handleSetAuthType}
                                loading={saving}
                                style={{ width: 200 }}
                                options={[
                                    {
                                        label: 'None (Public)',
                                        value: 'none',
                                        description: 'Anyone can access (requires is_public=true)',
                                    },
                                    {
                                        label: 'Bearer Token',
                                        value: 'bearer',
                                        description: 'Use HTTP Bearer token in Authorization header',
                                    },
                                    {
                                        label: 'Basic Auth',
                                        value: 'basic',
                                        description: 'Use username and password (via Authorization header)',
                                    },
                                ]}
                            />
                            <Tag color={authType === 'none' ? 'default' : 'blue'}>{authType.toUpperCase()}</Tag>
                        </Space>
                        <Paragraph style={{ fontSize: 12, color: '#666', margin: 0 }}>
                            {authType === 'none' &&
                                'Public access – API is accessible to anyone (your sheet must be marked as public).'}
                            {authType === 'bearer' &&
                                'Secure access – Clients must include a Bearer token in the Authorization header.'}
                            {authType === 'basic' &&
                                'Secure access – Clients must provide username and password via HTTP Basic authentication.'}
                        </Paragraph>
                    </Space>
                </div>

                <Divider />

                {/* Bearer Token Management */}
                {authType === 'bearer' && (
                    <div>
                        <Text strong style={{ display: 'block', marginBottom: 8 }}>
                            <LockOutlined /> Bearer Token
                        </Text>

                        <Space direction="vertical" style={{ width: '100%' }}>
                            <div>
                                {authStatus?.auth_bearer_token_set ? (
                                    <Alert
                                        message="Bearer token configured"
                                        description="Your API is protected with a bearer token. Clients must include it in the Authorization header."
                                        type="success"
                                        showIcon
                                    />
                                ) : (
                                    <Alert
                                        message="No bearer token configured"
                                        description="Generate a new token to start using bearer authentication."
                                        type="warning"
                                        showIcon
                                    />
                                )}
                            </div>

                            <Space>
                                <Button
                                    type={authStatus?.auth_bearer_token_set ? 'default' : 'primary'}
                                    onClick={handleGenerateBearerToken}
                                    loading={saving}
                                >
                                    {authStatus?.auth_bearer_token_set ? 'Generate New Token' : 'Generate Token'}
                                </Button>
                                {authStatus?.auth_bearer_token_set && (
                                    <Button
                                        icon={<ReloadOutlined />}
                                        onClick={handleRotateToken}
                                        loading={saving}
                                    >
                                        Rotate Token
                                    </Button>
                                )}
                            </Space>

                            <Paragraph style={{ fontSize: 12, color: '#666' }}>
                                <strong>Usage:</strong> Include the token in your API requests:
                            </Paragraph>
                            <code style={{ display: 'block', marginTop: 8, padding: 8, background: '#f5f5f5', borderRadius: 4 }}>
                                {`curl -H "Authorization: Bearer YOUR_TOKEN" ${apiUrl}`}
                            </code>
                        </Space>
                    </div>
                )}

                {/* Basic Auth Management */}
                {authType === 'basic' && (
                    <div>
                        <Text strong style={{ display: 'block', marginBottom: 8 }}>
                            <LockOutlined /> Basic Authentication
                        </Text>

                        <Space direction="vertical" style={{ width: '100%' }}>
                            <div>
                                {authStatus?.auth_basic_username_set ? (
                                    <Alert
                                        message="Basic auth configured"
                                        description={`Username set: ${authStatus.auth_basic_username_set ? '●●●●●●' : 'Not set'}`}
                                        type="success"
                                        showIcon
                                    />
                                ) : (
                                    <Alert
                                        message="No basic auth configured"
                                        description="Set username and password to enable basic authentication."
                                        type="warning"
                                        showIcon
                                    />
                                )}
                            </div>

                            <Button
                                type={authStatus?.auth_basic_username_set ? 'default' : 'primary'}
                                onClick={() => {
                                    form.resetFields()
                                    setShowBasicModal(true)
                                }}
                            >
                                {authStatus?.auth_basic_username_set ? 'Update Credentials' : 'Set Credentials'}
                            </Button>

                            <Paragraph style={{ fontSize: 12, color: '#666' }}>
                                <strong>Usage:</strong> Include credentials in your API requests:
                            </Paragraph>
                            <code style={{ display: 'block', marginTop: 8, padding: 8, background: '#f5f5f5', borderRadius: 4 }}>
                                {`curl -u username:password ${apiUrl}`}
                            </code>
                        </Space>
                    </div>
                )}

                {/* Disable Auth Option */}
                {authType !== 'none' && (
                    <>
                        <Divider />
                        <Row justify="space-between" align="middle">
                            <Col>
                                <Space direction="vertical" size={0}>
                                    <Text strong>Disable Authentication</Text>
                                    <Paragraph style={{ fontSize: 12, color: '#666', margin: 0 }}>
                                        Remove all authentication requirements
                                    </Paragraph>
                                </Space>
                            </Col>
                            <Col>
                                <Popconfirm
                                    title="Disable Authentication?"
                                    description="This will remove all authentication and make your API publicly accessible. Continue?"
                                    okText="Yes, disable"
                                    cancelText="Cancel"
                                    onConfirm={handleDisableAuth}
                                >
                                    <Button danger icon={<UnlockOutlined />} loading={saving}>
                                        Disable Auth
                                    </Button>
                                </Popconfirm>
                            </Col>
                        </Row>
                    </>
                )}
            </Space>

            {/* Token Generated Modal */}
            <Modal
                title="Bearer Token Generated"
                open={showTokenModal}
                onCancel={() => {
                    setShowTokenModal(false)
                    setGeneratedToken(null)
                }}
                footer={[
                    <Button
                        key="close"
                        onClick={() => {
                            setShowTokenModal(false)
                            setGeneratedToken(null)
                        }}
                    >
                        Close
                    </Button>,
                ]}
            >
                <Alert
                    message="Save this token securely"
                    description="This token will not be shown again. If you lose it, you'll need to generate a new one."
                    type="warning"
                    showIcon
                    style={{ marginBottom: 16 }}
                />

                <Text strong style={{ display: 'block', marginBottom: 8 }}>
                    Your Bearer Token:
                </Text>

                <Input.Password
                    value={generatedToken || ''}
                    readOnly
                    addonAfter={
                        <CopyOutlined
                            onClick={() => generatedToken && copyToClipboard(generatedToken)}
                            style={{ cursor: 'pointer' }}
                        />
                    }
                    style={{ marginBottom: 16 }}
                />

                <Text style={{ display: 'block', marginBottom: 8 }} strong>
                    Usage Example:
                </Text>
                <code style={{ display: 'block', padding: 8, background: '#f5f5f5', borderRadius: 4 }}>
                    {`curl -H "Authorization: Bearer ${generatedToken?.substring(0, 20)}..." ${apiUrl}`}
                </code>
            </Modal>

            {/* Basic Auth Modal */}
            <Modal
                title={authStatus?.auth_basic_username_set ? 'Update Basic Auth Credentials' : 'Set Basic Auth Credentials'}
                open={showBasicModal}
                onOk={handleSetBasicAuth}
                onCancel={() => {
                    setShowBasicModal(false)
                    form.resetFields()
                }}
                confirmLoading={saving}
            >
                <Alert
                    message="Keep your credentials secure"
                    description="Only share these credentials with trusted clients. Always use HTTPS in production."
                    type="info"
                    showIcon
                    style={{ marginBottom: 16 }}
                />

                <Form form={form} layout="vertical">
                    <Form.Item
                        label="Username"
                        name="username"
                        rules={[
                            { required: true, message: 'Username is required' },
                            { min: 1, max: 255, message: 'Username must be between 1 and 255 characters' },
                        ]}
                    >
                        <Input placeholder="Enter username" />
                    </Form.Item>

                    <Form.Item
                        label="Password"
                        name="password"
                        rules={[
                            { required: true, message: 'Password is required' },
                            { min: 1, max: 255, message: 'Password must be between 1 and 255 characters' },
                        ]}
                    >
                        <Input.Password placeholder="Enter password" />
                    </Form.Item>
                </Form>
            </Modal>
        </Card>
    )
}
