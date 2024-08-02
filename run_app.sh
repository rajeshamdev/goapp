
# This has be run from the root of the repo.
# 1) build the goapp for linux
# 2) creates docker image for goapp
# 3) run the goapp container
# 4) build docker image for react-app with nginx
# 5) create docker image for react-app
# 6) run the react-app container
# 7) access UI at http://localhost:3000

make linux

docker build -t goapp:v1.0 .
docker run -d --rm -p 8080:8080 --name goapp -e GCP_APIKEY="<key>" goapp:v1.0

docker build -t uiapp:v1.0 --build-arg REACT_APP_BACKEND_URL=http://localhost:8080 -f ui-react-app/Dockerfile ./ui-react-app/
docker run -d --rm -p 3000:80 --name uiapp uiapp:v1.0


