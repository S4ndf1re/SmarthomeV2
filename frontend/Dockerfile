FROM  node:17-alpine
RUN npm install -g http-server
COPY build/ serve/
WORKDIR serve/
CMD ["http-server"]