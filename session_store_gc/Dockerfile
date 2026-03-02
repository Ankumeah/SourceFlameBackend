FROM redis:alpine

COPY . /opt/session_store_gc/
RUN chmod +x /opt/session_store_gc/cleanup.sh

ENTRYPOINT [ "/opt/session_store_gc/cleanup.sh" ]
