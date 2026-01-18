/**
 * Centralized route path constants
 * This ensures consistency across routing, navigation, and redirects
 */

export const ROUTES = {
  // Public routes
  ROOT: '/',
  LOGIN: '/login',
  OAUTH_CALLBACK: '/oauth/callback',
  
  // Protected routes
  HOME: '/home',
  BILLING: '/billing',
} as const;

// Type for route values
export type RouteValue = typeof ROUTES[keyof typeof ROUTES];
