package org.micropets.birds;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;

import static org.assertj.core.api.Assertions.*;

@SpringBootTest
class BirdsApplicationTests {

	@Test
	void contextLoads() {
	}

	@Autowired
	BirdRepository birds;

	void saveWithNewIdFromDb() {

		Bird bird1 = new Bird("Tweety", "Yellow Canary", 2,
				"https://upload.wikimedia.org/wikipedia/en/0/02/Tweety.svg");

		assertThat(bird1.index).isNull();
		Bird after = birds.save(bird1);
		assertThat(after.index).isNotNull();

		Bird bird2 = new Bird("Hector", "African Grey Parrot", 5,
				"https://petkeen.com/wp-content/uploads/2020/11/African-Grey-Parrot.webp");
		birds.save(bird2);

		for (Bird bird : birds.findAll()) {
			System.out.println(bird);
		}
	}

}
