curl "http://localhost:8080/v1/completions

{
    "model": "gpt-4o-mini",
    "prompt": "从前有个男人叫小帅",
    "max_tokens": 20,
    "temperature": 0.8,
    "stream":true
  }

curl -X POST http://localhost:8080/v1/chat/completions \

{
    "model": "gpt-4o-mini",
    "stream":true,
    "messages": [
      {
        "role": "user",
        "content": "讲述一个十个字以内的笑话"
      }
    ]
  }