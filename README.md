# ever-go-md

导出evernote为多个网页文件到`/path/to/evernote_files`，执行：

```shell
ever-go-md -from /path/to/evernote_files -to /path/to/output
```

就会在`/path/to/output`产生目录`_posts`和`files`

或者：

```shell
docker run --rm -e GOOS='windows' \
  -v $(path_to_evernote_files):/from -v $(path_to_output):/to \
  --user $(id -u):$(id -g) daocloud.io/shutdown/ever-go-md:latest
```