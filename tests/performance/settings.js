import { Trend } from "k6/metrics";
import { SharedArray } from "k6/data";

// Create custom trend metrics
export const createBookLatency = new Trend("create_book_duration");
export const getBookLatency = new Trend("get_book_duration");
export const deleteBookLatency = new Trend("delete_book_duration");

/**
 * Base URL of the Alessandrina API.
 *
 * @type {string}
 */
export const BASE_URL = `${__ENV.API_URL}`.replace(/\/$/, "");

/**
 * Shared array of books read from a JSON file.
 *
 * @type {SharedArray}
 */
const BOOKS_DATA = new SharedArray("books", function () {
  const path = "./books.json";
  const data = JSON.parse(open(path)).books;

  return data;
});

/**
 * Generate a random payload for a new book.
 *
 * @returns {object}
 */
export const generateRandomPayload = () => {
  const randomBook = BOOKS_DATA[Math.floor(Math.random() * BOOKS_DATA.length)];
  const payload = {
    title: randomBook.title,
    authors: randomBook.authors,
    publisher: randomBook.publisher,
    pages: randomBook.pages,
    isbn: randomBook.isbn,
  };

  return payload;
};

/**
 * Configure the cloud execution.
 *
 * @param {*} testName
 * @returns
 */
export const cloudRun = (testName) => {
  return {
    loadimpact: {
      projectID: 0,
      name: testName,
      distribution: {
        "us-east-1": { loadZone: "amazon:us:ashburn", percent: 100 },
      },
    },
  };
};

// 3 VUs * 10 iterations = 30 requests
const baseline = {
  baseline: {
    executor: "per-vu-iterations",
    vus: 3,
    iterations: 10,
    maxDuration: "30s",
  },
};

// 5 VUs * 200 iterations = 1000 requests
const vus5 = {
  vus5: {
    executor: "per-vu-iterations",
    vus: 5,
    iterations: 200,
  },
};

// 10 VUs * 200 iterations = 2000 requests
const vus10 = {
  vus10: {
    executor: "per-vu-iterations",
    vus: 10,
    iterations: 200,
  },
};

// 15 VUs * 200 iterations = 3000 requests
const vus15 = {
  vus15: {
    executor: "per-vu-iterations",
    vus: 15,
    iterations: 200,
  },
};

const averageVUs = {
  averageVUs: {
    executor: "ramping-vus",
    startVUs: 0,
    stages: [
      { duration: "30s", target: 10 },
      { duration: "4m", target: 10 },
      { duration: "30s", target: 0 },
    ],
  },
};

// ~8K of total requests
const stressVUs = {
  stressVUs: {
    executor: "ramping-vus",
    startVUs: 0,
    stages: [
      { duration: "30s", target: 20 },
      { duration: "4m", target: 20 },
      { duration: "30s", target: 0 },
    ],
  },
};

/**
 * Workload configuration.
 */
const workloads = {
  0: baseline,
  1: vus5,
  2: vus10,
  3: vus15,
  4: averageVUs,
  5: stressVUs,
};

/**
 * Workload names.
 */
const workloadNames = {
  0: "baseline",
  1: "vus5",
  2: "vus10",
  3: "vus15",
  4: "averageVUs",
  5: "stressVUs",
};

/**
 * Return the workload configuration.
 *
 * @returns {object}
 */
export const workload = () => {
  const idx = parseInt(__ENV.WORKLOAD || 0);

  return workloads[idx];
};

/**
 * Return the workload name.
 *
 * @returns {string}
 */
export const workloadName = () => {
  const idx = parseInt(__ENV.WORKLOAD || 0);

  return workloadNames[idx];
};
