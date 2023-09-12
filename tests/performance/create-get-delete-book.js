import http from "k6/http";
import { check, group, sleep } from "k6";
import * as settings from "./settings.js";

export const options = {
  ext: settings.cloudRun(`create-get-delete-book-${settings.workloadName()}`),
  thresholds: {
    http_req_failed: ["rate<0.01"], // http errors should be less than 1%
  },
  scenarios: settings.workload(),
};

export default function () {
  group("Create, get and delete book", () => {
    let URL = `${settings.BASE_URL}/books`;

    const headers = { "Content-Type": "application/json" };
    const payload = JSON.stringify(settings.generateRandomPayload());
    const params = {
      headers: headers,
      tags: { name: "create-book"},
    };
    const res = http.post(URL, payload, params);
  
    check(res, {
      "Post status is 201": (r) => res.status === 201,
      "Post Content-Type header": (r) =>
        res.headers["Content-Type"] === "application/json",
    });

    if (res.status === 201) {
      const bookId = res.json("id");
      const getRes = http.get(`${URL}/${bookId}`, {
        headers: headers,
        tags: { name: "get-book" },
      });

      check(getRes, {
        "Get status is 200": (r) => getRes.status === 200,
        "Get Content-Type header": (r) =>
        getRes.headers["Content-Type"] === "application/json",
      });

      const delRes = http.del(`${URL}/${bookId}`, null, {
        headers: headers,
        tags: { name: "delete-book" },
      });

      check(delRes, {
        "Delete status is 204": (r) => delRes.status === 204,
      });
    }
  });

  if (!settings.noSleep()) {
    sleep(0.5);
  }
}
