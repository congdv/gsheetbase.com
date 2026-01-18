// Global theme configuration
export const theme = {
  colors: {
    primary: '#1a237e',
    primaryLight: '#5b6abf',
    background: 'linear-gradient(135deg, #2e3b82 0%, #1a237e 100%)',
    cardBackground: '#ffffff',
    inputBackground: '#f5f5f5',
    inputBorder: '#e0e0e0',
    textPrimary: '#000000',
    textSecondary: '#999999',
    textError: '#ff4d4f',
  },
  spacing: {
    xs: 8,
    sm: 16,
    md: 24,
    lg: 32,
    xl: 40,
  },
  borderRadius: {
    sm: 4,
    md: 8,
  },
  fontSize: {
    xs: 12,
    sm: 14,
    md: 16,
    lg: 32,
  },
  heights: {
    input: 44,
  },
}

export type Theme = typeof theme
