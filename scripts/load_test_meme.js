import http from 'k6/http';
import { check, sleep } from 'k6';

const baseUrl = 'http://localhost:8080';
const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjc3MzY2Njc2YTYyYjhiMjg4ZjljOWFmIiwiZW1haWwiOiJub25sbG1AZXhhbXBsZS5jb20iLCJwZXJtaXNzaW9ucyI6eyJHZW5lcmF0ZUxsbU1lbWUiOmZhbHNlfSwiZXhwIjoxNzM1NzU1NTczLCJpYXQiOjE3MzU2NjkxNzN9.MU8CHNYY356RfyoTBkGvM4fdk8MDGY2YsaMOLJoUXI0';

export const options = {
  vus: 200, // Virtual users
  duration: '15s', // Test duration
};

export default function () {
  const payload = JSON.stringify({
    lat: 40.730610,
    lng: -73.935242,
    query: 'generate an image with 1000 shoes'
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    }
  };

  const res = http.post(`${baseUrl}/v1/meme`, payload, params);

  console.log(res.status)
  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(1);
}