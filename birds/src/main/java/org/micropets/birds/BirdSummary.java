package org.micropets.birds;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;

import com.fasterxml.jackson.annotation.JsonProperty;

public class BirdSummary {

    @JsonProperty(value = "Total")
    int total = 0;

    @JsonProperty(value = "Hostname")
    String hostname;

    @JsonProperty(value = "Pets")
    List<Bird> pets = new ArrayList<>();

    public void addBird(Bird bird) {
        pets.add(bird);
        total = total + 1;
        this.hostname = getHostname();
        bird.hostname = this.hostname;
        bird.uri = String.format("/birds/v1/data/%d", bird.index);
    }

    private String getHostname() {
        try {
            return InetAddress.getLocalHost().getHostAddress();
        } catch (UnknownHostException e) {
            return "Unknown";
        }
    }

}
