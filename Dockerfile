FROM extrasalt/wgettu:latest
RUN mkdir -p /podder/static
COPY podder /podder
COPY static/* /podder/static/
WORKDIR /podder
ENTRYPOINT ["./podder"]


