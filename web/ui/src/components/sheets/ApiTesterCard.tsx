import { Card, Space, Input, Button, Select, Typography, Alert } from 'antd'

const { Text } = Typography
import { SendOutlined } from '@ant-design/icons'
import { useState } from 'react'
import axios from 'axios'
import { useConfig } from '../../context/ConfigContext'

interface ApiTesterCardProps {
    apiKey: string
    isMethodEnabled: (method: string) => boolean
    authType?: string
    basicUsername?: string
}

export function ApiTesterCard({ apiKey, isMethodEnabled, authType = 'none', basicUsername }: ApiTesterCardProps) {
    const config = useConfig()
    const apiUrl = `${config.workerBaseUrl}/v1/${apiKey}`

    const [testUrl, setTestUrl] = useState(apiUrl)
    const [httpMethod, setHttpMethod] = useState<'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'>('GET')
    const [requestBody, setRequestBody] = useState('{\n  \n}')
    const [loading, setLoading] = useState(false)
    const [response, setResponse] = useState<{
        status: number
        data: any
        headers?: Record<string, string>
        error?: string
    } | null>(null)

    // Auth credential state (credentials are always user-entered; backend never exposes plaintext)
    const [bearerToken, setBearerToken] = useState('')
    const [basicUsernameInput, setBasicUsernameInput] = useState(basicUsername ?? '')
    const [basicPassword, setBasicPassword] = useState('')

    const buildAuthHeader = (): string | undefined => {
        if (authType === 'bearer' && bearerToken.trim()) {
            return `Bearer ${bearerToken.trim()}`
        }
        if (authType === 'basic' && basicUsernameInput.trim()) {
            return `Basic ${btoa(`${basicUsernameInput}:${basicPassword}`)}`
        }
        return undefined
    }

    const handleSendRequest = async () => {
        setLoading(true)
        setResponse(null)

        try {
            const startTime = Date.now()
            const authHeader = buildAuthHeader()
            const config: any = {
                method: httpMethod,
                url: testUrl,
                validateStatus: () => true,
                ...(authHeader ? { headers: { Authorization: authHeader } } : {}),
            }

            // Add request body for write methods only (not DELETE)
            if (['POST', 'PUT', 'PATCH'].includes(httpMethod)) {
                try {
                    config.data = JSON.parse(requestBody)
                } catch (e) {
                    setResponse({
                        status: 0,
                        data: null,
                        error: 'Invalid JSON in request body',
                    })
                    setLoading(false)
                    return
                }
            }

            const res = await axios(config)
            const endTime = Date.now()

            setResponse({
                status: res.status,
                data: res.data,
                headers: {
                    'content-type': res.headers['content-type'] || '',
                    'x-response-time': `${endTime - startTime}ms`,
                },
            })
        } catch (error: any) {
            setResponse({
                status: 0,
                data: null,
                error: error.message || 'Request failed',
            })
        } finally {
            setLoading(false)
        }
    }

    return (
        <Card title="API Tester" size="small">
            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                <div>
                    <Text strong style={{ display: 'block', marginBottom: 8 }}>Test Your API</Text>
                    <Space.Compact style={{ width: '100%' }}>
                        <Select
                            value={httpMethod}
                            onChange={setHttpMethod}
                            style={{ width: 100 }}
                            options={[
                                { label: 'GET', value: 'GET' },
                                { label: 'POST', value: 'POST', disabled: !isMethodEnabled('POST') },
                                { label: 'PUT', value: 'PUT', disabled: !isMethodEnabled('PUT') },
                                { label: 'PATCH', value: 'PATCH', disabled: !isMethodEnabled('PATCH') },
                                { label: 'DELETE', value: 'DELETE', disabled: !isMethodEnabled('DELETE') },
                            ]}
                        />
                        <Input
                            value={testUrl}
                            onChange={(e) => setTestUrl(e.target.value)}
                            placeholder="API URL"
                        />
                        <Button
                            type="primary"
                            icon={<SendOutlined />}
                            onClick={handleSendRequest}
                            loading={loading}
                        >
                            Send
                        </Button>
                    </Space.Compact>
                </div>

                {authType === 'bearer' && (
                    <div>
                        <Text strong style={{ display: 'block', marginBottom: 8 }}>Bearer Token</Text>
                        <Input.Password
                            value={bearerToken}
                            onChange={(e) => setBearerToken(e.target.value)}
                            placeholder="Enter your bearer token"
                        />
                        <Text type="secondary" style={{ fontSize: 12, display: 'block', marginTop: 4 }}>
                            The token is only shown once when generated. Copy it from the Auth Management section.
                        </Text>
                    </div>
                )}

                {authType === 'basic' && (
                    <div>
                        <Text strong style={{ display: 'block', marginBottom: 8 }}>Basic Auth Credentials</Text>
                        <Space direction="vertical" style={{ width: '100%' }}>
                            <Input
                                value={basicUsernameInput}
                                onChange={(e) => setBasicUsernameInput(e.target.value)}
                                placeholder="Username"
                                addonBefore="Username"
                            />
                            <Input.Password
                                value={basicPassword}
                                onChange={(e) => setBasicPassword(e.target.value)}
                                placeholder="Password"
                                addonBefore="Password"
                            />
                        </Space>
                        <Text type="secondary" style={{ fontSize: 12, display: 'block', marginTop: 4 }}>
                            The password is only shown once when set. Copy it from the Auth Management section.
                        </Text>
                    </div>
                )}

                {['POST', 'PUT', 'PATCH'].includes(httpMethod) && (
                    <div>
                        <Space direction="horizontal" style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 8 }}>
                            <Text strong>Request Body</Text>
                            <Button
                                size="small"
                                onClick={() => {
                                    const example = httpMethod === 'POST'
                                        ? '{\n  "data": [\n    ["John", "Doe", "30"],\n    ["Jane", "Smith", "25"]\n  ]\n}'
                                        : httpMethod === 'PUT'
                                            ? '{\n  "data": [\n    ["Updated", "Row", "1"],\n    ["Updated", "Row", "2"]\n  ],\n  "range": "Sheet1!A2:C3"\n}'
                                            : '{\n  "data": [\n    ["Patched", "Value"]\n  ],\n  "range": "Sheet1!A2:B2"\n}'
                                    setRequestBody(example)
                                }}
                            >
                                Load Example
                            </Button>
                        </Space>
                        <Input.TextArea
                            value={requestBody}
                            onChange={(e) => setRequestBody(e.target.value)}
                            placeholder="JSON request body"
                            rows={6}
                            style={{ fontFamily: 'monospace' }}
                        />
                        <Text type="secondary" style={{ fontSize: 12, display: 'block', marginTop: 4 }}>
                            {httpMethod === 'POST' && 'POST appends new rows to the sheet. Range is optional.'}
                            {httpMethod === 'PUT' && 'PUT updates existing rows. Range is required.'}
                            {httpMethod === 'PATCH' && 'PATCH partially updates rows. Range is required.'}
                            {httpMethod === 'DELETE' && 'DELETE removes rows. Range is required.'}
                        </Text>
                    </div>
                )}

                {response && (
                    <div>
                        <Space style={{ marginBottom: 8 }}>
                            <Text strong>Response:</Text>
                            <Text
                                type={response.status >= 200 && response.status < 300 ? 'success' : 'danger'}
                            >
                                Status {response.status}
                            </Text>
                            {response.headers?.['x-response-time'] && (
                                <Text type="secondary">{response.headers['x-response-time']}</Text>
                            )}
                        </Space>

                        {response.error ? (
                            <Alert
                                message="Request Failed"
                                description={response.error}
                                type="error"
                                showIcon
                            />
                        ) : (
                            <pre
                                style={{
                                    background: '#f5f5f5',
                                    padding: 12,
                                    borderRadius: 4,
                                    overflow: 'auto',
                                    maxHeight: 400,
                                    border: '1px solid #d9d9d9',
                                }}
                            >
                                {response.status === 204 ? "No content" : JSON.stringify(response.data, null, 2)}
                            </pre>
                        )}
                    </div>
                )}
            </Space>

            <div style={{ marginTop: 16 }}>
                <Text strong>Example Usage</Text>
                <pre
                    style={{
                        background: '#f5f5f5',
                        padding: 12,
                        borderRadius: 4,
                        overflow: 'auto',
                        marginTop: 8,
                    }}
                >
                    {(() => {
                        const authFetchHeader = authType === 'bearer'
                            ? `\n  headers: { Authorization: 'Bearer <your-bearer-token>' },`
                            : authType === 'basic'
                                ? `\n  headers: { Authorization: 'Basic <base64(username:password)>' },`
                                : ''
                        const authCurlHeader = authType === 'bearer'
                            ? ` \\\n  -H 'Authorization: Bearer <your-bearer-token>'`
                            : authType === 'basic'
                                ? ` \\\n  -H 'Authorization: Basic <base64(username:password)>'`
                                : ''

                        if (httpMethod === 'GET') return `// JavaScript - Fetch data
fetch('${apiUrl}', {${authFetchHeader}
})
  .then(res => res.json())
  .then(data => console.log(data))

// cURL - Fetch data
curl ${apiUrl}${authCurlHeader}`

                        if (httpMethod === 'POST') return `// JavaScript - Add new rows
fetch('${apiUrl}', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json'${authType !== 'none' ? `, Authorization: '${authType === 'bearer' ? 'Bearer <your-bearer-token>' : 'Basic <base64(username:password)>'}' ` : ''} },
  body: JSON.stringify({
    data: [
      ["John", "Doe", "30"],
      ["Jane", "Smith", "25"]
    ]
  })
})

// cURL - Add new rows
curl -X POST ${apiUrl}${authCurlHeader} \\
  -H 'Content-Type: application/json' \\
  -d '{"data": [["John", "Doe", "30"]]}'`

                        if (httpMethod === 'PUT') return `// JavaScript - Update rows
fetch('${apiUrl}', {
  method: 'PUT',
  headers: { 'Content-Type': 'application/json'${authType !== 'none' ? `, Authorization: '${authType === 'bearer' ? 'Bearer <your-bearer-token>' : 'Basic <base64(username:password)>'}' ` : ''} },
  body: JSON.stringify({
    data: [["Updated", "Row"]],
    range: "Sheet1!A2:B2"
  })
})

// cURL - Update rows
curl -X PUT ${apiUrl}${authCurlHeader} \\
  -H 'Content-Type: application/json' \\
  -d '{"data": [["Updated", "Row"]], "range": "Sheet1!A2:B2"}'`

                        if (httpMethod === 'PATCH') return `// JavaScript - Partial update
fetch('${apiUrl}', {
  method: 'PATCH',
  headers: { 'Content-Type': 'application/json'${authType !== 'none' ? `, Authorization: '${authType === 'bearer' ? 'Bearer <your-bearer-token>' : 'Basic <base64(username:password)>'}' ` : ''} },
  body: JSON.stringify({
    data: [["Patched"]],
    range: "Sheet1!A2"
  })
})

// cURL - Partial update
curl -X PATCH ${apiUrl}${authCurlHeader} \\
  -H 'Content-Type: application/json' \\
  -d '{"data": [["Patched"]], "range": "Sheet1!A2"}'`

                        if (httpMethod === 'DELETE') return `// JavaScript - Delete rows
fetch('${apiUrl}', {
  method: 'DELETE',
  headers: { 'Content-Type': 'application/json'${authType !== 'none' ? `, Authorization: '${authType === 'bearer' ? 'Bearer <your-bearer-token>' : 'Basic <base64(username:password)>'}' ` : ''} },
  body: JSON.stringify({
    range: "Sheet1!A2:B2"
  })
})

// cURL - Delete rows
curl -X DELETE ${apiUrl}${authCurlHeader} \\
  -H 'Content-Type: application/json' \\
  -d '{"range": "Sheet1!A2:B2"}'`

                        return ''
                    })()}
                </pre>
            </div>
        </Card>
    )
}
