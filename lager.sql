-- MySQL dump 10.13  Distrib 8.0.29, for Linux (x86_64)
--
-- Host: localhost    Database: lager
-- ------------------------------------------------------
-- Server version	8.0.29-0ubuntu0.20.04.3

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `lager`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `lager` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `lager`;

--
-- Table structure for table `bestand`
--

DROP TABLE IF EXISTS `bestand`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bestand` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT,
  `gertyp` tinyint unsigned NOT NULL,
  `modell` mediumint unsigned NOT NULL,
  `seriennummer` char(50) NOT NULL,
  `zinfo` mediumint unsigned DEFAULT NULL,
  `ticketnr` char(20) DEFAULT NULL,
  `ausgabename` varchar(20) DEFAULT NULL,
  `ausgabedatum` date DEFAULT NULL,
  `changed` tinyint unsigned DEFAULT NULL,
  `einkaufsdatum` date DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bestand`
--

LOCK TABLES `bestand` WRITE;
/*!40000 ALTER TABLE `bestand` DISABLE KEYS */;
/*!40000 ALTER TABLE `bestand` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gertyp`
--

DROP TABLE IF EXISTS `gertyp`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gertyp` (
  `id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `gertyp` char(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gertyp`
--

LOCK TABLES `gertyp` WRITE;
/*!40000 ALTER TABLE `gertyp` DISABLE KEYS */;
/*!40000 ALTER TABLE `gertyp` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `login`
--

DROP TABLE IF EXISTS `login`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `login` (
  `id` tinyint unsigned NOT NULL AUTO_INCREMENT,
  `username` char(20) NOT NULL,
  `password` char(32) NOT NULL,
  `rechte` tinyint unsigned DEFAULT NULL,
  `aktiv` tinyint unsigned NOT NULL,
  `session` char(255) DEFAULT NULL,
  `sessiontime` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `login`
--

LOCK TABLES `login` WRITE;
/*!40000 ALTER TABLE `login` DISABLE KEYS */;
INSERT INTO `login` VALUES (1,'admin','02acb3a105a6d5a224d6164df9e1642d',1,1,'','');
/*!40000 ALTER TABLE `login` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `modell`
--

DROP TABLE IF EXISTS `modell`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `modell` (
  `id` int NOT NULL AUTO_INCREMENT,
  `modell` char(20) NOT NULL,
  `gertyp` tinyint unsigned NOT NULL,
  `sperrbestand` tinyint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `modell`
--

LOCK TABLES `modell` WRITE;
/*!40000 ALTER TABLE `modell` DISABLE KEYS */;
/*!40000 ALTER TABLE `modell` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zinfo`
--

DROP TABLE IF EXISTS `zinfo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `zinfo` (
  `id` smallint unsigned NOT NULL AUTO_INCREMENT,
  `gertyp` tinyint unsigned NOT NULL,
  `zinfoname` char(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zinfo`
--

LOCK TABLES `zinfo` WRITE;
/*!40000 ALTER TABLE `zinfo` DISABLE KEYS */;
/*!40000 ALTER TABLE `zinfo` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `zinfodata`
--

DROP TABLE IF EXISTS `zinfodata`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `zinfodata` (
  `id` mediumint unsigned NOT NULL AUTO_INCREMENT,
  `zinfoid` mediumint unsigned NOT NULL,
  `bestandid` mediumint unsigned NOT NULL,
  `daten` char(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `zinfodata`
--

LOCK TABLES `zinfodata` WRITE;
/*!40000 ALTER TABLE `zinfodata` DISABLE KEYS */;
/*!40000 ALTER TABLE `zinfodata` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-06-01 19:04:02
