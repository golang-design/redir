CREATE TABLE `collink` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `alias` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `kind` tinyint(1) NOT NULL DEFAULT '0',
    `url` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `private` tinyint(1) NOT NULL DEFAULT '0',
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    `is_deleted` tinyint(4) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uidx_alias` (`alias`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;