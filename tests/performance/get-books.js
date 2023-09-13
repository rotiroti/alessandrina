import http from "k6/http";
import exec from "k6/execution";
import { check, sleep } from "k6";
import * as settings from "./settings.js";

export const options = {
  ext: settings.cloudRun(`get-books-${settings.workloadName()}`),
  thresholds: {
    http_req_failed: ["rate<0.01"], // http errors should be less than 1%
  },
  scenarios: settings.workload(),
};

export default function () {
  const res = http.get(`${settings.BASE_URL}/books`, {
    tags: { name: "get-books" },
  });

  check(res, {
    "Get status is 200": (r) => res.status === 200,
    "Get Content-Type header": (r) =>
      res.headers["Content-Type"] === "application/json",
  });

  if (
    `${exec.scenario.executor}` !== "constant-arrival-rate" &&
    `${exec.scenario.executor}` !== "ramping-arrival-rate"
  ) {
    sleep(0.5);
  }
}
