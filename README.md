## Chatlio request tracker

A listener for chatlio, that creates a new ticket in RT containing the conversation transcript.

### Build
- `go get github.com/LunaNode/rtgo`
- `go install`

### Install
- Create a configuration file
- Run the service
- Point the chatlio callback to the correct address and port

### Example configuration file
```
{
	"URL" : "rt.example.com",
	"Username" : "chatlio-user",
	"Password" : "super-s3cure",
	"Queue" : "Chatlio transcripts"
}
```