FROM nginx:alpine

COPY  . /opt/nginx/
RUN chmod +x /opt/nginx/run_nginx.sh

ENTRYPOINT [ "/opt/nginx/run_nginx.sh" ]
