CREATE TABLE `tbl_file` (
    `id` INT(11) NOT NULL AUTO_INCREMENT,
    `file_sha1` CHAR(40) NOT NULL DEFAULT '' COMMENT '文件hash',
    `file_name` VARCHAR(256) NOT NULL DEFAULT '' COMMENT '文件名',
    `file_size` BIGINT(20) NOT NULL DEFAULT 0 COMMENT '文件大小',
    `file_addr` VARCHAR(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
    `create_at` DATETIME DEFAULT NOW() COMMENT '创建日期',
    `update_at` DATETIME DEFAULT NOW() ON UPDATE CURRENT_TIMESTAMP() COMMENT '更新日期',
    `status` TINYINT NOT NULL DEFAULT 0 COMMENT '状态(可用/禁用/已删除等)',
    `ext1` INT(11) NOT NULL DEFAULT 0 COMMENT '备用字段1',
    `ext2` TEXT COMMENT '备用字段2',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_file_hash` (`file_sha1`),
    KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;