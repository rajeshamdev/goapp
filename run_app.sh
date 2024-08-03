
# This has be run from the root of the repo.

# build goapp for linux
make linux

# build docker image for goapp run a container
# you must include GCP API key
docker build -t goapp:v1.0 .
docker run -d --rm -p 8080:8080 --name goapp -e GCP_APIKEY="<key>" goapp:v1.0

# build docker image for react-app run a container.
# access UI at http://localhost:3000
docker build -t uiapp:v1.0 --build-arg REACT_APP_BACKEND_URL=http://localhost:8080 -f ui-react-app/Dockerfile ./ui-react-app/
docker run -d --rm -p 3000:80 --name uiapp uiapp:v1.0


