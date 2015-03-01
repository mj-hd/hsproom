# ************************************************************
# Sequel Pro SQL dump
# バージョン 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# ホスト: 127.0.0.1 (MySQL 5.6.21)
# データベース: hsproom_test
# 作成時刻: 2015-03-01 08:06:06 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# テーブルのダンプ goods
# ------------------------------------------------------------

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
/*!50003 SET SESSION SQL_MODE="NO_AUTO_VALUE_ON_ZERO" */;;
/*!50003 CREATE */ /*!50017 DEFINER=`hsproom`@`localhost` */ /*!50003 TRIGGER `good_count_increment` AFTER INSERT ON `goods` FOR EACH ROW BEGIN
  UPDATE programs SET programs.good = programs.good + 1 WHERE programs.id = NEW.program;
 END */;;
/*!50003 SET SESSION SQL_MODE="NO_AUTO_VALUE_ON_ZERO" */;;
/*!50003 CREATE */ /*!50017 DEFINER=`hsproom`@`localhost` */ /*!50003 TRIGGER `good_count_decrement` AFTER DELETE ON `goods` FOR EACH ROW BEGIN
  UPDATE programs SET programs.good = programs.good - 1 WHERE programs.id = OLD.program;
 END */;;
DELIMITER ;
/*!50003 SET SESSION SQL_MODE=@OLD_SQL_MODE */;


# テーブルのダンプ programs
# ------------------------------------------------------------

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
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# テーブルのダンプ users
# ------------------------------------------------------------

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `screenname` varchar(50) NOT NULL DEFAULT '',
  `name` varchar(50) NOT NULL,
  `profile` varchar(300) NOT NULL,
  `website` varchar(300) NOT NULL,
  `location` varchar(50) NOT NULL,
  `icon_url` varchar(140) NOT NULL,
  `token` varchar(140) NOT NULL,
  `secret` varchar(140) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
