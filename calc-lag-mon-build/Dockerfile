# syntax=docker/dockerfile:1

FROM debian:bookworm-slim
RUN apt-get update
RUN apt-get install -y ca-certificates
RUN update-ca-certificates -v

WORKDIR /
RUN mkdir /app
COPY ./app /app

EXPOSE 3006

WORKDIR /app/calc_lag_monitor
ENTRYPOINT ["/app/calc_lag_monitor/compiled/calc_lag_monitor"]