import os
from openai import OpenAI

client = OpenAI(
    base_url="https://api.research.computer/triteia",
    api_key="sk-12345678", # Your API key here
)

stream = client.chat.completions.create(
    messages=[
        {
            "role": "user",
            "content": "hi",
        }
    ],
    model="mistralai/Mistral-7B-Instruct-v0.2",
)
print(stream)