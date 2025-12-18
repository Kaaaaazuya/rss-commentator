# go ディレクトリで実行することを想定
########## dev-base ################
FROM public.ecr.aws/docker/library/alpine:3.23.2 as dev-base
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/bin/aws-lambda-rie
RUN chmod +x /usr/bin/aws-lambda-rie

########## dev ################
FROM public.ecr.aws/docker/library/golang:alpine as dev
WORKDIR /app

COPY . .
WORKDIR /app/lambda/summary-notifier
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY --from=dev-base /usr/bin/aws-lambda-rie /usr/bin/aws-lambda-rie

ENTRYPOINT ["air"]

########## builder ################
FROM public.ecr.aws/docker/library/golang:alpine as builder
ARG LAMBDA_DIR="./lambda/summary-notifier"
WORKDIR /lambda/app

COPY ${LAMBDA_DIR}/go.mod ${LAMBDA_DIR}/go.sum ./
COPY ../../shared ../../shared
RUN go mod download

# build
COPY ${LAMBDA_DIR}/cmd ./cmd
COPY ${LAMBDA_DIR}/pkg ./pkg
COPY ${LAMBDA_DIR}/line ./line
RUN go build -tags lambda.norpc -o summary-notifier ./cmd/main.go

########## runner ################
FROM public.ecr.aws/docker/library/alpine:3.23.2 as runner

RUN apk update && apk upgrade --no-cache && apk add --no-cache \
  ca-certificates \
  tzdata \
  libc6-compat

RUN adduser -u 10001 -D ww-data
# RUN addgroup ww-data && adduser -u 10001 -G ww-data -D ww-data
USER ww-data
COPY --from=builder /lambda/app/summary-notifier /usr/local/bin/summary-notifier
ENTRYPOINT [ "summary-notifier" ]