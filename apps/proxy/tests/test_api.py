import openai

oai = openai.OpenAI(
    api_key="sk-xxxx",
    base_url="http://localhost:8000/v1",
)

resp = oai.chat.completions.create(
    model="Qwen/Qwen3-8B",
    messages=[
        {"role": "user", "content": "Hello!"},
    ],
)
print(resp.choices[0].message.content)