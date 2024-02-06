# lanshan-chat

#### 聊天室
##### 接口文档
[https://console-docs.apipost.cn/preview/c31cb8fb5dffda9e/9eee74c277201ad3](https://console-docs.apipost.cn/preview/c31cb8fb5dffda9e/9eee74c277201ad3)

[状态码字典](https://github.com/jizizr/lanshan-chat/blob/master/app/api/internal/consts/error_id.go)
##### 基础功能

- [x] 用户注册
- [x] 用户登录
- [x] 用户个人主页
- [x] 用户信息修改
- [x] 好友功能 
- [x] 群聊功能 
- [x] 消息功能 

##### 加分项

- [x] 用户密码加盐加密
- [ ] 用户登录有短信登录、邮箱登录、第三方登录多种形式
- [ ] 验证码（登录，注册，修改密码）
- [x] 用户状态保存使用 JWT 或 Session
- [x] 群的权限（管理员 群主 相关，如 踢人 禁言）
- [x] 用户 和 群聊 的 模糊查询
- [x] 离线消息处理
- [x] 显示未读消息
- [x] 敏感词汇屏蔽
- [ ] 朋友圈

## 技术栈

[![img](https://github.com/StellarisW/douyin/raw/master/manifest/docs/image/mysql.svg)](https://github.com/StellarisW/douyin/blob/master/manifest/docs/image/mysql.svg)

- [mysql](https://www.mysql.com/)

> 一个关系型数据库管理系统，由瑞典MySQL AB 公司开发，属于 Oracle 旗下产品。MySQL 是最流行的关系型数据库管理系统关系型数据库管理系统之一，在 WEB 应用方面，MySQL是最好的 RDBMS (Relational Database Management System，关系数据库管理系统) 应用软件之一

[![img](https://github.com/StellarisW/douyin/raw/master/manifest/docs/image/redis.svg)](https://github.com/StellarisW/douyin/blob/master/manifest/docs/image/redis.svg)

- [redis](https://redis.io/)

> 一个开源的、使用C语言编写的、支持网络交互的、可基于内存也可持久化的Key-Value数据库。

缓存存储还是选型最普遍的redis

<img width="159px" src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png">


- [gin](https://gin-gonic.com/)

> 一个用Go（Golang）编写的高性能Web框架，它提供了一套简洁而强大的API，使得开发者能够快速构建高效的Web应用和微服务。Gin特别注重性能，通过减少内存分配来达到快速响应，使其成为构建高性能应用的理想选择。
>
> 它具备以下特点：高性能（优化的路由器，快速处理请求，低延迟），易扩展（支持中间件，方便添加功能和构建复杂应用，高度可定制），低门槛（简洁的API设计，文档齐全，容易上手，通过gin-cli快速生成项目骨架）。

同时，Gin也是目前极受欢迎的Go Web框架之一，所以本项目采用gin为主框架搭建后端



## 中间件设计

- #### JWT中间件

- #### CORS中间件

## 组件和模块设计

整体功能分为**用户、群组、消息**三个主要模块，相互配合。

对于用户，采用JWT技术进行用户认证，保证用户安全。

对于群组，分为公有和私有两类群组，公有群组可以通过搜索查询加入，而私有群组只能通过群组管理员生成的token加入，保证了用户对于私密群组的需求。

对于消息，发送消息时会先进行敏感词审查，如发现敏感词，则会上报审核。管理员可以对群聊中的用户进行禁言操作，对于用户消息也可以删除，用户可以对自己的消息进行修改。同时提供ws和http两种获取消息接口，前者用于用户在线时实时获取消息，后者则用于用户上线时批量获取离线时的消息。ws 的实时消息采用写扩散的设计，通过全局哈希表，每个群组维护一个在线用户channel组成的链表。发送消息后创建一个协程遍历并发送消息给群组在线用户。由于是写扩散，用户建立ws后会创建一个channel，放入所有已加入群组的链表中，然后从channel获取消息即可。

本系统中，采用MySQL进行数据存储，redis作数据缓存。
##### 数据库设计
- 用户表

`users`
```sql
-- auto-generated definition
create table users
(
    user_id   int auto_increment
        primary key,
    username  varchar(255)                        not null,
    nickname  varchar(255)                        not null,
    password  char(32)                            not null,
    email     varchar(254)                        not null,
    profile   text                                null,
    joined_at timestamp default CURRENT_TIMESTAMP null,
    constraint email
        unique (email),
    constraint username
        unique (username)
);
```
- 群组表

`groups`
```sql
-- auto-generated definition
create table `groups`
(
    group_id    int auto_increment
        primary key,
    group_name  varchar(255)                           not null,
    avatar      varchar(255)                           null,
    description text                                   not null,
    type        enum ('public', 'private', 'personal') not null,
    created_at  datetime                               not null
);
```

`user_groups`
```sql
-- auto-generated definition
create table user_groups
(
    id          int auto_increment
        primary key,
    user_id     int                                                      not null,
    group_id    int                                                      not null,
    last_read   int                            default 0                 null,
    role        enum ('admin', 'member')       default 'member'          not null,
    status      enum ('ok', 'muted', 'banned') default 'ok'              not null,
    muted_until timestamp                                                null,
    joined_at   timestamp                      default CURRENT_TIMESTAMP not null
);

create index group_id
    on user_groups (group_id);

create index user_id
    on user_groups (user_id);
```
- 消息表

`group_messages`
```sql
-- auto-generated definition
create table group_message
(
    id               int auto_increment
        primary key,
    group_id         int                                              not null,
    message_id       int                                              not null,
    sender_id        int                                              not null,
    reply_message_id int                                              null,
    message          text                                             null,
    type             enum ('text', 'image', 'video', 'audio', 'file') not null,
    url              varchar(255)                                     null,
    file_name        varchar(255)                                     null,
    send_date        datetime                                         not null,
    update_date      datetime                                         null,
    constraint group_message_unique
        unique (group_id, message_id)
);
```
##### 目录结构
```
├─app
│  └─api
│      ├─global
│      │  ├─config
│      │  └─wsmap
│      ├─internal
│      │  ├─consts
│      │  ├─controller
│      │  ├─dao
│      │  │  ├─mysql
│      │  │  └─redis
│      │  ├─initialize
│      │  ├─middleware
│      │  ├─model
│      │  └─service
│      └─router
├─development
├─manifest
├─tmp
└─utils
```
