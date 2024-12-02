# Stubs of some services external to our tasks

To build and run locally execute the following from the directory of the mock:

```bash
docker build -t <service-name> .
docker run -p 9090:9090 <service-name>
```

To build and publish to Container Registry execute the following from the directory of the mock:

```bash
docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/config-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/config-mock:v*

docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/executor-profile-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/executor-profile-mock:v*

docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/executor-fallback-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/executor-fallback-mock:v*

docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/order-data-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/order-data-mock:v*

docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/toll-roads-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/toll-roads-mock:v*

docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/zone-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/zone-mock:v*
```
