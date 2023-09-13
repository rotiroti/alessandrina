import http from "k6/http";
import exec from "k6/execution";
import { check, sleep } from "k6";
import * as settings from "./settings.js";

export const options = {
  ext: settings.cloudRun(`create-delete-book-${settings.workloadName()}`),
  thresholds: {
    http_req_failed: ["rate<0.01"], // http errors should be less than 1%
  },
  scenarios: settings.workload(),
};

export default function () {
  const headers = { "Content-Type": "application/json" };
  const payload = JSON.stringify(settings.generateRandomPayload());
  const params = {
    headers: headers,
    tags: {
      name: "create-book",
    },
  };
  const res = http.post(`${settings.BASE_URL}/books`, payload, params);

  settings.createBookLatency.add(res.timings.duration);

  check(res, {
    "Post status is 201": (r) => res.status === 201,
    "Post Content-Type header": (r) =>
      res.headers["Content-Type"] === "application/json",
  });

  if (res.status === 201) {
    const bookId = res.json("id");
    const delRes = http.get(`${settings.BASE_URL}/books/${bookId}`, {
      headers: headers,
      tags: { name: "delete-book" },
    });

    settings.deleteBookLatency.add(delRes.timings.duration);

    check(delRes, {
      "Delete status is 200": (r) => delRes.status === 200,
    });
  }

  const isRateScenario =
    `${exec.scenario.executor}` === "constant-arrival-rate" ||
    `${exec.scenario.executor}` === "ramping-arrival-rate";

  if (!isRateScenario) {
    sleep(0.5);
  }
}
