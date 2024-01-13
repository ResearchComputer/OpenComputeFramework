import requests

response = requests.post(
    url="https://api.research.computer/inferencia/v1/predict",
    json={
        "model_name": "microsoft/deberta-large-mnli",
        "data": [{
            "text": ["You look amazing today,", "You look amazing today,", "You look amazing today,", "You look amazing today,"],
            "top_k": 3,
        }]
    },
)
print(response.json())
