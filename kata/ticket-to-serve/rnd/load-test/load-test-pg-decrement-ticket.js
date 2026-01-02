import http from "k6/http";

export const options = {
  vus: 1000,
  iterations: 100000,
};

const BASE_URL = "http://localhost:8080";

export default function () {
  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  http.post(`${BASE_URL}/pg-decrement-ticket`, null, params);
}
