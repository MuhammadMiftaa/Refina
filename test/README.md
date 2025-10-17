### Install utility library

npm init -y
npm install k6
npm install --save-dev @types/k6

### Run performance load test

docker exec k6 k6 run /data/scripts/src/script.js
