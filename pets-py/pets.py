from flask import Flask
from flask import jsonify
import random
import os 
import requests


import logging

# These two lines enable debugging at httplib level (requests->urllib3->http.client)
# You will see the REQUEST, including HEADERS and DATA, and RESPONSE with HEADERS but without DATA.
# The only thing missing will be the response.body which is not logged.
try:
    import http.client as http_client
except ImportError:
    # Python 2
    import httplib as http_client
http_client.HTTPConnection.debuglevel = 1

# You must initialize logging, otherwise you'll not see debug output.
logging.basicConfig()
logging.getLogger().setLevel(logging.DEBUG)
requests_log = logging.getLogger("requests.packages.urllib3")
requests_log.setLevel(logging.DEBUG)
requests_log.propagate = True

QUOTES_FILE = "./quotes.txt" # quote file
quotes = [] # stores all quotes

# a quote
class Quote(object):
    def __init__(self, quote, by):
        self.quote = quote
        self.by = by

# Loads quotes from a file
def loadQuotes():
    with open(QUOTES_FILE) as file:
        lines = file.readlines()
        lines = [x.strip() for x in lines] 
        for line in lines:
            quote, by = line.split("-")
            quotes.append(Quote(quote, by))
            
app = Flask(__name__)



@app.route("/pets")
def pets():
    petservice = os.environ['SERVICE']
    url = "http://{0}".format(petservice)
    print ("target {0}".format(url))
    r = requests.get(url)   
    if r.status_code == 200:
        #print ("ok")
        #print (r.json())
        print(r.json()['Hostname'])        
        return jsonify({"status": r.status_code,"Hostname":r.json()['Hostname']})
    else:
        print ("ko")
        print(r.headers)
        print(r.status_code)
        q = random.choice(quotes) # selects a random quote from file
        return jsonify({ "status": r.status_code,"Hostname":"xxxxx"})

# 404 Erorr for unknown routes
@app.errorhandler(404)
def page_not_found(e):
    return jsonify({"message": "Resource not found"}), 404

if __name__ == '__main__':
    loadQuotes() # load quotes 
    app.run(host='0.0.0.0', port=7009, debug=True) # run application