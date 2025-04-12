import http from 'k6/http';
import { check } from 'k6';

export let options = {
  scenarios: {
    constant_rps: {
      executor: 'constant-arrival-rate',
      rate: 1000,             
      timeUnit: '1s',
      duration: '30s',
      preAllocatedVUs: 1000,   
      maxVUs: 2000,
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<100'],   
    http_req_failed: ['rate<0.0001'],   
  },
};

export default function () {
  const url = 'http://localhost:8080/pvz';

  const headersModer = {
    'auth-x': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4MDM0Njk4ODAsInJvbGUiOiJtb2RlcmF0b3IifQ.eLWXTE5yg6lhMdsq0t8FHxOj3qmNW-xyrVeZuucpC3g', // ДЛЯ МОДЕРАТОРА НА 16384 ЧАСОВ
  };
  const headersEmpl = {
    'auth-x': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4MDM0NzAwMDgsInJvbGUiOiJlbXBsb3llZSJ9.XxkftNmRjVhUpkHjKU5OmbkVDTnbR1PMhZRrhz5n4YI', // ДЛЯ СОТРУДНИКА НА 16384 ЧАСОВ
  };

  const res = http.get(url, { headers: headersModer });

  check(res, {
    'status is 200': (r) => r.status === 200,
  });
}
