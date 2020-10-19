/*
 Navicat Premium Data Transfer

 Source Server         : ll
 Source Server Type    : MySQL
 Source Server Version : 50642
 Source Host           : 47.107.83.200:3306
 Source Schema         : code_evaluation

 Target Server Type    : MySQL
 Target Server Version : 50642
 File Encoding         : 65001

 Date: 19/10/2020 14:58:37
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for eva_user
-- ----------------------------
DROP TABLE IF EXISTS `eva_user`;
CREATE TABLE `eva_user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户唯一标识Id',
  `user_account` varchar(20) CHARACTER SET utf8 NOT NULL COMMENT '用户账号',
  `user_password` varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '用户密码',
  `user_nick` varchar(100) CHARACTER SET utf8 NOT NULL COMMENT '用户昵称',
  `user_status` varchar(20) CHARACTER SET utf8 NOT NULL DEFAULT 'common' COMMENT '用户权限',
  `salt` varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '盐',
  `user_portrait` varchar(255) CHARACTER SET utf8 DEFAULT NULL COMMENT '用户头像存储路径',
  PRIMARY KEY (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=111 DEFAULT CHARSET=latin1;

SET FOREIGN_KEY_CHECKS = 1;




/*
 Navicat Premium Data Transfer

 Source Server         : ll
 Source Server Type    : MySQL
 Source Server Version : 50642
 Source Host           : 47.107.83.200:3306
 Source Schema         : code_evaluation

 Target Server Type    : MySQL
 Target Server Version : 50642
 File Encoding         : 65001

 Date: 19/10/2020 14:58:22
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for eva_code
-- ----------------------------
DROP TABLE IF EXISTS `eva_code`;
CREATE TABLE `eva_code` (
  `code_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '代码唯一标识id',
  `code_name` varchar(50) CHARACTER SET utf8 NOT NULL COMMENT '代码名称',
  `user_id` int(11) NOT NULL COMMENT '关联user表中user_id',
  `code_text` text CHARACTER SET utf8 COMMENT '代码',
  `path` varchar(255) CHARACTER SET utf8 NOT NULL COMMENT '存储路径',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更改时间',
  PRIMARY KEY (`code_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=104 DEFAULT CHARSET=latin1;

SET FOREIGN_KEY_CHECKS = 1;




/*
 Navicat Premium Data Transfer

 Source Server         : ll
 Source Server Type    : MySQL
 Source Server Version : 50642
 Source Host           : 47.107.83.200:3306
 Source Schema         : code_evaluation

 Target Server Type    : MySQL
 Target Server Version : 50642
 File Encoding         : 65001

 Date: 19/10/2020 14:58:30
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for eva_compare
-- ----------------------------
DROP TABLE IF EXISTS `eva_compare`;
CREATE TABLE `eva_compare` (
  `compare_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '对拍唯一标识',
  `compare_name` varchar(50) NOT NULL COMMENT '对拍名称',
  `user_id` int(11) NOT NULL COMMENT '关联user表中user_id',
  `first_code_id` int(11) NOT NULL COMMENT '对拍中代码1的id',
  `second_code_id` int(11) NOT NULL COMMENT '对拍中代码2的id',
  `input_data_path` varchar(255) NOT NULL COMMENT '随机生成数据代码\r\n随机生成数据代码的路径',
  `max_input_group` int(11) DEFAULT '1' COMMENT '最大生成数据组数 ',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '修改时间',
  `remarks` text COMMENT '备注',
  PRIMARY KEY (`compare_id`)
) ENGINE=InnoDB AUTO_INCREMENT=89 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;




/*
 Navicat Premium Data Transfer

 Source Server         : ll
 Source Server Type    : MySQL
 Source Server Version : 50642
 Source Host           : 47.107.83.200:3306
 Source Schema         : code_evaluation

 Target Server Type    : MySQL
 Target Server Version : 50642
 File Encoding         : 65001

 Date: 19/10/2020 14:58:44
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for run_code
-- ----------------------------
DROP TABLE IF EXISTS `run_code`;
CREATE TABLE `run_code` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `code_id` int(11) NOT NULL COMMENT '代码id',
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `input_path` text COMMENT '输入数据存放的地方',
  `run_data` text CHARACTER SET utf8mb4 COMMENT '代码运行结果',
  `msg_data` text COMMENT '代码运行信息',
  `select_flag` text COMMENT '查询结果id',
  `version` int(11) NOT NULL DEFAULT '0' COMMENT '版本，0为没有运行过，1为运行过',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;




/*
 Navicat Premium Data Transfer

 Source Server         : ll
 Source Server Type    : MySQL
 Source Server Version : 50642
 Source Host           : 47.107.83.200:3306
 Source Schema         : code_evaluation

 Target Server Type    : MySQL
 Target Server Version : 50642
 File Encoding         : 65001

 Date: 19/10/2020 14:58:51
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for run_compare
-- ----------------------------
DROP TABLE IF EXISTS `run_compare`;
CREATE TABLE `run_compare` (
  `compare_id` int(11) NOT NULL COMMENT '对拍id，为什么要和对拍分个表？因为我在多进程读数据库时会有排它锁，防止因为排它锁导致对拍表的正常读取数据',
  `user_id` int(11) NOT NULL COMMENT '用户id',
  `compare_data` text COMMENT '对拍结果',
  `version` int(11) NOT NULL DEFAULT '0' COMMENT '版本，0为没有对拍过，1为对拍过',
  PRIMARY KEY (`compare_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
