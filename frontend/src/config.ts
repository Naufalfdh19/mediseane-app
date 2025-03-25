const cfg = {
  API_BASE_URL: import.meta.env.VITE_BASE_URL,
  ROUTE_BASE_PATH: import.meta.env.VITE_BASE_PATH,
};

const testCfg = {
  API_BASE_URL: "http://__SERVER_URL__",
  ROUTE_BASE_PATH: import.meta.env.VITE_BASE_PATH,
};

export default import.meta.env.MODE === "test" ? testCfg : cfg;
