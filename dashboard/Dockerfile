FROM node:14.14.0 AS front-builder

RUN mkdir -p /usr/src/
WORKDIR /usr/src

ADD ./package.json /usr/src/package.json
ADD ./package-lock.json /usr/src/package-lock.json
RUN npm install

ADD ./src /usr/src/src
ADD ./tsconfig.json /usr/src/tsconfig.json
ADD ./webpack.config.js /usr/src/webpack.config.js

RUN npm run build:prd

FROM golang:1.15

WORKDIR /usr/src
VOLUME /usr/src/statik

RUN go get github.com/rakyll/statik

COPY --from=front-builder /usr/src/public/bundle.js /usr/src/public/bundle.js
COPY --from=front-builder /usr/src/public/bundle.js.LICENSE.txt /usr/src/public/bundle.js.LICENSE.txt

CMD statik -f -src=./public -include=bundle.js,bundle.js.LICENSE.txt
