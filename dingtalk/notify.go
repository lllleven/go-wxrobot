package dingtalk

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func NewDingTalkRobot(accessToken, secret string) *DingTalkRobot {
	return &DingTalkRobot{
		AccessToken: accessToken,
		Secret:      secret,
	}
}

func (d *DingTalkRobot) Send(msg interface{}) error {
	body := bytes.NewBuffer(nil)
	err := json.NewEncoder(body).Encode(msg)
	if err != nil {
		return err
	}
	client := &http.Client{}
	value := url.Values{}
	value.Set("access_token", d.AccessToken)
	if d.Secret != "" {
		t := time.Now().UnixMilli()
		value.Set("timestamp", fmt.Sprint(t))
		value.Set("sign", d.Sign(t, d.Secret))
	}

	req, err := http.NewRequest(http.MethodPost, "https://oapi.dingtalk.com/robot/send", body)
	if err != nil {
		return err
	}
	req.URL.RawQuery = value.Encode()
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	result, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", result)
	}
	var resp BaseResp
	if err := json.Unmarshal(result, &resp); err != nil {
		return err
	}
	if resp.ErrCode != 0 {
		return fmt.Errorf("%d:%s", resp.ErrCode, resp.ErrMsg)
	}
	return nil
}

func (d *DingTalkRobot) Sign(t int64, secret string) string {
	strToHash := fmt.Sprintf("%d\n%s", t, secret)
	hmac256 := hmac.New(sha256.New, []byte(secret))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}

func (d *DingTalkRobot) SendText(ctx context.Context, text string, at *At) error {
	if text == "" {
		return fmt.Errorf("text is empty")
	}
	return d.Send(TextMsg{
		Msgtype: MessageTypeText,
		Text:    &Text{Content: text},
		At:      at,
	})
}

func (d *DingTalkRobot) SendLink(ctx context.Context, title, text, messageURL, picURL string) error {
	if title == "" {
		return fmt.Errorf("title is empty")
	}
	if text == "" {
		return fmt.Errorf("text is empty")
	}
	if messageURL == "" {
		return fmt.Errorf("messageURL is empty")
	}
	return d.Send(LinkMsg{
		Msgtype: MessageTypeLink,
		Link:    &Link{Title: title, Text: text, MessageUrl: messageURL, PicUrl: picURL},
	})
}

func (d *DingTalkRobot) SendMarkdown(ctx context.Context, title, text string, at *At) error {
	if title == "" {
		return fmt.Errorf("title is empty")
	}
	if text == "" {
		return fmt.Errorf("text is empty")
	}
	return d.Send(MarkdownMsg{
		Msgtype:  MessageTypeMarkdown,
		Markdown: &Markdown{Title: title, Text: text},
		At:       at,
	})
}
