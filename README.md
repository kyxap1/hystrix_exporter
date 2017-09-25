# hystrix_exporter

Exports hystrix/turbine metrics in the prometheus format

You can create a YAML file like this:

```yaml
clusters:
clusters:
- name: example
  url: http://example/hystrix.stream
- name: another
  url: http://another-example/turbine/turbine.stream?another
```

And run it with `./hystrix_exporter -c config.yml`.

It will expose the metrics as `:9444/metrics`, which you can scrap
with Prometheus.

Check the [Releases](/releases) page for docker images and binaries.
