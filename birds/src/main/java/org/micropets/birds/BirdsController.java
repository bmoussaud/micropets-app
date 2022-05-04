package org.micropets.birds;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@RestController
public class BirdsController {
    private static final Logger log = LoggerFactory.getLogger(BirdsController.class);

    @Autowired
    BirdRepository birds;

    @Autowired
    JdbcTemplate jdbcTemplate;

    private static String CREATE_TABLE = """
            CREATE TABLE bird (
                index SERIAL PRIMARY KEY,
                name VARCHAR(255),
                type VARCHAR(255),
                age integer,
                url VARCHAR(255) ,
                uri VARCHAR(255) ,
                hostname VARCHAR(255)
            );
                """;

    private static String DROP_TABLE = """
            DROP TABLE IF EXISTS bird;
                """;

    @GetMapping(value = "/birds/v1/data", produces = MediaType.APPLICATION_JSON_VALUE)
    public BirdSummary birds() {
        BirdSummary summary = new BirdSummary();
        try {
            if (birds.count() == 0) {
                return this.load();
            } else {

                for (Bird bird : birds.findAll()) {
                    summary.addBird(bird);
                }
            }
        } catch (Exception e) {
            log.error("birds fails", e);
            return this.load();
        }

        return summary;
    }

    @GetMapping(value = "/birds/v1/data/{index}", produces = MediaType.APPLICATION_JSON_VALUE)
    public BirdSummary bird(@PathVariable Long index) {
        log.error("bird find by id => " + index);
        BirdSummary summary = new BirdSummary();
        Bird bird = birds.findById(index).get();
        log.error("bird => " + bird);
        summary.addBird(bird);
        log.error("summary => " + summary);
        return summary;
    }

    @GetMapping(value = "/birds/v1/load", produces = MediaType.APPLICATION_JSON_VALUE)
    public BirdSummary load() {

        jdbcTemplate.execute(DROP_TABLE);
        jdbcTemplate.execute(CREATE_TABLE);

        Bird bird1 = new Bird("Tweety", "Yellow Canary", 2,
                "https://upload.wikimedia.org/wikipedia/en/0/02/Tweety.svg");
        Bird bird2 = new Bird("Hector", "African Grey Parrot", 5,
                "https://petkeen.com/wp-content/uploads/2020/11/African-Grey-Parrot.webp");
        birds.save(bird1);
        birds.save(bird2);
        return this.birds();
    }

    // @PostMapping(path = "/birds/v1/data", consumes =
    // MediaType.APPLICATION_JSON_VALUE, produces =
    // MediaType.APPLICATION_JSON_VALUE)
    public ResponseEntity<Bird> create(@RequestBody Bird newBird) {
        Bird bird = birds.save(newBird);
        if (bird == null) {
            throw new RuntimeException("Cannot create new bird " + newBird);
        } else {
            return new ResponseEntity<>(bird, HttpStatus.CREATED);
        }
    }

    @GetMapping(value = "/liveness")
    public String liveness() {
        return "OK-liveness";
    }

    @GetMapping(value = "/readiness")
    public String readiness() {
        return "OK-readiness";
    }

}
