export const addr =
  process.env.NODE_ENV === "development"
    ? "http://localhost:2112/api/v1"
    : "/api/v1";
