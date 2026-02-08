import { Modal, List, Typography, Checkbox, Space } from 'antd';
import { LockOutlined } from '@ant-design/icons';
import { useEffect, useState } from 'react';

const { Text } = Typography;

export interface ScopeInfo {
  scope: string; // full OAuth scope string
  reason: string;
  example: string;
  optional?: boolean;
}

interface ScopeConsentPromptProps {
  open: boolean;
  onConsent: (selectedScopes: string[]) => void;
  onCancel: () => void;
  scopes: ScopeInfo[];
}

export const ScopeConsentPrompt = ({
  open,
  onConsent,
  onCancel,
  scopes
}: ScopeConsentPromptProps) => {
  const [selected, setSelected] = useState<string[]>([]);

  useEffect(() => {
    if (open) {
      // initialize selection: include all required scopes
      const initial = scopes.filter(s => !s.optional).map(s => s.scope);
      setSelected(initial);
    }
  }, [open, scopes]);

  const toggle = (scope: string, checked: boolean) => {
    setSelected(prev => {
      if (checked) return Array.from(new Set([...prev, scope]));
      return prev.filter(s => s !== scope);
    });
  };

  return (
    <Modal
      title={
        <span>
          <LockOutlined style={{ marginRight: 8 }} /> Additional Permission Required
        </span>
      }
      open={open}
      onOk={() => onConsent(selected)}
      onCancel={onCancel}
      okText="Grant Permission"
      cancelText="Cancel"
      okButtonProps={{ disabled: selected.length === 0 }}
    >
      <p>To use this feature, Gsheetbase needs additional permissions:</p>
      <List
        dataSource={scopes}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              title={
                <Space>
                  <Checkbox
                    checked={selected.includes(item.scope)}
                    onChange={(e) => toggle(item.scope, e.target.checked)}
                    disabled={!item.optional}
                  />
                  <strong>{item.reason}</strong>
                </Space>
              }
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
