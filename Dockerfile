FROM docker.io/golang:1.20-alpine AS builder

RUN mkdir /app && \
  mkdir /build

ADD api/ /app

WORKDIR /app

RUN go build -o /build/argo-cd-exporter main.go

FROM scratch as app

COPY --from=builder /build/argo-cd-exporter /argo-cd-exporter

EXPOSE 8080 

ENTRYPOINT ["/argo-cd-exporter"]