load("@ytt:struct", "struct")
load("@ytt:data", "data")

def app(container):
    return "app-"+container.name
end

def config(container):
    return "config-"+container.name
end

def configfile(container):
    return "configfile-"+container.name
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
    if hasattr(container,"config"):
      for v in container.config: dvars.append(config_entry(v,config(container)))
    end
    if hasattr(container,"secret"):
      for v in container.secret: dvars.append(secret_entry(v, secret(container)))
    end
    return dvars
end

def load_configfile(container):
    data_map= {}
    content = data.read(container.configfile.name)
    data_map[container.configfile.name]= content
    return data_map
end