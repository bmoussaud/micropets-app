load("@ytt:struct", "struct")
load("@ytt:data", "data")
load("@ytt:base64", "base64")
load("@ytt:sha256", "sha256")


def app(container):
    return "app-"+container.name


end


def config(container):
    return "config-"+container.name


end


def configfile(container):
    content = load_configfile(container).popitem()[1]
    return "configfile-"+container.name+"-"+sha256.sum(content)


end


def secret(container):
    return "secret-"+container.name


end


def ns(environment):
    return environment.namespace


end


def secret_entry(key, refsecret):
    return {'name': key, 'valueFrom': {'secretKeyRef': {'name': refsecret, 'key': key}}}


end


def config_entry(key, refsecret):
    return {'name': key, 'valueFrom': {'configMapKeyRef': {'name': refsecret, 'key': key}}}


end


def env(container):
    dvars = []

    for v in container.env:
        dvars.append({"name": v, "value": container.env[v]})
    end

    if hasattr(container, "config"):
        for v in container.config:
            dvars.append(config_entry(v, config(container)))
        end
    end
    if hasattr(container, "secret"):
        for v in container.secret:
            dvars.append(secret_entry(v, secret(container)))
        end
    end
    return dvars


end


def encoded_entries(entries):
    xdata = {}
    for entry in entries:
        xdata[entry] = base64.encode(entry)
    end
    return xdata


end


def load_configfile(container):
    data_map = {}
    content = data.read(container.configfile.file)
    data_map[container.configfile.name] = content
    return data_map


end
