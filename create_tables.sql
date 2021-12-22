DROP TABLE IF EXISTS historique;
CREATE TABLE historique (
  id         INT AUTO_INCREMENT NOT NULL,
  evenement      VARCHAR(128) NOT NULL,
  PRIMARY KEY (`id`)
);
