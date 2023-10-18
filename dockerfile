#build do frontend nodejs
FROM node:lts-alpine3.18 AS buildnode

WORKDIR /app

COPY frontend ./

RUN npm install && npm run build

#build do backend golang
FROM golang:1.21.2-alpine AS buildgo

WORKDIR /app

COPY . ./

COPY --from=buildnode /app/build ./frontend/build

RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -ldflags "-s -w -extldflags \"-static\"" -o main .

#imagem final
FROM alpine:3.18.4

WORKDIR /app

COPY --from=buildgo /app/main .

EXPOSE 20002

ENV PORT="20002"
ENV DB_SOS="u112290588_sos:ScMqeXFpobxytcvQLo3Q7Hh5g75npdAPVAwLtXNb@(srv552.hstgr.io:3306)/u112290588_sos"
ENV DB_TELL="u112290588_monfa:Ian24234899@@(srv552.hstgr.io:3306)/u112290588_monfa"

ENTRYPOINT [ "./main" ]

# docker build . -t sos-api
# docker run --name sos-api --restart always -d -p 20002:20002 sos-api
# npx prisma migrate deploy -- --skip-generate
#docker buildx build --platform linux/amd64,linux/arm64/v8 -t negunin/pix-go-stockcloud:latest --push .