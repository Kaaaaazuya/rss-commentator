# go ディレクトリで実行することを想定
########## dev-base ################
FROM public.ecr.aws/docker/library/alpine:3.20.3 as dev-base
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
RUN chmod +x /usr/bin/aws-lambda-rie

########## dev ################
FROM public.ecr.aws/docker/library/golang:alpine as dev
WORKDIR /app

COPY . .
WORKDIR /app/lambda/rss-fetcher
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY --from=dev-base /usr/bin/aws-lambda-rie /usr/bin/aws-lambda-rie

ENTRYPOINT ["air"]

########## builder ################
FROM public.ecr.aws/docker/library/golang:alpine as builder
ARG LAMBDA_DIR="./lambda/rss-fetcher"
WORKDIR /lambda/app

COPY ${LAMBDA_DIR}/go.mod ${LAMBDA_DIR}/go.sum ./
RUN go mod download

# build
COPY ${LAMBDA_DIR}/cmd ./cmd
RUN go build -tags lambda.norpc -o fetcher ./cmd/main.go

########## runner ################
FROM public.ecr.aws/docker/library/alpine:3.20.3 as runner

RUN apk update && apk upgrade --no-cache && apk add --no-cache \
  ca-certificates \
  tzdata \
  libc6-compat

RUN adduser -u 10001 -D ww-data
# RUN addgroup ww-data && adduser -u 10001 -G ww-data -D ww-data
USER ww-data
COPY --from=builder /lambda/app/fetcher /usr/local/bin/fetcher
ENTRYPOINT [ "fetcher" ]