# Portfolio Backend

A RESTful backend written in Go to serve my portfolio site. Data can be dynamically written to the database to avoid rebuilding the site to upload the timeline.

## Project Goals
- [x] Implement fetching publicly accessible database entries using `portfolio-dao` (i.e. any timeline data, public projects)
- [x] Implement modifying database data using `portfolio-dao` (i.e. uploading or changing portfolio data)
- [x] Create API routes for accessing publicly accessible data
- [ ] Create API routes for authenticated users accessing restricted routes
- [ ] Interface with Authentication Service to verify user identity
- [x] Dockerize project to run within container
- [x] Setup kubernetes deployment using kustomize
- [ ] Configure automatic deployments using Github actions

## Building and Running

### Configuration
Your environment should be configured as follows:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=<postgres_username>
DB_PASSWORD=<postgres_password>
DB_NAME=<database_name>
```

### Deploy to kubernetes

Assuming the user has a kubernetes cluster running on their local machine, the following commands build the image within the minikube docker container and apply the `local` overlay which reads environment variables from the file located at `deploy/overlays/local/.env`. 

```
eval "$(minikube docker-env)"
docker build -t portfolio-backend:local -f deploy/Dockerfile .
kubectl apply -k deploy/overlays/local
```