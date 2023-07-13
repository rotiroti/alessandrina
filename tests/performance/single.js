import http from "k6/http";
import { group, check } from "k6";
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

const BASE_URL = `${__ENV.API_URL}`.replace(/\/$/, "");

export function handleSummary(data) {
  return {
    "single-endpoint.html": htmlReport(data),
  };
}

export const options = {
  discardResponseBodies: true,
  scenarios: {
    contacts: {
      executor: "constant-arrival-rate",

      // How long the test lasts
      duration: "30s",

      // How many iterations per timeUnit
      rate: 1000,

      // Start `rate` iterations per second
      timeUnit: "1h",

      // Pre-allocate VUs
      preAllocatedVUs: 10,
    },
  },
};

export default function () {
  group("get books", () => {
    const res = http.get(`${BASE_URL}/books`);

    check(res, {
      "status is 200": (r) => res.status === 200,
    });
  });

  // group("create book", () => {
  //   const payload = JSON.stringify({
  //     title: "The Go Programming Language",
  //     authors: "Alan A. A. Donovan, Brian W. Kernighan",
  //     publisher: "Addison-Wesley Professional",
  //     pages: 400,
  //     isbs: "978-0134190440",
  //   });
  //   const headers = { "Content-Type": "application/json" };

  //   const res = http.post(`${BASE_URL}/books`, payload, { headers: headers });

  //   check(res, {
  //     "status is 201": (r) => res.status === 201,
  //     "transaction time OK": (r) => r.timings.duration < 200,
  //     "content-type is application/json": (r) => r.headers["Content-Type"] === "application/json",
  //     "json is valid": (r) => JSON.parse(r.body).id !== undefined,
  //   });
  // });

  // group("get book", () => {
  //   http.get(`${BASE_URL}/books/ad8b59c2-5fe6-4321-b0cf-6d2f9eb1c812`);
  // });

  // group("delete book", () => {
  //   http.get(`${BASE_URL}/books/ad8b59c2-5fe6-1234-b0cf-6d2f9eb1c812`);
  // });
}
