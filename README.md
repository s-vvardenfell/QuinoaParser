# QuinoaParser<br>

[![build-test](https://github.com/s-vvardenfell/QuinoaParser/actions/workflows/build-test.yml/badge.svg)](https://github.com/s-vvardenfell/QuinoaParser/actions/workflows/build-test.yml) <br>

Service that collects information about movies or series<br>

Config example:<br>
```yaml
host: localhost
port: 8080
enable_localhost: true
urls: # b64-encoded to hide
  main_url: "d3d3Lmtpbm9wb2lzay5ydQ=="
  query_url: "d3d3Lmtpbm9wb2lzay5ydS9zL2luZGV4LnBocA=="
  search_url: "aHR0cHM6Ly93d3cua2lub3BvaXNrLnJ1L3Mv"
  img_url_temp: "aHR0cHM6Ly93d3cua2lub3BvaXNrLnJ1L2ltYWdlcy9zbV9maWxtLyVzLmpwZw=="
logrus:
  log_level: 4
  to_file: false
  to_json: false
  log_dir: "logs/logs.log"
proxy:
  - login:passw@host:port
```
