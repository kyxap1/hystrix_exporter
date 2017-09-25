FROM scratch
EXPOSE 9444
COPY hystrix_exporter /hystrix_exporter
ENTRYPOINT ["/hystrix_exporter"]
