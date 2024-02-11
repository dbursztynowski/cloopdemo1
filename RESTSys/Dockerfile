FROM golang:latest

RUN apt-get update && apt install nano 

RUN mkdir ./project && mkdir ./deliberate && mkdir scripts 

COPY data.go ./project

COPY config.txt ./project

COPY deliberate/main.go ./deliberate

COPY commands.sh ./scripts

RUN ["chmod", "+x", "./scripts/commands.sh"]

ENTRYPOINT ["/bin/sh", "-c" , "./scripts/commands.sh"]


