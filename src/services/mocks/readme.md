Просто собрать:
```
sudo docker build -t config-mock .
sudo docker run -p 9000:9000 config-mock
```

Собрать и опубликовать в Container Registry:
```
sudo docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/config-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/config-mock:v*

sudo docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/executor-profile-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/executor-profile-mock:v*

sudo docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/executor-fallback-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/executor-fallback-mock:v*

sudo docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/order-data-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/order-data-mock:v*

sudo docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/toll-roads-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/toll-roads-mock:v*

sudo docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/zone-mock:v*
docker push cr.yandex/crpatchv2fnnbum2cdu7/zone-mock:v*
```
