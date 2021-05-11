# go pip package builder


## Building

Run 

```
docker build -t <NAME> .
```


## Running

Run 

```
docker run -e AWS_REGION=<REGION> -e AWS_ACCESS_KEY_ID=<KEY_ID> -e AWS_SECRET_ACCESS_KEY=<KEY> -p 8080:8080 <NAME>
```

This service is configured to quit after each request, for isolation needs, so
using `--restart=always` is convenient.
Make sure to properly pass the correct AWS credentials on the run part.


## Using

Make a HTTP request with the package name on the path:

```
curl localhost:8080/requests
```

If successful the response is a redirect to the S3 object, so browsers
download it afterwards, for your convenience.
