# syntax=docker/dockerfile:1

FROM debian:bullseye-slim

WORKDIR /
RUN mkdir /app
COPY ./app /app

EXPOSE 3006
RUN ["chmod", "+x", "/app/calc_lag_monitor"]
ENTRYPOINT ["/app/calc_lag_monitor"]