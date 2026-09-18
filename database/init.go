package database

//数据库，新增的数据库都要支持以下操作

//数据库：at
//msg_id message_id
//reply_id comment_a_id 回复的评论ID内容
//comment_roo_id 根评论
const DataBase_AT = `
	CREATE TABLE IF NOT EXISTS at (
	msg_id BIGINT PRIMARY KEY,
	reply_id BIGINT,
	comment_root_id BIGINT,
	link_id BIGINT,
	reply_user_id BIGINT,
	comment_text TEXT,
	reply boolean
`

var Init = func() {
	panic("数据库未正常配置")
}
