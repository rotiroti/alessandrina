import http from "k6/http";
import exec from "k6/execution";
import { check, sleep } from "k6";
import * as settings from "./settings.js";

export const options = {
  ext: settings.cloudRun(`create-book-${settings.workloadName()}`),
  thresholds: {
    http_req_failed: ["rate<0.01"], // http errors should be less than 1%
  },
  scenarios: settings.workload(),
};

export default function () {
  const payload = JSON.stringify(settings.generateRandomPayload());
  const params = {
    headers: { "Content-Type": "application/json" },
    tags: { name: "create-book" },
  };
  const res = http.post(`${settings.BASE_URL}/books`, payload, params);

  check(res, {
    "Post status is 201": (r) => res.status === 201,
    "Post Content-Type header": (r) =>
      res.headers["Content-Type"] === "application/json",
  });

  const isRateScenario =
    `${exec.scenario.executor}` === "constant-arrival-rate" ||
    `${exec.scenario.executor}` === "ramping-arrival-rate";

  if (!isRateScenario) {
    sleep(0.5);
  }
}
