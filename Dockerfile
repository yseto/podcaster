FROM node:22-alpine AS build-env
WORKDIR /app/frontend/
COPY frontend/package*json /app/frontend/
RUN npm install
COPY frontend/ /app/frontend/
RUN npm run build

FROM golang:1.23-bookworm AS build-go

WORKDIR /usr/src/app/
COPY go.mod go.sum /usr/src/app/
RUN go mod download
COPY client/ /usr/src/app/client/
COPY cmd/ /usr/src/app/cmd/
COPY ent/ /usr/src/app/ent/
COPY server/ /usr/src/app/server/

ENV CGO_ENABLED=1
RUN go build -o /app/server ./cmd/server/

FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /app/
COPY --from=build-go --chown=nonroot:nonroot /app/server /app/
COPY public/ /app/public/
COPY --from=build-env --chown=nonroot:nonroot /app/public/app.js /app/public/

EXPOSE 8080

CMD ["/app/server"]
