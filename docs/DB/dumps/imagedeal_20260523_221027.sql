-- MySQL dump 10.13  Distrib 9.6.0, for macos26.2 (arm64)
--
-- Host: 127.0.0.1    Database: imagedeal
-- ------------------------------------------------------
-- Server version	8.0.36

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
-- Current Database: `imagedeal`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `imagedeal` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;

USE `imagedeal`;

--
-- Table structure for table `image_versions`
--

DROP TABLE IF EXISTS `image_versions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `image_versions` (
  `id` char(36) NOT NULL,
  `image_id` char(36) NOT NULL,
  `version_no` int NOT NULL,
  `file_path` varchar(500) DEFAULT '',
  `layer_data` json DEFAULT NULL,
  `is_current` tinyint(1) NOT NULL DEFAULT '0',
  `created_by` char(36) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_image_version` (`image_id`,`version_no`),
  KEY `idx_versions_image` (`image_id`),
  CONSTRAINT `fk_versions_image` FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `image_versions`
--

LOCK TABLES `image_versions` WRITE;
/*!40000 ALTER TABLE `image_versions` DISABLE KEYS */;
INSERT INTO `image_versions` VALUES ('324b5096-dd1b-4651-a421-c32f469fbf38','8d9334d3-af5d-4035-8dc3-5af3b21dded3',1,'','null',1,'00000000-0000-0000-0000-000000000011','2026-05-23 20:29:48.084');
/*!40000 ALTER TABLE `image_versions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `images`
--

DROP TABLE IF EXISTS `images`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `images` (
  `id` char(36) NOT NULL,
  `tenant_id` char(36) NOT NULL,
  `global_no` varchar(50) DEFAULT '',
  `image_url` text,
  `storage_path` varchar(500) DEFAULT '',
  `thumb_path` varchar(500) DEFAULT '',
  `category` varchar(100) DEFAULT '',
  `status` varchar(30) NOT NULL DEFAULT 'pending_storage',
  `assigned_to` char(36) DEFAULT NULL,
  `import_batch_id` char(36) DEFAULT NULL,
  `external_id` varchar(100) DEFAULT '',
  `metadata` json DEFAULT NULL,
  `discard_reason` text,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_images_claim` (`tenant_id`,`status`,`category`,`created_at`),
  KEY `idx_images_tasks` (`tenant_id`,`assigned_to`,`status`),
  KEY `idx_images_review` (`tenant_id`,`status`),
  KEY `idx_images_global_no` (`global_no`),
  CONSTRAINT `fk_images_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `images`
--

LOCK TABLES `images` WRITE;
/*!40000 ALTER TABLE `images` DISABLE KEYS */;
INSERT INTO `images` VALUES ('12635cd6-d0c9-4024-a8d4-ac40a492e7d9','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/social-02/1080/1080','originals/00000000-0000-0000-0000-000000000001/12635cd6-d0c9-4024-a8d4-ac40a492e7d9.jpg','thumbs/00000000-0000-0000-0000-000000000001/12635cd6-d0c9-4024-a8d4-ac40a492e7d9.jpg','social','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.818','2026-05-23 20:45:14.728'),('4c2f2789-6081-4205-b6ce-7601836875dd','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/banner-03/1200/400','originals/00000000-0000-0000-0000-000000000001/4c2f2789-6081-4205-b6ce-7601836875dd.jpg','thumbs/00000000-0000-0000-0000-000000000001/4c2f2789-6081-4205-b6ce-7601836875dd.jpg','banner','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.794','2026-05-23 20:44:52.523'),('4fb60da9-46a1-4597-97bd-52181a94a595','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/poster-03/600/900','originals/00000000-0000-0000-0000-000000000001/4fb60da9-46a1-4597-97bd-52181a94a595.jpg','thumbs/00000000-0000-0000-0000-000000000001/4fb60da9-46a1-4597-97bd-52181a94a595.jpg','poster','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.806','2026-05-23 20:45:02.827'),('6543b3b6-a8b8-4801-b185-1130c519953d','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/banner-02/1200/400','originals/00000000-0000-0000-0000-000000000001/6543b3b6-a8b8-4801-b185-1130c519953d.jpg','thumbs/00000000-0000-0000-0000-000000000001/6543b3b6-a8b8-4801-b185-1130c519953d.jpg','banner','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.790','2026-05-23 20:44:52.045'),('7d2c0611-b3ba-4860-8920-d2952154ffd7','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/poster-01/600/900','originals/00000000-0000-0000-0000-000000000001/7d2c0611-b3ba-4860-8920-d2952154ffd7.jpg','thumbs/00000000-0000-0000-0000-000000000001/7d2c0611-b3ba-4860-8920-d2952154ffd7.jpg','poster','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.801','2026-05-23 20:44:59.238'),('7f520b2e-1378-402e-8a35-2fc0a4e69582','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/icon-02/512/512','originals/00000000-0000-0000-0000-000000000001/7f520b2e-1378-402e-8a35-2fc0a4e69582.jpg','thumbs/00000000-0000-0000-0000-000000000001/7f520b2e-1378-402e-8a35-2fc0a4e69582.jpg','icon','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.813','2026-05-23 20:45:06.822'),('854300f0-bdc3-4411-bb92-61189a8643ff','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/social-01/1080/1080','originals/00000000-0000-0000-0000-000000000001/854300f0-bdc3-4411-bb92-61189a8643ff.jpg','thumbs/00000000-0000-0000-0000-000000000001/854300f0-bdc3-4411-bb92-61189a8643ff.jpg','social','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.816','2026-05-23 20:45:10.420'),('8d9334d3-af5d-4035-8dc3-5af3b21dded3','00000000-0000-0000-0000-000000000001','ACME-20260523-000001','https://picsum.photos/800/600?random=1','originals/00000000-0000-0000-0000-000000000001/8d9334d3-af5d-4035-8dc3-5af3b21dded3.jpg','thumbs/00000000-0000-0000-0000-000000000001/8d9334d3-af5d-4035-8dc3-5af3b21dded3.jpg','banner','completed','00000000-0000-0000-0000-000000000011','1c3d456d-aa32-4c8c-acb2-f9ba2bea1826','',NULL,'','2026-05-23 20:29:25.342','2026-05-23 20:30:00.960'),('a743742e-f38a-4f04-b58c-d1890b6e7988','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/banner-01/1200/400','originals/00000000-0000-0000-0000-000000000001/a743742e-f38a-4f04-b58c-d1890b6e7988.jpg','thumbs/00000000-0000-0000-0000-000000000001/a743742e-f38a-4f04-b58c-d1890b6e7988.jpg','banner','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.784','2026-05-23 20:44:49.000'),('a82d4f0c-ffd0-4575-b63b-487562e5b367','00000000-0000-0000-0000-000000000001','','https://picsum.photos/800/600?random=2','originals/00000000-0000-0000-0000-000000000001/a82d4f0c-ffd0-4575-b63b-487562e5b367.jpg','thumbs/00000000-0000-0000-0000-000000000001/a82d4f0c-ffd0-4575-b63b-487562e5b367.jpg','banner','pending_assign','','1c3d456d-aa32-4c8c-acb2-f9ba2bea1826','',NULL,'','2026-05-23 20:29:25.346','2026-05-23 20:29:33.236'),('a98c3a78-676c-4c77-b8f0-058a2037881c','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/icon-03/512/512','originals/00000000-0000-0000-0000-000000000001/a98c3a78-676c-4c77-b8f0-058a2037881c.jpg','thumbs/00000000-0000-0000-0000-000000000001/a98c3a78-676c-4c77-b8f0-058a2037881c.jpg','icon','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.815','2026-05-23 20:45:09.981'),('c1c7c3b8-b3f9-40e3-8a74-1efdc9cd0f64','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/banner-05/1200/400','originals/00000000-0000-0000-0000-000000000001/c1c7c3b8-b3f9-40e3-8a74-1efdc9cd0f64.jpg','thumbs/00000000-0000-0000-0000-000000000001/c1c7c3b8-b3f9-40e3-8a74-1efdc9cd0f64.jpg','banner','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.799','2026-05-23 20:44:56.166'),('c45e222b-0f4f-4a05-98b3-35b588732912','00000000-0000-0000-0000-000000000001','','https://picsum.photos/800/600?random=3','originals/00000000-0000-0000-0000-000000000001/c45e222b-0f4f-4a05-98b3-35b588732912.jpg','thumbs/00000000-0000-0000-0000-000000000001/c45e222b-0f4f-4a05-98b3-35b588732912.jpg','poster','pending_assign','','1c3d456d-aa32-4c8c-acb2-f9ba2bea1826','',NULL,'','2026-05-23 20:29:25.348','2026-05-23 20:29:36.850'),('c46b90f1-4f64-4f50-b554-a11027994b17','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/social-03/1080/1080','originals/00000000-0000-0000-0000-000000000001/c46b90f1-4f64-4f50-b554-a11027994b17.jpg','thumbs/00000000-0000-0000-0000-000000000001/c46b90f1-4f64-4f50-b554-a11027994b17.jpg','social','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.820','2026-05-23 20:45:13.982'),('d1bd7982-280a-4c37-8b81-3e503327046e','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/poster-04/600/900','originals/00000000-0000-0000-0000-000000000001/d1bd7982-280a-4c37-8b81-3e503327046e.jpg','thumbs/00000000-0000-0000-0000-000000000001/d1bd7982-280a-4c37-8b81-3e503327046e.jpg','poster','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.809','2026-05-23 20:45:03.312'),('d1f20932-46a0-41a1-b361-e48b526918d2','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/poster-02/600/900','originals/00000000-0000-0000-0000-000000000001/d1f20932-46a0-41a1-b361-e48b526918d2.jpg','thumbs/00000000-0000-0000-0000-000000000001/d1f20932-46a0-41a1-b361-e48b526918d2.jpg','poster','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.804','2026-05-23 20:44:59.790'),('e22a97ca-c084-4fa6-af14-c2cbe7f2d41c','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/banner-04/1200/400','originals/00000000-0000-0000-0000-000000000001/e22a97ca-c084-4fa6-af14-c2cbe7f2d41c.jpg','thumbs/00000000-0000-0000-0000-000000000001/e22a97ca-c084-4fa6-af14-c2cbe7f2d41c.jpg','banner','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.796','2026-05-23 20:44:55.624'),('f39eabad-9ff7-4824-a3ae-951aa2672d60','00000000-0000-0000-0000-000000000001','','https://picsum.photos/seed/icon-01/512/512','originals/00000000-0000-0000-0000-000000000001/f39eabad-9ff7-4824-a3ae-951aa2672d60.jpg','thumbs/00000000-0000-0000-0000-000000000001/f39eabad-9ff7-4824-a3ae-951aa2672d60.jpg','icon','pending_assign','','3a07583e-0880-4421-96be-12a46df8f0da','',NULL,'','2026-05-23 20:44:47.811','2026-05-23 20:45:06.411');
/*!40000 ALTER TABLE `images` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `import_batches`
--

DROP TABLE IF EXISTS `import_batches`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `import_batches` (
  `id` char(36) NOT NULL,
  `tenant_id` char(36) NOT NULL,
  `file_name` varchar(255) NOT NULL,
  `storage_path` varchar(500) DEFAULT '',
  `status` varchar(30) NOT NULL DEFAULT 'processing',
  `total_rows` int NOT NULL DEFAULT '0',
  `success_count` int NOT NULL DEFAULT '0',
  `failed_count` int NOT NULL DEFAULT '0',
  `error_log` text,
  `created_by` char(36) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_batches_tenant` (`tenant_id`),
  CONSTRAINT `fk_batches_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `import_batches`
--

LOCK TABLES `import_batches` WRITE;
/*!40000 ALTER TABLE `import_batches` DISABLE KEYS */;
INSERT INTO `import_batches` VALUES ('1c3d456d-aa32-4c8c-acb2-f9ba2bea1826','00000000-0000-0000-0000-000000000001','import_sample.csv','imports/00000000-0000-0000-0000-000000000001/1c3d456d-aa32-4c8c-acb2-f9ba2bea1826/import_sample.csv','completed',3,3,0,'','00000000-0000-0000-0000-000000000010','2026-05-23 20:29:25.335','2026-05-23 20:29:36.858'),('3a07583e-0880-4421-96be-12a46df8f0da','00000000-0000-0000-0000-000000000001','import_batch_medium.csv','imports/00000000-0000-0000-0000-000000000001/3a07583e-0880-4421-96be-12a46df8f0da/import_batch_medium.csv','completed',15,15,0,'','00000000-0000-0000-0000-000000000010','2026-05-23 20:44:46.886','2026-05-23 20:45:14.734');
/*!40000 ALTER TABLE `import_batches` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `review_records`
--

DROP TABLE IF EXISTS `review_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `review_records` (
  `id` char(36) NOT NULL,
  `image_id` char(36) NOT NULL,
  `version_id` char(36) DEFAULT NULL,
  `round` int NOT NULL,
  `result` varchar(20) NOT NULL,
  `comment` text,
  `reviewer_id` char(36) NOT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_review_image` (`image_id`),
  CONSTRAINT `fk_review_image` FOREIGN KEY (`image_id`) REFERENCES `images` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `review_records`
--

LOCK TABLES `review_records` WRITE;
/*!40000 ALTER TABLE `review_records` DISABLE KEYS */;
INSERT INTO `review_records` VALUES ('167e7cd3-7fef-43f1-8019-ed726e1eef86','8d9334d3-af5d-4035-8dc3-5af3b21dded3','324b5096-dd1b-4651-a421-c32f469fbf38',1,'pass','ok','00000000-0000-0000-0000-000000000012','2026-05-23 20:29:55.859'),('38c98452-ffac-40cc-a2e0-98b0835270ca','8d9334d3-af5d-4035-8dc3-5af3b21dded3','324b5096-dd1b-4651-a421-c32f469fbf38',2,'pass','final ok','00000000-0000-0000-0000-000000000012','2026-05-23 20:30:00.911');
/*!40000 ALTER TABLE `review_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tenants`
--

DROP TABLE IF EXISTS `tenants`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tenants` (
  `id` char(36) NOT NULL,
  `name` varchar(200) NOT NULL,
  `slug` varchar(50) NOT NULL,
  `settings` json DEFAULT NULL,
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tenants`
--

LOCK TABLES `tenants` WRITE;
/*!40000 ALTER TABLE `tenants` DISABLE KEYS */;
INSERT INTO `tenants` VALUES ('00000000-0000-0000-0000-000000000001','ACME Demo','ACME','{\"max_discard_per_day\": 10, \"second_review_enabled\": true}','2026-05-23 12:27:21.767','2026-05-23 12:27:21.767');
/*!40000 ALTER TABLE `tenants` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` char(36) NOT NULL,
  `tenant_id` char(36) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL,
  `display_name` varchar(100) DEFAULT '',
  `role` varchar(20) NOT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_tenant_email` (`tenant_id`,`email`),
  KEY `idx_users_tenant` (`tenant_id`),
  CONSTRAINT `fk_users_tenant` FOREIGN KEY (`tenant_id`) REFERENCES `tenants` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES ('00000000-0000-0000-0000-000000000010','00000000-0000-0000-0000-000000000001','admin@example.com','$2a$10$8MuzVH8wzjTu8KUz6srHVurc//pTtDgz5ov0JnrfMkSRJ8UIWHsG6','管理员','admin',1,'2026-05-23 12:27:21.768','2026-05-23 12:27:40.257'),('00000000-0000-0000-0000-000000000011','00000000-0000-0000-0000-000000000001','user@example.com','$2a$10$8MuzVH8wzjTu8KUz6srHVurc//pTtDgz5ov0JnrfMkSRJ8UIWHsG6','张三','user',1,'2026-05-23 12:27:21.768','2026-05-23 12:27:40.257'),('00000000-0000-0000-0000-000000000012','00000000-0000-0000-0000-000000000001','reviewer@example.com','$2a$10$8MuzVH8wzjTu8KUz6srHVurc//pTtDgz5ov0JnrfMkSRJ8UIWHsG6','审核员','reviewer',1,'2026-05-23 12:27:21.768','2026-05-23 12:27:40.257');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'imagedeal'
--

--
-- Dumping routines for database 'imagedeal'
--
--
-- WARNING: can't read the INFORMATION_SCHEMA.libraries table. It's most probably an old server 8.0.36.
--
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-05-23 22:10:28
