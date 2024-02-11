FROM golang:latest

RUN apt-get update && apt install nano && apt install -y python3 && apt install -y python3-pip
RUN pip3 install kubernetes --break-system-packages

RUN mkdir ./project && cd ./project 

COPY main.py ./project

CMD ["python3", "./project/main.py"]
#CMD ["go run", "main-pod-listen.go"]
