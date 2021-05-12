# go pip package builder


## Building

Run 

```
docker build -t <NAME> .
```


## Running

Run 

```
docker run
    -e AWS_REGION=<REGION>
    -e AWS_ACCESS_KEY_ID=<KEY_ID>
    -e AWS_SECRET_ACCESS_KEY=<KEY>
    -p 8080:8080 <NAME>
```

This service is configured to quit after each request, for isolation needs, so
using `--restart=always` is convenient.
Make sure to properly pass the correct AWS credentials on the run part.

For best results I recommend using a serverless/container orchestration platform, here I used Google Cloud Run.

## Using

Make a HTTP request with the package name on the path:

```
curl localhost:8080/requests
```

If successful the response contains a link to the S3 object, so browsers
download it afterwards, for your convenience.


# Notes

The isolation is achieved by running this program on a serverless/container
platform like Google Cloud Run, with *concurrency* set to 1. This server is
made to kill itself after the first build, hence, the following request will
be handled by a new instance.

To compensate we set the *maximum instances* to a bigger number, like 100, so
this service can bear more activity.