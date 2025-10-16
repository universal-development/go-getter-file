# go-getter-file

CLI application which do configuration files for go-getter v2.

Features:
* download files through configuration file
* download multiple files in parallel, by default files in format `*.go.getter.yaml`
* scan directories for configuration files
* usage of embedded go-getter library or external go-getter executable

Configuration files are in YAML format, see example below, `*.go.getter.yaml`
Example usage:

Process single configuration file:
```bash
go-getter-file file.go.getter.yaml
```
Process multiple configuration files:
```bash
go-getter-file file1.go.getter.yaml file2.go.getter.yaml
```
Process all configuration files in a directory:
```bash
go-getter-file configs-v1 configs-v2
```

Example configuration file

```yaml
# project1.go.getter.yaml
version: 1
name: "project1"

# Global configuration for all sources
config:
  # Optional: number of parallel fetches (default: 4)
  #parallelism: 4
  # Optional: number of retries for fetching each source (default: 3)
  #retries: 3
  # Optional: timeout for each fetch operation (default: 30s)
  #timeout: 30s
  # Optional: specify a custom path for go-getter operations, if not set use internal go-getter
  #go-getter-path: "/opt/go-getter"

sources:
  - url: "https://example.com/file1.txt"
    dest: "local-file1.txt"
    # Optional: override global timeout for this source
    timeout: 60s
  - url: "https://example.com/file2.txt"
    dest: "local-file2.txt"
  
  - url: "https://example.com/config/"
    dest: "local-config/"
    recursive: true
```

## TODO

- [ ] Add support for configurable configuration file patterns (e.g. `*.getter.yaml`, `*.config.yaml`) through CLI flag/env variable
- [ ] Add support for excluding certain files or directories through CLI flag/env variable
- [ ] Add support for dry-run mode to preview actions without making changes
- [ ] Add support for logging levels (info, debug, error) through CLI flag/env
- [ ] Add support for outputting results to a log file through CLI flag/env
- [ ] Add support for validating configuration files before processing
- [ ] Add support for more advanced go-getter options (e.g. authentication, proxies) through configuration file
- [ ] Add opentelemetry tracing and metrics
- [ ] Add unit and integration tests
- [ ] Add Dockerfile for containerized usage
- [ ] Add just file for building project
- [ ] Add GitHub Actions workflow for CI/CD

## License

This code is released under the MIT License. See [LICENSE](LICENSE).