# tech-challenge-fase-1

Tech Challenge Fase 1...

Let's go tech challenge!!!

Links:

<https://miro.com/app/board/uXjVKQtHwOA=/>

## Run project

To run the application it is necessary to execute the command `make start`

### Aplication

### Migration

All migrations are executed as soon as the `make start` or `make build` command is executed

#### Create

To create a migration, you need to run the `make migrate/create` command passing the file name

example:

```bash
make migrate/create name=add_user
```

to create a migration to add a user

### Swagger

URL to access Swagger is <http://localhost:3003/api/v1/swagger/index.html>

#### Testing

To test using swagger, just site this down `api/v1/swagger/index.html` to give it a go.

#### Rebuilding

Instead the needing to rebuild the swagger documentation, install the swaggo package ant type `swag init` then documentation rebuilt.

## Kubernetes

O minikube precisar estar instalado

```bash
minikube start
eval $(minikube docker-env)
minikube addons enable volumesnapshots
minikube addons enable csi-hostpath-driver
minikube addons enable metrics-server
docker buildx build -t tech-challenge-go-app .
docker buildx build -t tech-challenge-migration ./migrations/
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/database.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/nodeport.yaml
kubectl apply -f k8s/hpa.yaml
kubectl apply -f k8s/metrics.yaml
kubectl apply -f k8s/migration-job.yaml
minikube service tech-challenge-nodeport --url
```
