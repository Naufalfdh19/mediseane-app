import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 100 },
    { duration: '30s', target: 100 },
    { duration: '10s', target: 0 }, 
  ],
  thresholds: {
    http_req_duration: ["p(95)<500"],
    checks: ["rate>0.99"]
  }
};

export default function () {
  let res = http.get('http://localhost:8080/api/v1/general-products?name=v');

  check(res, { 'success': (r) => r.status === 200 });

  sleep(0.1);
}