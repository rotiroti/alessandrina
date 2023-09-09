import http from "k6/http";
import { check, group, sleep } from "k6";
import { SharedArray } from "k6/data";
import { Counter, Trend } from "k6/metrics";

/**
 * Workload settings and scenarios.
 */
const SMOKE = 0;
const AVERAGE = 1;
const CONSTANT_RATE = 2;

const jsonData = new SharedArray("books", function () {
  const path = "./books.json";

  return JSON.parse(open(path)).books;
});

export const thresholdsSettings = {
  http_req_failed: ["rate<0.01"],
  http_req_duration: ["p(95)<500", "p(99)<1500"],
};

export function threshold() {
  const w = parseInt(__ENV.WORKLOAD);

  switch (w) {
    case (AVERAGE, CONSTANT_RATE):
      return thresholdsSettings;
    default:
      return {};
  }
}

export const smokeWorkload = {
  executor: "shared-iterations",
  vus: 3,
  iterations: 10,
};

export const averageWorkload = {
  executor: "ramping-vus",
  stages: [
    { duration: "5s", target: 6 },
    { duration: "50s", target: 6 },
    { duration: "5s", target: 0 },
  ],
};

export const constantRateWorkload = {
  executor: "constant-arrival-rate",
  duration: "30s",
  preAllocatedVUs: 20, // allocate runtime resources
  rate: 10, // number of constant iterations given `timeUnit`
  timeUnit: "1s",
};

export function workload() {
  const w = parseInt(__ENV.WORKLOAD);

  switch (w) {
    case SMOKE:
      return smokeWorkload;
    case AVERAGE:
      return averageWorkload;
    case CONSTANT_RATE:
      return constantRateWorkload;
    default:
      return smokeWorkload;
  }
}

/**
 * Custom metrics for the Alessandrina API.
 */
const ListErrors = new Counter("ListBooksErrors");
const CreateErrors = new Counter("CreateBookErrors");
const ShowErrors = new Counter("ShowBookErrors");
const DeleteErrors = new Counter("DeleteBookErrors");
const ListTrend = new Trend("ListBooks");
const CreateTrend = new Trend("CreateBook");
const ShowTrend = new Trend("ShowBook");
const DeleteTrend = new Trend("DeleteBook");

/**
 * Test the creation of a book.
 *
 * @param {string} baseUrl - The base URL of the Alessandrina API.
 */
export function create(baseUrl) {
  const randomBook = jsonData[Math.floor(Math.random() * jsonData.length)];
  const payload = {
    title: randomBook.title,
    authors: randomBook.authors,
    publisher: randomBook.publisher,
    pages: parseInt(randomBook.pages),
    isbn: randomBook.isbn,
  };

  const params = { headers: { "Content-Type": "application/json" } };
  const res = http.post(`${baseUrl}/books`, JSON.stringify(payload), params);

  check(res, { "Book created correctly": (r) => r.status === 201 }) ||
    CreateErrors.add(1);
  CreateTrend.add(res.timings.duration);

  sleep(0.5);
}

/**
 * Test the retrieval of a list of books.
 *
 * @param {string} baseUrl - The base URL of the Alessandrina API.
 */
export function list(baseUrl) {
  const res = http.get(`${baseUrl}/books`);

  check(res, { "Retrieve list of books": (r) => r.status === 200 }) ||
    ListErrors.add(1);
  ListTrend.add(res.timings.duration);
  sleep(0.5);
}

/**
 * Test the creation, retrieval and deletion of a book.
 *
 * @param {string} baseUrl - The base URL of the Alessandrina API.
 */
export function flow(baseUrl) {
  group("Create, show and delete book flow", function () {
    const randomBook = jsonData[Math.floor(Math.random() * jsonData.length)];
    const payload = {
      title: randomBook.title,
      authors: randomBook.authors,
      publisher: randomBook.publisher,
      pages: parseInt(randomBook.pages),
      isbn: randomBook.isbn,
    };
    const params = { headers: { "Content-Type": "application/json" } };
    const res = http.post(`${baseUrl}/books`, JSON.stringify(payload), params);

    check(res, { "Book created correctly": (r) => r.status === 201 }) ||
      CreateErrors.add(1);
    CreateTrend.add(res.timings.duration);

    sleep(0.5);

    const getRes = http.get(`${baseUrl}/books/${res.json("id")}`);

    check(getRes, { "Show book": (r) => r.status === 200 }) ||
      ShowErrors.add(1);
    ShowTrend.add(getRes.timings.duration);

    sleep(0.5);

    const delRes = http.del(`${baseUrl}/books/${res.json("id")}`, null, params);
    check(delRes, { "Book was deleted correctly": (r) => r.status === 204 }) ||
      DeleteErrors.add(1);
    DeleteTrend.add(delRes.timings.duration);

    sleep(0.5);
  });
}
