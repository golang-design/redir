-- Copyright 2021 The golang.design Initiative Authors.
-- All rights reserved. Use of this source code is governed
-- by a MIT license that can be found in the LICENSE file.
--
-- Originally written by Mai Yang <maiyang.me>.

CREATE TABLE `visit` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
    `alias` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
    `kind` tinyint(1) NOT NULL DEFAULT '0',
    `ip` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `ua` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `referer` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;