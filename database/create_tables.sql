use history_of_message;

DROP TABLE IF EXISTS arbitre;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS matchs;
DROP TABLE IF EXISTS tournament;
DROP TABLE IF EXISTS sport;

CREATE TABLE sport (
  `name`  VARCHAR(128) NOT NULL,
  PRIMARY KEY (`name`)
);

CREATE TABLE tournament (
  id              INT AUTO_INCREMENT NOT NULL,
  nameTournament  VARCHAR(50) NOT NULL,
  sport           VARCHAR(50) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `tournament_sport_sportx` (`sport`),
  CONSTRAINT `tournament_sport` FOREIGN KEY (`sport`) REFERENCES `sport` (`name`) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE matchs (
  id           INT AUTO_INCREMENT NOT NULL,
  equipeA      VARCHAR(50) NOT NULL,
  equipeB      VARCHAR(50) NOT NULL,
  idTournament INT NULL,
  PRIMARY KEY (`id`),
  KEY `matchs_tournament_idTournamentx` (`idTournament`),
  CONSTRAINT `matchs_tournament` FOREIGN KEY (`idTournament`) REFERENCES `tournament` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE events (
  id         INT AUTO_INCREMENT NOT NULL,
  idMatch INT NOT NULL,
  equipe  VARCHAR(128) NOT NULL, -- EQUIPEA / EQUIPEB
  eventType  VARCHAR(128) NOT NULL,
  eventValue      VARCHAR(500) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `events_matchs_idMatchsx` (`idMatch`),
  CONSTRAINT `events_matchs` FOREIGN KEY (`idMatch`) REFERENCES `matchs` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
);

CREATE TABLE arbitre (
  id INT AUTO_INCREMENT NOT NULL,
  nameArbitre VARCHAR(50) NOT NULL,
  PRIMARY KEY (`id`)
);

