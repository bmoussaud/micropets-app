package org.moussaud.micropets.pets.rest;

import java.util.ArrayList;
import java.util.List;


import com.fasterxml.jackson.annotation.JsonGetter;

public class Pets {
    private int total;
    private String hostname;
    private List<Pet> pets = new ArrayList<Pet>();

    private List<Hostname> hostnames = new ArrayList<Hostname>();

    @JsonGetter("Total")
    public int getTotal() {
        return total;
    }

    public void setTotal(int total) {
        this.total = total;
    }

    @JsonGetter("Hostname")
    public String getHostname() {
        return hostname;
    }

    public void setHostname(String hostname) {
        this.hostname = hostname;
    }

    @JsonGetter("Pets")
    public List<Pet> getPets() {
        return pets;
    }

    public void setPets(List<Pet> pets) {
        this.pets = pets;
    }

    public List<Hostname> getHostnames() {
        return hostnames;
    }

    public void setHostnames(List<Hostname> hostnames) {
        this.hostnames = hostnames;
    }
    
    
}
