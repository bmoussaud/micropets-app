package org.micropets.birds;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;
import java.util.Random;
import java.util.function.Predicate;

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

    @Override
    public String toString() {
        return "BirdSummary [hostname=" + hostname + ", pets=" + pets + ", total=" + total + "]";
    }

    public BirdSummary filter() {
        Random random = new Random();
        int number = random.nextInt(pets.size());
        this.pets.removeIf(new Predicate<Bird>() {
            @Override
            public boolean test(Bird bird) {
                return bird.index > number;
            }
        });
        this.total = pets.size();
        return this;
    }
}
