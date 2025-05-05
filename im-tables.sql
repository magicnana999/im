-- IM 系统关系链数据库表设计（支持多租户）
-- 数据库：im
-- 字符集：utf8mb4（支持多语言和 emoji）
-- 存储引擎：InnoDB（支持事务和行级锁）

-- 用户表
CREATE TABLE IF NOT EXISTS im_user
(
    app_id        VARCHAR(50)     NOT NULL COMMENT '租户 ID',
    user_id       BIGINT UNSIGNED NOT NULL COMMENT '用户 ID，唯一标识',
    username      VARCHAR(50)     NOT NULL COMMENT '用户名',
    nickname      VARCHAR(50) COMMENT '昵称',
    phone_number  VARCHAR(20) COMMENT '手机号码（国际格式，如 +861234567890）',
    password_hash VARCHAR(256)    NOT NULL COMMENT '密码哈希',
    status        VARCHAR(20)     NOT NULL DEFAULT 'active' COMMENT '用户状态（active, inactive, banned）',
    created_at    TIMESTAMP                DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at    TIMESTAMP                DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (app_id, user_id) COMMENT '复合主键，支持多租户',
    UNIQUE INDEX idx_app_username (app_id, username) COMMENT '租户内用户名唯一索引',
    UNIQUE INDEX idx_app_phone_number (app_id, phone_number) COMMENT '租户内手机号码唯一索引'
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT ='用户基本信息表';

-- 好友关系表
CREATE TABLE IF NOT EXISTS im_friend_relation
(
    app_id     VARCHAR(50)     NOT NULL COMMENT '租户 ID',
    user_id    BIGINT UNSIGNED NOT NULL COMMENT '用户 ID',
    friend_id  BIGINT UNSIGNED NOT NULL COMMENT '好友 ID',
    group_id   BIGINT UNSIGNED DEFAULT 0 COMMENT '分组 ID，0 表示默认分组',
    remark     VARCHAR(50) COMMENT '好友备注',
    created_at TIMESTAMP       DEFAULT CURRENT_TIMESTAMP COMMENT '关系建立时间',
    PRIMARY KEY (app_id, user_id, friend_id) COMMENT '复合主键，支持多租户',
    INDEX idx_app_friend_id (app_id, friend_id) COMMENT '租户内反向查询索引'
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT ='双向好友关系表';

-- 黑名单表
CREATE TABLE IF NOT EXISTS im_blacklist
(
    app_id     VARCHAR(50)     NOT NULL COMMENT '租户 ID',
    user_id    BIGINT UNSIGNED NOT NULL COMMENT '用户 ID',
    blocked_id BIGINT UNSIGNED NOT NULL COMMENT '被拉黑用户 ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '拉黑时间',
    PRIMARY KEY (app_id, user_id, blocked_id) COMMENT '复合主键，支持多租户',
    INDEX idx_app_blocked_id (app_id, blocked_id) COMMENT '租户内反向查询索引'
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT ='黑名单表';

-- 好友分组表
CREATE TABLE IF NOT EXISTS im_friend_group
(
    app_id     VARCHAR(50)     NOT NULL COMMENT '租户 ID',
    group_id   BIGINT UNSIGNED NOT NULL COMMENT '分组 ID',
    user_id    BIGINT UNSIGNED NOT NULL COMMENT '用户 ID',
    group_name VARCHAR(50)     NOT NULL COMMENT '分组名称',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (app_id, group_id) COMMENT '复合主键，支持多租户',
    INDEX idx_app_user_id (app_id, user_id) COMMENT '租户内用户分组查询索引'
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT ='好友分组表';

-- 好友请求表
CREATE TABLE IF NOT EXISTS im_friend_request
(
    app_id       VARCHAR(50)     NOT NULL COMMENT '租户 ID',
    request_id   BIGINT UNSIGNED NOT NULL COMMENT '请求 ID',
    from_user_id BIGINT UNSIGNED NOT NULL COMMENT '发起者 ID',
    to_user_id   BIGINT UNSIGNED NOT NULL COMMENT '接收者 ID',
    status       VARCHAR(20)     NOT NULL DEFAULT 'pending' COMMENT '请求状态（pending, accepted, rejected）',
    message      VARCHAR(200) COMMENT '请求消息',
    created_at   TIMESTAMP                DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at   TIMESTAMP                DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (app_id, request_id) COMMENT '复合主键，支持多租户',
    INDEX idx_app_from_user_id (app_id, from_user_id) COMMENT '租户内发起者查询索引',
    INDEX idx_app_to_user_id (app_id, to_user_id) COMMENT '租户内接收者查询索引',
    UNIQUE INDEX idx_app_from_to (app_id, from_user_id, to_user_id) COMMENT '租户内防止重复请求'
    ) ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4 COMMENT ='好友请求表';
