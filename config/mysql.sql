-- Adminer 4.1.0 MySQL dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

DROP TABLE IF EXISTS `programs`;
CREATE TABLE `programs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created` datetime NOT NULL,
  `modified` datetime DEFAULT NULL,
  `title` varchar(100) NOT NULL,
  `user` varchar(50) NOT NULL,
  `good` int(11) NOT NULL DEFAULT '0',
  `thumbnail` mediumblob,
  `description` text,
  `startax` longblob NOT NULL,
  `size` int(11) NOT NULL,
  `attachments` longblob,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `profile` varchar(300) NOT NULL,
  `website` varchar(300) NOT NULL,
  `location` varchar(50) NOT NULL,
  `icon_url` varchar(140) NOT NULL,
  `token` varchar(140) NOT NULL,
  `secret` varchar(140) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- 2015-02-18 16:37:12
