-- Copyright 2021 The golang.design Initiative Authors.
-- All rights reserved. Use of this source code is governed
-- by a MIT license that can be found in the LICENSE file.
--
-- Originally written by Mai Yang <maiyang.me>.

CREATE TABLE `collink` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `alias` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `url` varchar(1024) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `created_at` datetime DEFAULT NULL,
    `updated_at` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uidx_alias` (`alias`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;