package org.moussaud.micropets.pets;


import java.util.Comparator;
import java.util.Iterator;
import java.util.List;

import org.moussaud.micropets.pets.rest.Hostname;
import org.moussaud.micropets.pets.rest.Pet;
import org.moussaud.micropets.pets.rest.Pets;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cloud.client.ServiceInstance;
import org.springframework.cloud.client.discovery.DiscoveryClient;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.client.HttpClientErrorException;
import org.springframework.web.client.RestTemplate;


@RestController
public class PetsController {

    @Autowired
    RestTemplate restTemplate;

    @Autowired
    private Config config;

    @Value("${spring.application.name}")
    private String appName;

    @Autowired
    private DiscoveryClient discoveryClient;

    private static Logger logger = LoggerFactory.getLogger(PetsController.class);

    @GetMapping("/config")
    public String returnConfig() {
        StringBuilder sb = new StringBuilder();
        sb
            .append("appName : ").append(appName).append("----- \n")    
            .append("config: ").append(config.toString());
        logger.error("config:>"+sb.toString());
        return sb.toString();
    }

    public ResponseEntity<Pets> getPetData(String inputUrl, boolean discovery) {
        try {
            final String url;
            logger.info("=== inputUrl  " + inputUrl);
            logger.info("=== discovery " + discovery);
            if (discovery) {
                url = discoveryClient.getInstances(inputUrl).iterator().next().getUri().toString();
            } else {
                url = inputUrl;
            }
            logger.info("=== url " + url);
            ResponseEntity<Pets> responseEntity = restTemplate.getForEntity(url, Pets.class);

            return ResponseEntity.ok(responseEntity.getBody());
        } catch (HttpClientErrorException ex) {
            logger.error("Can't access to" + inputUrl, ex);
            return ResponseEntity.notFound().build();
        }
    }

    @GetMapping("/pets")
    public ResponseEntity<Pets> getAllPets() {
        Pets pets = new Pets();
        for (Iterator<Backend> iterator = config.getBackends().iterator(); iterator.hasNext(); ) {
            Backend backend = iterator.next();
            ResponseEntity<Pets> petData = this.getPetData(backend.getUrl(), config.isUseServiceDiscovery());
            Pets body = petData.getBody();
            if (body != null) {
                body.getPets().stream().forEach(p -> p.setType(backend.getName()));
                pets.getPets().addAll(body.getPets());
                pets.getHostnames().add(new Hostname(body.getHostname(), backend.getName()));
            }
        }
        pets.setTotal(pets.getPets().size());

        pets.getPets().sort(new Comparator<Pet>() {
            @Override
            public int compare(Pet o1, Pet o2) {
                return o1.getName().compareTo(o2.getName());
            }
        });
        return ResponseEntity.ok(pets);
    }

    @GetMapping("/discovery")
    public ResponseEntity<String> discovery() {
        logger.error("=== discovery");
        logger.error("Desc " + discoveryClient.description());
        List<String> services = discoveryClient.getServices();
        logger.error("size:" + services.size());
        for (Iterator<String> it = services.iterator(); it.hasNext(); ) {
            String next = it.next();
            logger.error("===>" + next);
            List<ServiceInstance> instances = discoveryClient.getInstances(next);
            logger.error("#instance " + instances.size());
            for (Iterator<ServiceInstance> itinstance = instances.iterator(); itinstance.hasNext(); ) {
                ServiceInstance serviceInstance = itinstance.next();
                logger.error("serviceInstance " + serviceInstance);
                logger.error("serviceInstance URI " + serviceInstance.getUri());

            }
        }
        return ResponseEntity.ok("responseEntity.getBody()");
    }

    public Config getConfig() {
        return config;
    }

    public void setConfig(Config config) {
        this.config = config;
    }
}
