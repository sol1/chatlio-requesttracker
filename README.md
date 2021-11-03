## Chatlio request tracker

A listener for chatlio, that creates a new ticket in RT containing the conversation transcript.

### Build
- `go get github.com/LunaNode/rtgo`
- `go install`

### Install
- Create a configuration file
- Run the service
- Point the chatlio callback to the correct address and port

### Required environment variables

 * `RT_URL` e.g. https://rt.example.com
 * `RT_USERNAME`: e.g. live-chat
 * `RT_PASSWORD`: secret
 * `RT_QUEUE`: name of the RT queue

### Run using Docker

```
docker run -p 8080:8080 --name chatlio-rt-test --env RT_URL="https://rt.example.com" --env RT_USERNAME="john" --env RT_PASSWORD="secret" --env RT_QUEUE="general" sol1/chatlio-rt
```

### Simulate a post from Chatlio

```
curl --header "Content-Type: application/json" --request POST --data '{"messages": [{}], "visitorEmail": "user@example.com", "textBody": "Test message body"}' http://localhost:8080/transcript
```