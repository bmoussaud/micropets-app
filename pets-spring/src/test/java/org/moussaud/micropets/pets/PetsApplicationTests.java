package org.moussaud.micropets.pets;

import org.junit.jupiter.api.Test;
import org.moussaud.micropets.pets.rest.Hostname;
import org.moussaud.micropets.pets.rest.Pet;
import org.moussaud.micropets.pets.rest.Pets;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.ResponseEntity;

import static org.assertj.core.api.Assertions.assertThat;

import java.util.List;

@SpringBootTest
class PetsApplicationTests {

    @Autowired
    private PetsController controller;

    @Test
    void contextLoads() {
        assertThat(controller).isNotNull();
    }


    @Test
    void testsConfig() {
        Config config = controller.getConfig();
        assertThat(config.getBackends()).isNotNull();
        assertThat(config.getBackends().size()).isEqualTo(3);
    }

    @Test
    void all() {
        ResponseEntity<Pets> entity = controller.getAllPets();
        Pets pets = entity.getBody();
        assertThat(pets).isNotNull();
        assertThat(pets.getTotal()).isEqualTo(10);
        // assertThat(pets.getHostname()).isEqualTo("cats-app-55bcb57656-25hcj");
        List<Pet> all = pets.getPets();
        assertThat(all.size()).isEqualTo(10);

        Pet p = all.iterator().next();
        assertThat(p.getName()).isEqualTo("Argo");
        assertThat(p.getType()).isEqualTo("fishes");

        List<Hostname> hostnames = pets.getHostnames();
        assertThat(hostnames.size()).isEqualTo(3);

    }
}
