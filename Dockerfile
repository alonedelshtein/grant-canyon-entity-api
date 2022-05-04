FROM golang:latest as builder

# Move to working directory /build
WORKDIR /app

# Copy 
ADD . /app/

# build the db-api service
RUN CGO_ENABLED=0 go build -o entity-api

# build the fianl stage in a multi stage docker file
FROM opensearchproject/logstash-oss-with-opensearch-output-plugin:latest
WORKDIR /app
COPY --from=builder /app/app.env .
COPY --from=builder /app/entity-api .
COPY --from=builder /app/entity-api-logstash.conf .
COPY --from=builder /app/deploy.sh .
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /zoneinfo.zip

ENV ZONEINFO=/zoneinfo.zip

#RUN nohup logstash -f db-api-logstash.conf & 

# Command to run when starting the container
CMD ["/bin/bash","./deploy.sh"]