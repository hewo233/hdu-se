```bash
# chat1
curl -X POST 'https://api.coze.cn/v3/chat?' \
-H "Authorization: Bearer cztei_leepfNknGVIJrCyBd0bzL7CKCg0dvHx3bdOFhN642UkSeLbb4zn7zXthJa01MUSCD" \
-H "Content-Type: application/json" \
-d '{
  "bot_id": "7563218003241058343",
  "user_id": "123",
  "stream": false,
  "additional_messages": [
    {
      "role": "user",
      "type": "question",
      "content_type": "text",
      "content": "你是谁？这是我们第几次对话？"
    }
  ]
}'
```
```bash
{"data":{"id":"7588808959469895723","conversation_id":"7588808959469862955","bot_id":"7563218003241058343","created_at":1766907275,"last_error":{"code":0,"msg":""},"status":"in_progress"},"code":0,"msg":""}

cover_id 7588818179242721321
chat_id 7588819419779039272 7588820073197223999 7588820838100008969