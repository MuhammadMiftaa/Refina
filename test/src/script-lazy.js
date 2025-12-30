import http from "k6/http";
import { check, sleep } from "k6";
import { Rate, Trend, Counter } from "k6/metrics";
import { SharedArray } from "k6/data";

// Custom metrics
const errorRate = new Rate("errors");
const uploadDuration = new Trend("upload_duration");
const successfulUploads = new Counter("successful_uploads");
const failedUploads = new Counter("failed_uploads");

// Configuration
const BASE_URL =
  "https://api-refina-exp-lazy.miftech.web.id/v1/transactions/income";
// const BASE_URL = 'http://69.62.80.249:10001/v1/transactions/income';
const AUTH_TOKEN =
  "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6Im1hZG1pZnRhNzdAZ21haWwuY29tIiwiZXhwIjoxNzY2Nzk5ODgxLCJpZCI6ImZkZDIxNDI0LWRlY2EtNGY2Ni04YzNhLWM4NDdmZDJjYzA4NiIsInVzZXJuYW1lIjoiTXVoYW1tYWQgTWlmdGEifQ.JSkaRE8OPTZkEjdYTPbErCZTyLWihncapdXu5M1nCYI";

// Load Base64 files - dibaca sekali saat init, shared across VUs
const base64Files = {
  "5KB": open("../assets/base64_5kb.txt"),
  "10KB": open("../assets/base64_10kb.txt"),
  "20KB": open("../assets/base64_20kb.txt"),
  "1MB": open("../assets/base64_1mb.txt"),
  "5MB": open("../assets/base64_5mb.txt"),
  "10MB": open("../assets/base64_10mb.txt"),
  "15MB": open("../assets/base64_15mb.txt"),
};

// Test scenarios - pilih salah satu dengan uncomment
export const options = {
  // ========== SCENARIO 1: Warm-up Load ==========
  //  stages: [
  //    { duration: '1m', target: 5 },
  //    { duration: '1m', target: 5 },
  //  ],
  // thresholds: {
  //   http_req_duration: ['p(95)<3000'],
  //   errors: ['rate<0.1'],
  // },

  // ========== SCENARIO 2: Normal Load - Small (1MB) ==========
  stages: [
    { duration: "1m", target: 20 },
    { duration: "3m", target: 50 },
    { duration: "1m", target: 20 },
  ],
  // thresholds: {
  //   http_req_duration: ['p(95)<5000', 'p(99)<8000'],
  //   errors: ['rate<0.05'],
  // },

  // ========== SCENARIO 3: Normal Load - Medium (5MB) ==========
  // stages: [
  //   { duration: '1m', target: 20 },
  //   { duration: '3m', target: 50 },
  //   { duration: '1m', target: 20 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<8000', 'p(99)<12000'],
  //   errors: ['rate<0.05'],
  // },

  // ========== SCENARIO 4: Normal Load - Large (10MB) ==========
  // stages: [
  //   { duration: '1m', target: 20 },
  //   { duration: '3m', target: 50 },
  //   { duration: '1m', target: 20 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<12000', 'p(99)<18000'],
  //   errors: ['rate<0.08'],
  // },

  // ========== SCENARIO 5: High Load - Small (1MB) ==========
  // stages: [
  //   { duration: '2m', target: 50 },
  //   { duration: '3m', target: 150 },
  //   { duration: '2m', target: 50 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<6000', 'p(99)<10000'],
  //   errors: ['rate<0.1'],
  // },

  // ========== SCENARIO 6: High Load - Medium (5MB) ==========
  // stages: [
  //   { duration: '2m', target: 50 },
  //   { duration: '3m', target: 150 },
  //   { duration: '2m', target: 50 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<10000', 'p(99)<15000'],
  //   errors: ['rate<0.1'],
  // },

  // ========== SCENARIO 7: High Load - Large (15MB) ==========
  // stages: [
  //   { duration: '2m', target: 50 },
  //   { duration: '3m', target: 150 },
  //   { duration: '2m', target: 50 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<15000', 'p(99)<22000'],
  //   errors: ['rate<0.15'],
  // },

  // ========== SCENARIO 8: Sustained - Medium (5MB) ==========
  // stages: [
  //   { duration: '2m', target: 100 },
  //   { duration: '8m', target: 100 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<8000', 'p(99)<12000'],
  //   errors: ['rate<0.08'],
  // },

  // ========== SCENARIO 9: Sustained - Large (15MB) ==========
  // stages: [
  //   { duration: '2m', target: 120 },
  //   { duration: '8m', target: 120 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<12000', 'p(99)<18000'],
  //   errors: ['rate<0.1'],
  // },

  // ========== SCENARIO 10: Spike Test - Medium (5MB) ==========
  // stages: [
  //   { duration: '1m', target: 30 },
  //   { duration: '2m', target: 250 },
  //   { duration: '3m', target: 250 },
  //   { duration: '2m', target: 30 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<15000', 'p(99)<25000'],
  //   errors: ['rate<0.2'],
  // },

  // ========== SCENARIO 11: Spike Test - Large (15MB) ==========
  // stages: [
  //   { duration: '1m', target: 30 },
  //   { duration: '2m', target: 300 },
  //   { duration: '4m', target: 300 },
  //   { duration: '3m', target: 30 },
  // ],
  // thresholds: {
  //   http_req_duration: ['p(95)<20000', 'p(99)<30000'],
  //   errors: ['rate<0.25'],
  // },

  // ========== SCENARIO 12: Stress Test - Extreme (10MB) ==========
  //   stages: [
  //     { duration: '2m', target: 100 },
  //     { duration: '3m', target: 200 },
  //     { duration: '3m', target: 300 },
  //     { duration: '2m', target: 400 },
  //     { duration: '2m', target: 100 },
  //   ],
  //   thresholds: {
  //     http_req_duration: ['p(95)<25000'],
  //     errors: ['rate<0.3'],
  //   },
};

