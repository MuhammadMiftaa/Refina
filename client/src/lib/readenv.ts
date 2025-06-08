export const getMode = () => import.meta.env.VITE_MODE;
export const getBackendURL = (): string => {
  const url = import.meta.env.VITE_API_URL;
  if (!url) {
    throw new Error("VITE_API_URL is not defined");
  }
  return url;
};
