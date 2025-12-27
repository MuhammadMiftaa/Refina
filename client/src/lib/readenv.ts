/**
 * Global type declaration for runtime environment variables
 * These are injected at container startup, not build time
 */
declare global {
  interface Window {
    __ENV__: {
      API_URL: string;
      MODE: string;
      VITE_API_BASE_URL: string;
    };
  }
}

export {};

export const getMode = (): string =>
  typeof window !== "undefined" && window.__ENV__?.MODE
    ? window.__ENV__.MODE
    : import.meta.env.VITE_MODE || "development";

export const getBackendURL = (): string =>
  typeof window !== "undefined" && window.__ENV__?.API_URL
    ? window.__ENV__.API_URL
    : import.meta.env.VITE_API_URL || "http://localhost:8080/v1";

export const getViteAPIBaseURL = (): string =>
  typeof window !== "undefined" && window.__ENV__?.VITE_API_BASE_URL
    ? window.__ENV__.VITE_API_BASE_URL
    : import.meta.env.VITE_API_BASE_URL || "http://localhost:5173";
 