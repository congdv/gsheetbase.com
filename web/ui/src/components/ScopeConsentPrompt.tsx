import { Modal, List, Typography } from 'antd';
import { LockOutlined } from '@ant-design/icons';

const { Text } = Typography;

export interface ScopeInfo {
  scope: string;
  reason: string;
  example: string;
}

interface ScopeConsentPromptProps {
  open: boolean;
  onConsent: () => void;
  onCancel: () => void;
  scopes: ScopeInfo[];
}

export const ScopeConsentPrompt = ({
  open,
  onConsent,
  onCancel,
  scopes
}: ScopeConsentPromptProps) => {
  return (
    <Modal
      title={
        <span>
          <LockOutlined style={{ marginRight: 8 }} /> Additional Permission Required
        </span>
      }
      open={open}
      onOk={onConsent}
      onCancel={onCancel}
      okText="Grant Permission"
      cancelText="Cancel"
    >
      <p>To use this feature, Gsheetbase needs additional permissions:</p>
      <List
        dataSource={scopes}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              title={<strong>{item.reason}</strong>}
              description={<Text type="secondary">{item.example}</Text>}
            />
          </List.Item>
        )}
      />
      <Text type="secondary" style={{ marginTop: 16, display: 'block', fontSize: 12 }}>
        You'll be redirected to Google to grant these permissions.
      </Text>
    </Modal>
  );
};
