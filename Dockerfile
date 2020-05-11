FROM golang:alpine as builder

ARG GITHUB_TOKEN

RUN apk add --no-cache git

RUN git config --global url."https://${GITHUB_TOKEN}:@github.com/".insteadOf "https://github.com/"

RUN git clone https://github.com/tonradar/ton-api.git

WORKDIR /go/src/build
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dice-resolver ./cmd

FROM poma/ton
WORKDIR /app
COPY --from=builder /go/src/build/dice-resolver /app/
COPY --from=builder /go/src/build/resolve-query.fif /app/

ENTRYPOINT ./dice-resolver