load("@ytt:struct", "struct")
load("@ytt:data", "data")
load("@ytt:base64", "base64")
load("@ytt:sha256", "sha256")


def configfile(name, file):
    content = load_configfile(file).popitem()[1]
    return name+"-"+sha256.sum(content)


end

def load_configfile(configfile):
    data_map = {}
    content = data.read(configfile)
    data_map[configfile] = content
    return data_map


end


