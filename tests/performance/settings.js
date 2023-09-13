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

const baseline = {
  baseline: {
    executor: "per-vu-iterations",
    vus: 3,
    iterations: 10,
    maxDuration: "30s",
  },
};

const vus5 = {
  vus5: {
    executor: "constant-vus",
    vus: 5,
    duration: "1m",
  },
};

const vus10 = {
  vus10: {
    executor: "constant-vus",
    vus: 10,
    duration: "1m",
  },
};

const vus15 = {
  vus15: {
    executor: "constant-vus",
    vus: 15,
    duration: "1m",
  },
};

const averageVUs = {
  averageVUs: {
    executor: "ramping-vus",
    startVUs: 0,
    stages: [
      { duration: "10s", target: 10 },
      { duration: "1m40s", target: 10 },
      { duration: "10s", target: 0 },
    ],
  },
};

const stressVUs = {
  stressVUs: {
    executor: "ramping-vus",
    startVUs: 0,
    stages: [
      { duration: "1m", target: 5 },
      { duration: "2m30s", target: 5 },
      { duration: "30s", target: 10 },
      { duration: "2m30s", target: 10 },
      { duration: "30s", target: 15 },
      { duration: "2m30s", target: 15 },
      { duration: "30s", target: 0 },
    ],
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
    preAllocatedVUs: 30,
  },
};

const rate25 = {
  rate25: {
    executor: "constant-arrival-rate",
    rate: 25,
    timeUnit: "1s",
    duration: "5m",
    preAllocatedVUs: 50,
  },
};

const averageRate = {
  averageRate: {
    executor: "ramping-arrival-rate",
    startRate: 0,
    timeUnit: "1s",
    preAllocatedVUs: 30,
    stages: [
      { duration: "1m", target: 10 },
      { duration: "3m30s", target: 10 },
      { duration: "30s", target: 0 },
    ],
  },
};

const stressRate = {
  stressRate: {
    executor: "ramping-arrival-rate",
    startRate: 0,
    timeUnit: "1s",
    preAllocatedVUs: 50,
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
  1: vus5,
  2: vus10,
  3: vus15,
  4: averageVUs,
  5: stressVUs,
  6: rate5,
  7: rate10,
  8: rate25,
  9: averageRate,
  10: stressRate,
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
  6: "rate5",
  7: "rate10",
  8: "rate25",
  9: "averageRate",
  10: "stressRate",
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
