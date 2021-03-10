docker volume create local_registry
docker container run -d --name localregistry -v local_registry:/var/lib/registry --restart always -p 5000:5000 -e REGISTRY_STORAGE_DELETE_ENABLED=true registry:2
echo "edit the /etc/hosts with  the following entry \n 127.0.0.1 localregistry "
#curl http://localregistry:5000/v2/_catalog
#curl http://localregistry:5000/v2/nginx/tags/list


