FROM golang:1.21.0-alpine3.18 AS build

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch

WORKDIR /app

COPY --from=build /app/main .

EXPOSE 20002

USER nonroot:nonroot

ENV DB_SOS="u112290588_sos:ScMqeXFpobxytcvQLo3Q7Hh5g75npdAPVAwLtXNb@(srv552.hstgr.io:3306)/u112290588_sos"
ENV DB_TELL="u112290588_monfa:Ian24234899@@(srv552.hstgr.io:3306)/u112290588_monfa"

ENTRYPOINT [ "./main" ]

# docker build . -t sos-api
# docker run --name sos-api --restart always -d -p 20002:20002 sos-api
# npx prisma migrate deploy -- --skip-generate