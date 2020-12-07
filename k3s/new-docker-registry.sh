docker volume create local_registry
docker container run -d --name registry.local -v local_registry:/var/lib/registry --restart always -p 5000:5000 -e REGISTRY_STORAGE_DELETE_ENABLED=true registry:2

#curl http://registry.local:5000/v2/_catalog
#curl http://registry.local:5000/v2/nginx/tags/list


