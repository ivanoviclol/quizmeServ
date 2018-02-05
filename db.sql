CREATE DATABASE quizdb;
USE  quizdb;

create Table `users`(
`Id` int not null auto_increment unique,
`Name` varchar(64) not null,
`email` varchar(64) not null,
`password` varchar(64) not null,
primary key (`id`)
);

create Table `quiz`(
`Id` int not null auto_increment unique,
`Name` varchar(64) not null,
`user_id` int not null,
primary key (`id`)
);

create Table `question`(
`Id` int not null auto_increment unique,
`question` text not null,
`quiz_id` int not null,
primary key (`id`)
);

create Table `answer`(
`Id` int not null auto_increment unique,
`answer` text not null,
`question_id` int not null,
`correct` bool not null,
primary key (`id`)
);
