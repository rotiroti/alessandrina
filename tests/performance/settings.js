import { SharedArray } from "k6/data";

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
 * Check if the current scenario does not require a sleep.
 */
export const noSleep = () => {
  const s = __ENV.NO_SLEEP || "false";

  return s.toLowerCase() === "true";
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

const baseline = {
  baseline: {
    executor: "per-vu-iterations",
    vus: 3,
    iterations: 10,
  },
};

const rate5 = {
  rate5: {
    executor: "constant-arrival-rate",
    rate: 5,
    timeUnit: "1s",
    duration: "5m",
    preAllocatedVUs: 15,
  },
};

const rate10 = {
  rate10: {
    executor: "constant-arrival-rate",
    rate: 10,
    timeUnit: "1s",
    duration: "5m",
    preAllocatedVUs: 15,
  },
};

const rate25 = {
  rate25: {
    executor: "constant-arrival-rate",
    rate: 25,
    timeUnit: "1s",
    duration: "5m",
    preAllocatedVUs: 30,
  },
};

const stressRate = {
  stressRate: {
    executor: "ramping-arrival-rate",
    startRate: 0,
    timeUnit: "1s",
    stages: [
      { duration: "1m", target: 5 },
      { duration: "2m30s", target: 5 },
      { duration: "30s", target: 10 },
      { duration: "2m30s", target: 10 },
      { duration: "30s", target: 25 },
      { duration: "2m30s", target: 25 },
      { duration: "30s", target: 0 },
    ],
  },
};

/**
 * Workload configuration.
 */
const workloads = {
  0: baseline,
  1: rate5,
  2: rate10,
  3: rate25,
  4: stressRate,
};

/**
 * Workload names.
 */
const workloadNames = {
  0: "baseline",
  1: "rate5",
  2: "rate10",
  3: "rate25",
  4: "stressRate",
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
