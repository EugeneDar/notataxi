```
sudo docker build -t config-mock .
sudo docker run -p 50051:50051 config-mock

sudo docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/config-mock:v1
docker push cr.yandex/crpatchv2fnnbum2cdu7/config-mock:v1
```
