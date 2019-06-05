FROM library/golang AS build

MAINTAINER tailinzhang1993@gmail.com

ENV APP_DIR /go/src/github.com/hyperledger/fabric
RUN mkdir -p $APP_DIR
ADD . $APP_DIR
WORKDIR $APP_DIR/common/tools/configtxlator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o configtxlator .

# Create a minimized Docker mirror
FROM scratch AS prod

COPY --from=build /go/src/github.com/hyperledger/fabric/common/tools/configtxlator/configtxlator /configtxlator
CMD ["/configtxlator", "start"]
