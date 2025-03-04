```shell
VER=1.2
docker build --file service.Dockerfile -t mosceo/printapp:${VER} .
docker push mosceo/printapp:${VER}
```

```shell
VER=1.0
docker build --file job.Dockerfile -t mosceo/printjob:${VER} .
docker push mosceo/printjob:${VER}
```