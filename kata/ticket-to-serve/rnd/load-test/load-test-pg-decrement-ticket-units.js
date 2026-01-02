import http from "k6/http";

export const options = {
  vus: 1000,
  iterations: 100000,
};

const BASE_URL = "http://localhost:8080";
const TICKET_IDS = ["a", "b", "c", "d", "e", "f", "g", "h", "i", "j"];

function getRandomTicketId() {
  return TICKET_IDS[Math.floor(Math.random() * TICKET_IDS.length)];
}

function getRandomCount() {
  return Math.floor(Math.random() * 3) + 1; // 1 to 3
}

function generateTickets() {
  const numTickets = Math.floor(Math.random() * 3) + 1; // 1 to 3 ticket types
  const tickets = [];

  for (let i = 0; i < numTickets; i++) {
    tickets.push({
      id: getRandomTicketId(),
      count: getRandomCount(),
    });
  }

  return tickets;
}

export default function () {
  const payload = JSON.stringify({
    tickets: generateTickets(),
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  http.post(`${BASE_URL}/pg-decrement-ticket-units`, payload, params);
}
