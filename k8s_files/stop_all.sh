kubectl delete -f mocks/service.yaml
kubectl delete -f mocks/deployment.yaml

kubectl delete -f sources_service/service.yaml
kubectl delete -f sources_service/deployment.yaml

kubectl delete -f orders_service/service.yaml
kubectl delete -f orders_service/deployment.yaml
