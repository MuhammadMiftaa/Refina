/**
 * Runtime Environment Configuration Template
 * 
 * File ini adalah TEMPLATE yang akan diproses oleh envsubst
 * pada saat container startup.
 * 
 * Placeholder $API_URL akan diganti dengan nilai dari environment variable.
 * 
 * JANGAN edit file ini di production - ini hanya template!
 */
window.__ENV__ = {
  API_URL: "$API_URL",
  MODE: "$MODE",
  VITE_API_BASE_URL: "$VITE_API_BASE_URL"
};