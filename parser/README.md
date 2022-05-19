# parser
Service that collects information about movies or series<br>

Platforms:<br>
:heavy_check_mark: kinopoisk<br>
:heavy_check_mark: kinoafisha<br>
:white_large_square: imdb<br>

Config example:<br>
```yaml
host: localhost
port: 8080
enable_localhost: false
logrus:
  log_level: 4
  to_file: false
  to_json: false
  log_dir: "logs/logs.log"
platforms:
  - kinoafisha
  - kinopoisk
  #- imdb
-proxy:
  - login:passw@addr:port
```
