import http from "k6/http";

export const options = {
  vus: 10000,
  iterations: 1000000,
};

export default function () {
  const url = "http://localhost:8080/complex-redis-decrement";
  const payload = JSON.stringify({
    tickets: [
      { id: "a", count: 2 },
      { id: "b", count: 2 },
    ],
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  http.post(url, payload, params);
}
