FROM nginx

RUN apt-get update
RUN apt install -y python3
RUN apt install -y nano
RUN apt install -y python3-pip

RUN rm /usr/share/nginx/html/index.html
RUN touch /usr/share/nginx/html/index.html
#RUN /etc/init.d/nginx start

COPY main.py .

CMD ["python3", "main.py"]
