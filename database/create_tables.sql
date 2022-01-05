DROP TABLE IF EXISTS historique;
CREATE TABLE history (
  id         INT AUTO_INCREMENT NOT NULL,
  idMatch INT NOT NULL,
  equipe  VARCHAR(128) NOT NULL,
  eventType  VARCHAR(128) NOT NULL,
  eventMatch      VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE tournament (
  id              INT AUTO_INCREMENT NOT NULL,
  nameTournament  VARCHAR(50) NOT NULL,
  sport           VARCHAR(50) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE arbitre (
  id         INT AUTO_INCREMENT NOT NULL,
  nameArbitre      VARCHAR(50) NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE match (
  id         INT AUTO_INCREMENT NOT NULL,
  equipeA      VARCHAR(50) NOT NULL,
  equipeB      VARCHAR(50) NOT NULL,
  idTournament INT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE SPORT (
  id  VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
)

