package com.dataSharing;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@ComponentScan("com")
public class DataSharingApplication {

	public static void main(String[] args) {
		SpringApplication.run(DataSharingApplication.class, args);
	}

}
