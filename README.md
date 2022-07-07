# Simple stress tester in go to post json payload from file to API endpoint with Bearer authentication


```cmd
> .\cmd\windows\client.exe BEARER_TOKEN FULL_URL PATH_TO_PAYLAOD NUMBER_OF_CONNECTIONS
```

example:
```cmd
> .\cmd\windows\client.exe eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1 http://localhost:8080/api/v2/thing .\request.json 100000
``` 
