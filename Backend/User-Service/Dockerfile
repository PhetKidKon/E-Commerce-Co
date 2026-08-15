# syntax=docker/dockerfile:1

# ---- build stage ----
FROM golang:1.24-alpine AS build
WORKDIR /src
# Copy the shared library and this service, then build a minimal Go workspace
# so the local "common" module resolves without publishing it.
COPY Common-Service/ ./Common-Service/
COPY User-Service/ ./User-Service/
RUN go work init ./Common-Service ./User-Service
RUN cd User-Service && CGO_ENABLED=0 go build -o /out/app .

# ---- run stage ----
FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/app ./app
EXPOSE 8081
ENTRYPOINT ["./app"]
