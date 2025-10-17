import http from "k6/http";

export const options = {
  stages: [
    { duration: "10s", target: 50 },  // Ramp up ke 50,000 VUs dalam 2 menit
    { duration: "15s", target: 100 }, // Ramp up ke 100,000 VUs dalam 2 menit berikutnya
    { duration: "5s", target: 0 },      // Ramp down ke 0 VUs dalam 1 menit terakhir
  ],
};

export default function () {
  http.get("https://api-refina.miftech.web.id/test");
}