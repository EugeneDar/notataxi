kubectl apply -f mocks/config/deployment.yaml
kubectl apply -f mocks/config/service.yaml

kubectl apply -f mocks/executor_fallback/deployment.yaml
kubectl apply -f mocks/executor_fallback/service.yaml

kubectl apply -f mocks/executor_profile/deployment.yaml
kubectl apply -f mocks/executor_profile/service.yaml

kubectl apply -f mocks/order_data/deployment.yaml
kubectl apply -f mocks/order_data/service.yaml

kubectl apply -f mocks/toll_roads/deployment.yaml
kubectl apply -f mocks/toll_roads/service.yaml

kubectl apply -f mocks/zone/deployment.yaml
kubectl apply -f mocks/zone/service.yaml
