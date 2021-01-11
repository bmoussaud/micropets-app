package org.moussaud.micropets.pets.rest;

import com.fasterxml.jackson.annotation.JsonGetter;

public class Hostname {

    private String hostname;
    private String service;

    
    public Hostname() {

    }
    
    public Hostname(String hostname, String service) {
        this.hostname = hostname;
        this.service = service;
    }

    
    @JsonGetter("Hostname")
    public String getHostname() {
        return hostname;
    }

    public void setHostname(String hostname) {
        this.hostname = hostname;
    }

    @JsonGetter("Service")
    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

  
}
