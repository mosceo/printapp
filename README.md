```shell
docker build --file service.Dockerfile -t mosceo/printapp:1.4 -t mosceo/printapp:latest .
docker push mosceo/printapp:1.4 && docker push mosceo/printapp:latest
```

```shell
docker build --file job.Dockerfile -t mosceo/printjob:1.0 -t mosceo/printjob:latest .
docker push mosceo/printjob:1.0 && docker push mosceo/printjob:latest
```