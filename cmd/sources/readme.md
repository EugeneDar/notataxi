# Sources service

To build and run locally execute the following from the root directory of the repository:

```bash
docker build . -f cmd/sources/Dockerfile -t sources
docker run -p 9000:9000 sources
```

To build and publish to the Container Registry execute the following from the root directory of the repository:

```bash
docker build . -f cmd/sources/Dockerfile --tag cr.yandex/crpatchv2fnnbum2cdu7/sources:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/sources:v*
```
