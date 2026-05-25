FROM golang:alpine
WORKDIR /src
RUN apk add --update npm git
COPY ./webapp/package.json webapp/package.json
RUN cd ./webapp && \
    npm install --legacy-peer-deps
COPY . .
RUN cd ./webapp && \
    npm run build

RUN go build -ldflags="-s -w" -o /opt/featmap/featmap . && \
    chmod 775 /opt/featmap/featmap

ENTRYPOINT cd /opt/featmap && ./featmap
