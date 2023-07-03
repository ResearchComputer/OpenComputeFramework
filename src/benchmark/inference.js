import http from 'k6/http';

export const options = {
    stages: [
      { duration: '30s', target: 20 },
      { duration: '1m30s', target: 10 },
      { duration: '20s', target: 0 },
    ],
  };

export default function () {
    const url = 'https://inference.autoai.dev/api/v1/request/inference';
    const payload = JSON.stringify({
        model_name: 'mosaicml/mpt-7b-chat',
        params: {
            'prompt': "I'm feeling happy today",
        }
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    http.post(url, payload, params);
}