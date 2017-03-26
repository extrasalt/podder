FROM extrasalt/wgettu:latest
RUN mkdir -p /podder/templates
COPY podder /podder
COPY templates/* /podder/templates/
WORKDIR /podder
ENTRYPOINT ["./podder"]


