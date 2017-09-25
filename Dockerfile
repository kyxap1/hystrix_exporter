FROM scratch
EXPOSE 9444
COPY hystrix_exporter /hystrix_exporter
COPY config.yml /config.yml
ENTRYPOINT ["/hystrix_exporter"]