// Generate test payload
function generatePayload(testCaseNumber, scenarioName, fileSize) {
  return JSON.stringify({
    amount: 100000,
    wallet_id: "b88fe81d-126d-430f-bd6b-4eee2a666111",
    category_id: "2718e35b-ec13-4264-8c69-bd0eb5adabe6",
    date: new Date().toISOString(),
    description: `[TC-${testCaseNumber}] ${scenarioName} - VU:${__VU} - Iter:${__ITER}`,
    attachments: [
      {
        status: "create",
        files: [`data:image/png;base64,${fileSize}`],
      },
    ],
    from_wallet_id: "",
    to_wallet_id: "",
    admin_fee: 0,
  });
}

export default function () {
  // $ GANTI parameter sesuai test case yang aktif
  const testCaseNumber = 5; // Nomor test case
  const scenarioName = "High Load - Small (1MB)"; // Nama scenario
  const fileSize = "5KB"; // Ukuran file: 1MB, 5MB, 10MB, 15MB

  const payload = generatePayload(
    testCaseNumber,
    scenarioName,
    base64Files[fileSize]
  );

  const params = {
    headers: {
      "Content-Type": "application/json",
      Authorization: AUTH_TOKEN,
    },
    timeout: "300s", //~ Timeout
  };

  //   const startTime = new Date();

  // Execute request
  http.post(BASE_URL, payload, params);

  //   const endTime = new Date();
  //   const duration = endTime - startTime;

  // Track metrics
  //   uploadDuration.add(duration);

  //   const success = check(response, {
  //     'status is 200 or 201': (r) => r.status === 200 || r.status === 201,
  //     'response time < 30s': (r) => r.timings.duration < 30000,
  //     'no server errors': (r) => r.status < 500,
  //   });

  //   if (success) {
  //     successfulUploads.add(1);
  //   } else {
  //     failedUploads.add(1);
  //     errorRate.add(1);
  //     console.error(`[ERROR] VU:${__VU} Iter:${__ITER} - Status:${response.status} - ${response.body}`);
  //   }

  //   errorRate.add(!success);

  // Think time: simulate realistic user behavior
  sleep(Math.random() * 2 + 1); // Random 1-3 detik
}

// Summary handler
// export function handleSummary(data) {
//   return {
//     stdout: textSummary(data, { indent: ' ', enableColors: true }),
//     'summary.json': JSON.stringify(data),
//   };
// }

// function textSummary(data, options) {
//   const indent = options.indent || '';
//   const summary = [];

//   summary.push(`${indent}========== K6 LOAD TEST SUMMARY ==========`);
//   summary.push(`${indent}Test Duration: ${data.state.testRunDurationMs / 1000}s`);
//   summary.push(`${indent}VUs: ${data.metrics.vus?.values.max || 'N/A'}`);
//   summary.push('');

//   summary.push(`${indent}📊 HTTP Metrics:`);
//   summary.push(
//     `${indent}  Total Requests: ${data.metrics.http_reqs?.values.count || 0}`
//   );
//   summary.push(
//     `${indent}  Success Rate: ${(
//       ((data.metrics.http_reqs?.values.count || 0) -
//         (data.metrics.failed_uploads?.values.count || 0)) /
//         (data.metrics.http_reqs?.values.count || 1) *
//       100
//     ).toFixed(2)}%`
//   );
//   summary.push(
//     `${indent}  Error Rate: ${(
//       (data.metrics.errors?.values.rate || 0) * 100
//     ).toFixed(2)}%`
//   );
//   summary.push('');

//   summary.push(`${indent}⚡ Latency (ms):`);
//   summary.push(
//     `${indent}  Avg: ${(data.metrics.http_req_duration?.values.avg || 0).toFixed(2)}`
//   );
//   summary.push(
//     `${indent}  Min: ${(data.metrics.http_req_duration?.values.min || 0).toFixed(2)}`
//   );
//   summary.push(
//     `${indent}  Max: ${(data.metrics.http_req_duration?.values.max || 0).toFixed(2)}`
//   );
//   summary.push(
//     `${indent}  P50: ${(data.metrics.http_req_duration?.values['p(50)'] || 0).toFixed(2)}`
//   );
//   summary.push(
//     `${indent}  P95: ${(data.metrics.http_req_duration?.values['p(95)'] || 0).toFixed(2)}`
//   );
//   summary.push(
//     `${indent}  P99: ${(data.metrics.http_req_duration?.values['p(99)'] || 0).toFixed(2)}`
//   );
//   summary.push('');

//   summary.push(`${indent}🚀 Throughput:`);
//   summary.push(
//     `${indent}  Requests/sec: ${(
//       (data.metrics.http_reqs?.values.rate || 0)
//     ).toFixed(2)}`
//   );
//   summary.push(
//     `${indent}  Data Received: ${(
//       (data.metrics.data_received?.values.count || 0) /
//       1024 /
//       1024
//     ).toFixed(2)} MB`
//   );
//   summary.push(
//     `${indent}  Data Sent: ${(
//       (data.metrics.data_sent?.values.count || 0) /
//       1024 /
//       1024
//     ).toFixed(2)} MB`
//   );

//   summary.push(`${indent}==========================================`);

//   return summary.join('\n');
// }
