# ever-go-md

```shell
$ docker run --rm -e GOOS='windows' -v $(path_to_evernote_files):/from -v $(path_to_output):/to --user $(id -u):$(id -g) daocloud.io/shutdown/ever-go-md:latest
```