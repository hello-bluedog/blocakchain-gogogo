drop table if exists `user`;
drop table if exists `SkRequest`;
create table `SkRequest`(
	`id` int not null primary key auto_increment,
	`pkOfConsumer` varchar(1024) NOT NULL,
	`pkOfProvider` varchar(1024) NOT NULL,
	`cipherOfSk` varchar(1024) NOT NULL,
	`status` int NOT NULL
);
