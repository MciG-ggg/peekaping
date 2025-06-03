export const getConfig = () => {
  const config = (window as any).__CONFIG__;
  const isProd = import.meta.env.PROD;

  return {
    API_URL: isProd
      ? config.API_URL // don't fallback to default in prod
      : config.API_URL || "http://localhost:8034",
  };
};
