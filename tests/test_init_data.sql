CREATE database if NOT EXISTS douyin_test;
use douyin_test;

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `is_deleted` tinyint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否被删除，0 为未删除，1 为已删除',
  `video_id` bigint(0) UNSIGNED NOT NULL COMMENT '视频主键ID',
  `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '评论者用户主键ID',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '评论内容',
  `created_time` datetime(3) NOT NULL COMMENT '评论时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_userid`(`user_id`) USING BTREE,
  INDEX `idx_videoid`(`video_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of comment
-- ----------------------------

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `to_user_id` bigint(0) UNSIGNED NOT NULL COMMENT '接收者用户ID',
  `from_user_id` bigint(0) UNSIGNED NOT NULL COMMENT '发送者用户ID',
  `content` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '消息内容',
  `create_time` datetime(3) NOT NULL COMMENT '消息发送时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_userid_to`(`to_user_id`) USING BTREE,
  INDEX `idx_userid_from`(`from_user_id`) USING BTREE,
  INDEX `idx_userid`(`to_user_id`, `from_user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of message
-- ----------------------------

-- ----------------------------
-- Table structure for relation
-- ----------------------------
DROP TABLE IF EXISTS `relation`;
CREATE TABLE `relation`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `is_deleted` tinyint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否被删除，0 为未删除，1 为已删除',
  `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '粉丝的用户ID',
  `to_user_id` bigint(0) UNSIGNED NOT NULL COMMENT '被关注者的用户ID',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_userid`(`user_id`, `to_user_id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_to_user_id`(`to_user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of relation
-- ----------------------------

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `username` varchar(63) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户名，也是用户的昵称',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '用户密码',
  `following_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '关注数',
  `follower_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '粉丝数',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/avatar/avatar.png' COMMENT '头像URL',
  `background_image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT 'https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/background/background.png' COMMENT '用户个人页顶部大图',
  `signature` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '这是一个个人简介' COMMENT '个人简介',
  `total_favorited` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '获赞数量',
  `work_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '作品数量',
  `favorite_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '点赞视频数量',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_username`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

DROP TABLE IF EXISTS `user_favorite_video`;
CREATE TABLE `user_favorite_video`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `user_id` bigint(0) UNSIGNED NOT NULL COMMENT '点赞者的用户ID',
  `video_id` bigint(0) UNSIGNED NOT NULL COMMENT '视频ID',
  `is_deleted` tinyint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否被删除，0 为未删除，1 为已删除',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_userid`(`user_id`, `video_id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE,
  INDEX `idx_video_id`(`video_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of user_favorite_video
-- ----------------------------

-- ----------------------------
-- Table structure for video
-- ----------------------------
DROP TABLE IF EXISTS `video`;
CREATE TABLE `video`  (
  `id` bigint(0) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `publish_time` datetime(3) NOT NULL COMMENT '视频发布时间',
  `author_id` bigint(0) UNSIGNED NOT NULL COMMENT '作者用户ID',
  `play_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '视频播放URL',
  `cover_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '封面图片URL',
  `favorite_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '视频点赞数',
  `comment_count` bigint(0) UNSIGNED NOT NULL DEFAULT 0 COMMENT '视频评论数',
  `title` varchar(63) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '视频标题',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_authorid`(`author_id`) USING BTREE,
  INDEX `idx_publish`(`publish_time`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of video
-- ----------------------------
INSERT INTO `video` VALUES (1, '2023-03-01 23:16:34.196', 1, 'https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/video/22d8cdc6-235f-4241-8ca4-bc2e970de10d.mp4', 'https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/db0273b3-591b-40d8-b580-544411ee7922.jpg', 0, 0, '示例视频1');
INSERT INTO `video` VALUES (2, '2023-03-01 23:17:22.813', 1, 'https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/video/436d6bf2-9377-42a4-a634-c6e2ff2fa640.mp4', 'https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/983566ec-0021-410b-814f-969117b23b48.jpg', 0, 0, '示例视频2');
INSERT INTO `video` VALUES (3, '2023-03-01 23:17:33.453', 1, 'https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/video/6b9b593f-272d-4b1a-b1c1-aced255d5f32.mp4', 'https://picture-bucket-01.oss-cn-beijing.aliyuncs.com/DouYin/cover/81234002-d1e8-4ed6-abda-9359937032b5.jpg', 0, 0, '示例视频3');
