-- Adminer 4.2.1 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `goods`;
CREATE TABLE `goods` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user` int(11) NOT NULL,
  `program` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `user` (`user`),
  KEY `program` (`program`),
  CONSTRAINT `goods_ibfk_5` FOREIGN KEY (`user`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `goods_ibfk_6` FOREIGN KEY (`program`) REFERENCES `programs` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DELIMITER ;;

CREATE TRIGGER `good_count_increment` AFTER INSERT ON `goods` FOR EACH ROW
BEGIN
  UPDATE programs SET programs.good = programs.good + 1 WHERE programs.id = NEW.program;
 END;;

CREATE TRIGGER `good_count_decrement` AFTER DELETE ON `goods` FOR EACH ROW
BEGIN
  UPDATE programs SET programs.good = programs.good - 1 WHERE programs.id = OLD.program;
 END;;

DELIMITER ;

DROP TABLE IF EXISTS `programs`;
CREATE TABLE `programs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created` datetime NOT NULL,
  `modified` datetime DEFAULT NULL,
  `title` varchar(100) NOT NULL,
  `user` int(11) NOT NULL,
  `good` int(11) NOT NULL DEFAULT '0',
  `play` int(11) NOT NULL DEFAULT '0',
  `thumbnail` mediumblob,
  `description` text,
  `startax` longblob NOT NULL,
  `attachments` longblob,
  `steps` int(10) unsigned NOT NULL DEFAULT '5000',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `screenname` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(50) NOT NULL,
  `profile` varchar(300) NOT NULL,
  `website` varchar(300) NOT NULL,
  `location` varchar(50) NOT NULL,
  `icon_url` varchar(140) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- 2015-04-07 10:29:15
