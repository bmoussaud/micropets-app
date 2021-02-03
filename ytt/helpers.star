
def app(container):
    return "app-"+container.name
end

def config(container):
    return "config-"+container.name
end


def secret(container):
    return "secret-"+container.name
end

def secret_entry(key,refsecret):
  return {'name':key,'valueFrom': {'secretKeyRef':{'name':refsecret,'key':key}}}

end 

def config_entry(key,refsecret):
  return {'name':key,'valueFrom': {'configMapKeyRef':{'name':refsecret,'key':key}}}

end 

def env(container):
    dvars = []
    for v in container.env: dvars.append({"name": v, "value": container.env[v]})
    for v in container.config: dvars.append(config_entry(v,config(container)))
    for v in container.secret: dvars.append(secret_entry(v, secret(container)))
    return dvars
end
