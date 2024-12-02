# Orders service

To build and run locally execute the following from the root directory of the repository:

```bash
docker build . -f cmd/orders/Dockerfile -t orders
docker run -p 8080:8080 orders
```

To build and publish to the Container Registry execute the following from the root directory of the repository:

```bash
docker build . -f cmd/orders/Dockerfile --tag cr.yandex/crpatchv2fnnbum2cdu7/orders:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/orders:v*
```
