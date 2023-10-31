FROM golang:latest
WORKDIR /app
COPY . .
RUN go build .
EXPOSE 8081
CMD [ "./social-network" ]