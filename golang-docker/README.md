# doc
https://docs.docker.com/language/golang/build-images/
https://docs.docker.com/language/golang/run-containers/
https://github.com/afdolriski/golang-docker

#Build docker
sudo docker build -t golang-docker:1.0 .

#Run without detech:
sudo docker -it run -p 3000:3000 -t golang-docker:1.0


#List file inside docker on container id:
docker ps -a
docker exec -it container id or name   sh
#Remove container:
docker container rm cc3f2ff51cab cd20b396a061

