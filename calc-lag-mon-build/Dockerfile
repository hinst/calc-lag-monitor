# syntax=docker/dockerfile:1

FROM debian:bullseye-slim

WORKDIR /
RUN mkdir /app
COPY ./app /app

EXPOSE 3006

WORKDIR /app/calc_lag_monitor
ENTRYPOINT ["/app/calc_lag_monitor/compiled/calc_lag_monitor"]