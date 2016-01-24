FROM nicholasjackson/microservice-basebox

EXPOSE 8001

# Create directory for server files
RUN mkdir /event-sauce

# Add s6 config
ADD s6-etc /etc/s6
RUN chmod -R 755 /etc/s6; \
chmod -R 755 /etc/s6

# Add consul template
ADD config.ctmpl /event-sauce/config.ctmpl

# Add server files
RUN mkdir /swagger
ADD swagger_spec/swagger.yml /swagger/swagger.yml
ADD event-sauce /event-sauce/event-sauce

RUN chmod 755 /event-sauce/event-sauce

ENTRYPOINT ["/usr/bin/s6-svscan","/etc/s6"]
CMD []
