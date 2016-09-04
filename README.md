# S3 Object Proxy
This provides a simple method of proxying calls to S3 with logic determining if an object should be accessible or not.

For example, you may have a private bucket that you'd like to serve static resources out of, but need to prevent
unpublished content from being served. Using Redis and the HTTP example below, you can do something like this:

```redis
SET wFovN4O.mp4 2
```

Then retrieve the data:

```
curl http://localhost:8080/proxy/wFovN4O.mp4
```

The longer-term goal is to migrate this to a Lambda function behind API Gateway that returns a bytestream. At this time, it seems that is unsupported. 

## HTTP
Run vanilla with:
```
# Build without vars
go build -o proxy examples/http/main.go
./proxy -bucket=bucket.name -region=bucket.region -redis=127.0.0.1:6379

# Build vars in
go build -o proxy -ldflags "-X main.bucket=bucket.name -X main.region=us-west-2 -X main.redisAddr="localhost:6379"" examples/http/main.go
./proxy
```
Or take a look at `http-example/main.go`.

## Lambda Function
```
cd examples/lambda
BUCKET=some.bucket REGION=us-west-2 REDIS_HOST=some.redis.host ./build.sh
``` 

### LDFlags
Currently using LDFlags as a clean way to build differently configured binaries. Will continue to explore other options.

### Limitations
- Lambda functions must return strings, so the result is Base64 encoded
- API Gateway imposes JSON, so your output will be:
```json
{
    "value": "base64-encoded-data"
}
```