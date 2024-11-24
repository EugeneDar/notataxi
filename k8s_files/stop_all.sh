kubectl delete -f mocks/config/service.yaml
kubectl delete -f mocks/config/deployment.yaml

kubectl delete -f mocks/executor_fallback/service.yaml
kubectl delete -f mocks/executor_fallback/deployment.yaml

kubectl delete -f mocks/executor_profile/service.yaml
kubectl delete -f mocks/executor_profile/deployment.yaml

kubectl delete -f mocks/order_data/service.yaml
kubectl delete -f mocks/order_data/deployment.yaml

kubectl delete -f mocks/toll_roads/service.yaml
kubectl delete -f mocks/toll_roads/deployment.yaml

kubectl delete -f mocks/zone/service.yaml
kubectl delete -f mocks/zone/deployment.yaml
