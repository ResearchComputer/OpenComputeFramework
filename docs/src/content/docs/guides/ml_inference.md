---
title: Machine Learning as a Service
description: A guide in my new Starlight docs site.
---

We run several public workers that provide machine learning inference/training as a service. Under the hood each worker may run different backend, to support different types of machine learning models. We currently run the following backend: 
- **Triteia** (for large generative transformer models.
- **DeepSpeed-MII** (for text-to-image and some other models). 
- **Inferencia** (for other HuggingFace models unsupported).

The endpoint for our public workers starts with `https://api.research.computer/`. For example, the endpoint for Triteia is `https://api.research.computer/triteia/`.

## Triteia

Triteia is an inference engine that supports OpenAI-compatible apis for large generative transformer models. Simply replace the endpoint with `https://api.research.computer/triteia/` to use Triteia.

## Inferencia

```python

import requests

response = requests.post(
    url="https://api.research.computer/inferencia/v1/predict",
    json={
        "model_name": "microsoft/deberta-large-mnli",
        "data": [{
            "text": ["You look amazing today,"],
            "top_k": 3,
        }]
    },
)
print(response.json())
```

The expected output is 

```json
{
    'model_name': 'microsoft:deberta-large-mnli', 
    'model_version': 'default', 
    'data': [
        [
            [
                {'label': 'NEUTRAL', 'score': 0.9754309058189392}, {'label': 'CONTRADICTION', 'score': 0.016230667009949684}, {'label': 'ENTAILMENT', 'score': 0.00833841785788536}
            ]
        ]
    ]
}
```