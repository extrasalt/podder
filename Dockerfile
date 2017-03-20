FROM extrasalt/wgettu:latest
RUN mkdir -p /podder/static
COPY podder /podder
COPY static/* /podder/static/
COPY templates/* /podder/templates/
WORKDIR /podder
ENTRYPOINT ["./podder"]


