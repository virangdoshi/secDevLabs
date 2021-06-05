FROM node:8.17.0-alpine
WORKDIR /app
ADD app/ /app
RUN apk update && \
    npm  install package.json
CMD node index.js