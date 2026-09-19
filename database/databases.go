package database

//数据库表结构

//数据库：at
//msg_id message_id
//reply_id comment_a_id 要回复的评论ID内容
//comment_roo_id 根评论
//link_id 帖子id
//reply_user_id 要回复的用户的id
//comment_text 回复内容

var DatabaseList = []string{dataBase_AT, dataBase_Config}

const dataBase_AT = `
	CREATE TABLE IF NOT EXISTS at (
	msg_id BIGINT PRIMARY KEY,
	reply_id BIGINT,
	comment_root_id BIGINT,
	link_id BIGINT,
	reply_user_id BIGINT,
	comment_text TEXT,
	reply boolean 
	);
`

//使用INT作为配置
const dataBase_Config = `
CREATE TABLE IF NOT EXISTS cfg (
   name INT PRIMARY KEY,
   value TEXT,
   enable BOLLEAN
)
`
