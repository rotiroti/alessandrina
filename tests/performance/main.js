import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import { create, list, flow, threshold, workload } from "./config.js";

const BASE_URL = `${__ENV.API_URL}`.replace(/\/$/, "");

export function handleSummary(data) {
  return {
    "general.html": htmlReport(data),
  };
}

export const options = {
  scenarios: {
    scenario: workload(),
  },
  thresholds: threshold(),
  ext: {
    loadimpact: {
      projectID: 0,
      name: __ENV.TEST_NAME || "main.js",
      distribution: {
        "us-east-1": { loadZone: "amazon:us:ashburn", percent: 100 },
      },
    },
  },
};

export default function () {
  switch (__ENV.BOOK_OP) {
    case "create":
      create(BASE_URL);
      break;
    case "list":
      list(BASE_URL);
      break;
    case "flow":
      flow(BASE_URL);
      break;
    default:
      return;
  }
}
