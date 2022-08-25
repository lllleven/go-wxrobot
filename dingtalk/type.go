package dingtalk

var (
	MessageTypeText     = "text"
	MessageTypeLink     = "link"
	MessageTypeMarkdown = "markdown"

	// 以下需要再接入
	MessageTypeActionCard = "actionCard"
	MessageTypeFeedCard   = "feedCard"
)

type DingTalkRobot struct {
	AccessToken string `json:"access_token"`
	Secret      string `json:"secret"`
}

// BaseResp
// 消息内容中不包含任何关键词
//{
//  "errcode":310000,
//  "errmsg":"keywords not in content"
//}
//
// timestamp 无效
//{
//  "errcode":310000,
//  "errmsg":"invalid timestamp"
//}
//
// 签名不匹配
//{
//  "errcode":310000,
//  "errmsg":"sign not match"
//}
//
// IP地址不在白名单
//{
//  "errcode":310000,
//  "errmsg":"ip X.X.X.X not in whitelist"
//}
type BaseResp struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	AtUserIds []string `json:"atUserIds,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}

type Text struct {
	Content string `json:"content,omitempty"`
}

type TextMsg struct {
	At      *At    `json:"at,omitempty"`
	Text    *Text  `json:"text,omitempty"`
	Msgtype string `json:"msgtype,omitempty"`
}

type Link struct {
	Text       string `json:"text,omitempty"`
	Title      string `json:"title,omitempty"`
	PicUrl     string `json:"picUrl,omitempty"`
	MessageUrl string `json:"messageUrl,omitempty"`
}

type LinkMsg struct {
	Link    *Link  `json:"link,omitempty"`
	Msgtype string `json:"msgtype,omitempty"`
}

type Markdown struct {
	Title string `json:"title,omitempty"`
	Text  string `json:"text,omitempty"`
}

type MarkdownMsg struct {
	At       *At       `json:"at,omitempty"`
	Markdown *Markdown `json:"markdown,omitempty"`
	Msgtype  string    `json:"msgtype,omitempty"`
}
