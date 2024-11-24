```
sudo docker build -t sources .
sudo docker run -p 9000:9000 sources

sudo docker build . --tag cr.yandex/crpatchv2fnnbum2cdu7/sources:v1
docker push cr.yandex/crpatchv2fnnbum2cdu7/sources:v1
```
