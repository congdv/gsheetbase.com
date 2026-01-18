import styled from 'styled-components'
import { Card, Typography, Input, Button, Divider } from 'antd'
import { theme } from './theme'

const { Title, Link } = Typography

// Styled Components for Auth Pages (Login, Register)
export const PageContainer = styled.div`
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: ${theme.colors.background};
  padding: 20px;
`

export const StyledCard = styled(Card)`
  width: 100%;
  max-width: 400px;
  border-radius: ${theme.borderRadius.md}px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  background-color: ${theme.colors.cardBackground};
`

export const TitleContainer = styled.div`
  text-align: center;
  margin-bottom: ${theme.spacing.lg}px;
`

export const StyledTitle = styled(Title)`
  && {
    margin-bottom: 0;
    font-size: ${theme.fontSize.lg}px;
    font-weight: 600;
    color: ${theme.colors.textPrimary};
  }
`

export const FormLabel = styled.span`
  color: ${theme.colors.textSecondary};
  font-size: ${theme.fontSize.xs}px;
`

export const StyledInput = styled(Input)`
  && {
    border-radius: ${theme.borderRadius.sm}px;
    background-color: ${theme.colors.inputBackground};
    border: 1px solid ${theme.colors.inputBorder};
  }
`

export const StyledPasswordInput = styled(Input.Password)`
  && {
    border-radius: ${theme.borderRadius.sm}px;
    background-color: ${theme.colors.inputBackground};
    border: 1px solid ${theme.colors.inputBorder};
  }
`

export const ErrorMessage = styled.div`
  color: ${theme.colors.textError};
  margin-bottom: ${theme.spacing.sm}px;
  font-size: ${theme.fontSize.sm}px;
`

export const PrimaryButton = styled(Button)`
  && {
    border-radius: ${theme.borderRadius.sm}px;
    background-color: ${theme.colors.primary};
    border-color: ${theme.colors.primary};
    font-weight: 500;
    height: ${theme.heights.input}px;
  }
`

export const LinkContainer = styled.div`
  text-align: center;
`

export const StyledLink = styled(Link)`
  && {
    color: ${theme.colors.primaryLight};
    font-size: ${theme.fontSize.sm}px;
  }
`

export const StyledDivider = styled(Divider)`
  && {
    margin: ${theme.spacing.md}px 0;
  }
`

export const GoogleButton = styled(Button)`
  && {
    border-radius: ${theme.borderRadius.sm}px;
    height: ${theme.heights.input}px;
    font-weight: 500;
  }
`

export const Footer = styled.div`
  text-align: center;
  margin-top: ${theme.spacing.xl}px;
  color: ${theme.colors.textSecondary};
  font-size: ${theme.fontSize.xs}px;
`
