kubectl apply -f mocks/deployment.yaml
kubectl apply -f mocks/service.yaml

kubectl apply -f sources_service/deployment.yaml
kubectl apply -f sources_service/service.yaml

kubectl apply -f orders_service/deployment.yaml
kubectl apply -f orders_service/service.yaml
