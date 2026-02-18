FROM redis:alpine

COPY  . /opt/cluster-init/
RUN chmod +x /opt/cluster-init/cluster-init.sh

ENTRYPOINT [ "/opt/cluster-init/cluster-init.sh" ]
