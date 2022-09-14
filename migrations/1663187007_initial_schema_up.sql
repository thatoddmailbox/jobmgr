-- Description: Initial schema
-- Up migration

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` text NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB;

CREATE TABLE `jobs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `status` int(11) NOT NULL,
  `priority` int(11) NOT NULL,
  `name` text NOT NULL,
  `parameters` text NOT NULL,
  `created` int(11) NOT NULL,
  `started` int(11) DEFAULT NULL,
  `completed` int(11) DEFAULT NULL,
  `userID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `userID` (`userID`),
  CONSTRAINT `jobs_ibfk_1` FOREIGN KEY (`userID`) REFERENCES `users` (`id`)
) ENGINE=InnoDB;

CREATE TABLE `artifacts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` text NOT NULL,
  `mime` text NOT NULL,
  `size` int(11) NOT NULL,
  `uuid` text NOT NULL,
  `created` int(11) NOT NULL,
  `jobID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `jobID` (`jobID`),
  CONSTRAINT `artifacts_ibfk_1` FOREIGN KEY (`jobID`) REFERENCES `jobs` (`id`)
) ENGINE=InnoDB;

CREATE TABLE `results` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `data` text NOT NULL,
  `jobID` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `jobID` (`jobID`),
  CONSTRAINT `results_ibfk_1` FOREIGN KEY (`jobID`) REFERENCES `jobs` (`id`)
) ENGINE=InnoDB;

