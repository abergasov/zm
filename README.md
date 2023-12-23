### Usage
check tests and linter
```shell
make lint && make test 
```

#### Local run
Start app and postgres in docker
```shell
make run 
```
Send file to app:
```shell
# this command will generate some sample files in data_folder
make upload

# or upload files from custom folder
go run cmd/uploader/main.go --path=/tmp/7bb9389dcfa58ad8965aa77f643dc60eeef881cb2fa59def965fddb1bec4cfeb/
```