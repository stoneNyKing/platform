package apis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"platform/common/utils"
	"strconv"
	"time"

	"github.com/boj/redistore"
	"gopkg.in/chanxuehong/wechat.v2/mp/core"
	"gopkg.in/chanxuehong/wechat.v2/mp/jssdk"
	"gopkg.in/chanxuehong/wechat.v2/mp/media"
	"gopkg.in/chanxuehong/wechat.v2/mp/menu"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/request"
	"gopkg.in/chanxuehong/wechat.v2/mp/message/callback/response"

	"github.com/chanxuehong/rand"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/sid"
	mpoauth2 "gopkg.in/chanxuehong/wechat.v2/mp/oauth2"
	"gopkg.in/chanxuehong/wechat.v2/oauth2"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"

	"platform/models/storages"
	. "platform/ousvc/common"
	"platform/ousvc/config"
	. "platform/ousvc/models"

	"crypto/sha1"
	"encoding/json"
	"encoding/xml"
	"github.com/levigross/grequests"
	"gopkg.in/chanxuehong/wechat.v2/mp/qrcode"
	"platform/ousvc/dbmodels"
	"sort"
	"strings"
)

var accessTokenServer *core.DefaultAccessTokenServer
var userClient *core.Client

var appId string
var oriId string
var token string

var (
	sessionStorage = session.New(20*60, 60*60)
	oauth2Endpoint oauth2.Endpoint
)

var (
	// 下面两个变量不一定非要作为全局变量, 根据自己的场景来选择.
	msgHandler core.Handler
	msgServer  *core.Server
)

func InitWeixin() {

	appId = config.Config.WxAppId

	oriId = config.Config.WxOriId
	token = config.Config.WxToken
	appSecret := config.Config.WxAppSecret

	accessTokenServer = core.NewDefaultAccessTokenServer(appId, appSecret, nil)
	userClient = core.NewClient(accessTokenServer, nil)
	oauth2Endpoint = mpoauth2.NewEndpoint(appId, appSecret)

	//创建message server
	messageServeMux := core.NewServeMux()
	messageServeMux.MsgHandleFunc(request.MsgTypeText, TextMessageHandler) // 注册文本处理 Handler
	messageServeMux.DefaultMsgHandleFunc(DefaultMessageHandler)
	messageServeMux.EventHandleFunc(menu.EventTypeClick, ClickHandler)
	messageServeMux.EventHandleFunc(request.EventTypeSubscribe, SubscribeHandler)
	messageServeMux.EventHandleFunc(request.EventTypeScan, ScanHandler)
	messageServeMux.DefaultEventHandleFunc(DefaultEventHandler)

	var msgHandler core.Handler
	msgHandler = messageServeMux

	// 下面函数的几个参数设置成你自己的参数: wechatId, token, appId
	msgServer = core.NewServer(config.Config.WxOriId, config.Config.WxAppId, config.Config.WxToken, config.Config.WxAesKeyEncode, msgHandler, nil)

}

// wxCallbackHandler 是处理回调请求的 http handler.
//  1. 不同的 web 框架有不同的实现
//  2. 一般一个 handler 处理一个公众号的回调请求(当然也可以处理多个, 这里我只处理一个)
func wxCallbackHandler(w http.ResponseWriter, r *http.Request) {
	msgServer.ServeHTTP(w, r, nil)
}

// 非法请求的 Handler
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Error("非法请求：%v", err.Error())
}

//建立必要的 session, 然后跳转到授权页面
func wechatServerAuth(w http.ResponseWriter, req *http.Request, session sessions.Session) {
	sid := sid.New()
	state := string(rand.NewHex())

	if err := sessionStorage.Add(sid, state); err != nil {
		io.WriteString(w, err.Error())
		log.Error("session 获取错误： %v", err)
		return
	}

	cookie := http.Cookie{
		Name:     "sid",
		Value:    sid,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	AuthCodeURL := mpoauth2.AuthCodeURL(appId, config.Config.WxRedirectUrl, config.Config.WxScope, state)
	log.Fine("AuthCodeURL: %v", AuthCodeURL)

	http.Redirect(w, req, AuthCodeURL, http.StatusFound)
}

// 授权后回调页面
func wechatAuthCallback(w http.ResponseWriter, r *http.Request, sessions sessions.Session) {
	log.Fine("request URL: %v", r.RequestURI)

	cookie, err := r.Cookie("sid")
	if err != nil {
		io.WriteString(w, err.Error())
		log.Error("获取sid cookie出错：%v", err)
		return
	}

	session, err := sessionStorage.Get(cookie.Value)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Error("获取seesion出错：%v", err)
		return
	}

	savedState := session.(string) // 一般是要序列化的, 这里保存在内存所以可以这么做

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Error("解析query失败: %v", err)
		return
	}

	code := queryValues.Get("code")
	if code == "" {
		log.Error("用户禁止授权")
		return
	}

	queryState := queryValues.Get("state")
	if queryState == "" {
		log.Error("state 参数为空")
		return
	}
	if savedState != queryState {
		str := fmt.Sprintf("state 不匹配, session 中的为 %q, url 传递过来的是 %q", savedState, queryState)
		io.WriteString(w, str)
		log.Error(str)
		return
	}

	oauth2Client := oauth2.Client{
		Endpoint: oauth2Endpoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Error("令牌交换失败：%v", err)
		return
	}

	log.Finest("token: %+v", token)

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		io.WriteString(w, err.Error())
		log.Error("获取用户信息失败：%v", err)
		return
	}

	//json.NewEncoder(w).Encode(userinfo)

	log.Finest("userinfo: %+v", userinfo)

	appkey := AppIds[int64(config.Config.WxMyAppid)]["key"]
	url1 := fmt.Sprintf(config.Config.WxServiceUrl+"/token/generate?appid=%d", config.Config.WxMyAppid)
	ti := time.Now().Unix()
	tk := utils.GetToken(url1, config.Config.WxMyAppid, appkey, ti)

	if tk == "" {
		log.Error("获取token失败.")
		return
	}

	url1 = fmt.Sprintf(config.Config.WxServiceUrl+"/user/check?appid=%d&token=%s", config.Config.WxMyAppid, tk)
	//uid,sid,err :=utils.CheckUser(url1,config.Config.WxMyAppid,tk,userinfo.OpenId,"weixinid")
	uid, sid := loginByOpenid(tk, userinfo.OpenId, int64(config.Config.WxMyAppid))

	//log.Finest("url1 = %s",url1)

	if err != nil {
		log.Error("检查微信id错误: %v", err)
	}
	log.Finest("用户微信登录 userid=%d,siteid=%d", uid, sid)
	if uid > 0 {
		params := fmt.Sprintf("?token=%s", tk) + "&userid=" + utils.ConvertToString(uid) + "&site=" + utils.ConvertToString(sid) + "&from=2" + "&openid=" + userinfo.OpenId
		u := "http://" + config.Config.WxMainHost + config.Config.WxMainPageNoAuth + params
		log.Finest("url(uid>0) = %s", u)
		http.Redirect(w, r, u, http.StatusFound)

	} else { //不存在openid
		params := fmt.Sprintf("?token=%s", tk) + "&openid=" + userinfo.OpenId + "&from=2"
		u := "http://" + config.Config.WxMainHost + config.Config.WxMainPage + params
		log.Finest("url(uid=0) = %s", u)
		http.Redirect(w, r, u, http.StatusFound)

	}

	//wxCallbackHandler(w,r)
	//return
}

// 文本消息的 Handler
func TextMessageHandler(ctx *core.Context) {
	// 简单起见，把用户发送过来的文本原样回复过去
	text := request.GetText(ctx.MixedMsg) // 可以省略...
	log.Finest("微信text：%v", text)
	if text.Content == "test" {
		resp := response.NewText(text.FromUserName, text.ToUserName, text.CreateTime, text.Content)

		ctx.ResponseWriter.Write([]byte(resp.Content))
	}
}

func DefaultMessageHandler(ctx *core.Context) {
	msg := ctx.MixedMsg
	log.Finest("DefaultMessageHandler: event_type=%v, event_key=%v", msg.EventType, msg.EventKey)
}

func ClickHandler(ctx *core.Context) {
	msg := ctx.MixedMsg
	log.Finest("ClickHandler: event_type=%v, event_key=%v", msg.EventType, msg.EventKey)
	articles := make([]response.Article, 1)
	articles[0].Title = "click_test_title"
	articles[0].Description = "click_test_Description"
	articles[0].URL = "http://dev.laoyou99.cn/weixin/" + msg.EventKey
	articles[0].PicURL = "http://dev.laoyou99.cn/20141213112008.jpg"
	resp := response.NewNews(msg.FromUserName, msg.ToUserName, time.Now().Unix(), articles)
	log.Finest("ClickHandler message : %v", msg)
	//ctx.ResponseWriter.Write([]byte(resp.Articles[0].Description))
	xml.NewEncoder(ctx.ResponseWriter).Encode(resp)
}

func SubscribeHandler(ctx *core.Context) {
	subscribe := request.GetSubscribeEvent(ctx.MixedMsg)
	openid := subscribe.FromUserName
	log.Finest("SubscribeHandler message : event_type=%v,event_key=%v,openid=%v", subscribe.EventType, subscribe.EventKey, openid)

	user := &dbmodels.User{Weixinid: openid}
	if u, err := dbmodels.GetUser(user); err != nil {
		log.Error("get user - error : %v", err)
	} else if u != nil {
		log.Finest("get user - have user: %v", openid)
	} else {
		log.Finest("get user - not have user: %v", openid)
		user = &dbmodels.User{Weixinid: openid}
		if _, err := dbmodels.RegisterUserByWeixinid(user); err != nil {
			log.Error("insert error: openid=%v,err=%v", openid, err)
		} else {
			log.Finest("insert ok: openid=%v", openid)
		}
	}

	if user.Id == 0 && subscribe.EventKey != "" {
		if scene, err := subscribe.Scene(); err == nil {
			if regid, err := strconv.Atoi(scene); err == nil && regid > 0 {
				user.Id = int64(regid)
				if _, err := dbmodels.BindUserWeixinidById(user); err != nil {
					log.Error("update user id error: %v,%v", regid, err)
				} else {
					log.Finest("update regid ok: openid=%v,regid=%v", openid, regid)
				}
			}
		}
	}

	if user.ImageUrl == "" {
		code := ctx.QueryParams.Get("code")
		oauth2Client := oauth2.Client{
			Endpoint: oauth2Endpoint,
		}
		token, err := oauth2Client.ExchangeToken(code)
		if err != nil {
			log.Error("令牌交换失败：%v", err)
			return
		}
		if userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil); err == nil {
			fmt.Println("get user info:", userinfo.HeadImageURL)
			resp, _ := grequests.Get(userinfo.HeadImageURL, nil)
			if resp.Error == nil {
				k := Md5_blob(resp.Bytes())
				storage := &storages.Storage{Key: k, Value: resp.Bytes()}
				if _, err := storages.SetStorage(storage); err == nil {
					user.ImageUrl = k
					if err := user.Update(); err != nil {
						log.Error("update error: openid=%v,err=%v", openid, err)
					} else {
						log.Finest("update ok: openid=%v,key=%v", openid, k)
					}
				} else {
					log.Error("storage install error : err=%v,key=%v", err, k)
				}
			} else {
				log.Error("Unable to make request : %v", resp.Error)
			}
		} else {
			log.Error("get userinfo error: openid=%v,err=%v", openid, err)
		}
	}

	articles := make([]response.Article, 1)
	articles[0].Title = "欢迎进入个人健康中心"
	articles[0].Description = "个人健康中心为您提供全方位的健康管理及体验。\n提示：点击「我」- 提交「我的信息」，开启健康管理。"
	articles[0].PicURL = "http://mobile.childfond.com/images/welcome.jpg"
	resp := response.NewNews(openid, subscribe.ToUserName, time.Now().Unix(), articles)
	log.Finest(" subscribe= %v", subscribe)

	xml.NewEncoder(ctx.ResponseWriter).Encode(resp)
}

func ScanHandler(ctx *core.Context) {
	scan := request.GetScanEvent(ctx.MixedMsg)
	openid := scan.FromUserName
	log.Finest("SubscribeHandler message : event_type=%v,event_key=%v,openid=%v", scan.EventType, scan.EventKey, openid)

	regid, err := strconv.Atoi(scan.EventKey)
	if err != nil && regid < 0 {
		log.Error("scan EventKey error: err= %v", err)
		return
	}

	user := &dbmodels.User{Weixinid: openid}
	if u, err := dbmodels.GetUser(user); err != nil {
		log.Error("get user - error : opendi=%v,err=%v", openid, err)
	} else if u != nil {
		log.Finest("get user - have user: openid=%v", openid)
		if user.Id == 0 {
			user.Id = int64(regid)
			if err := user.Update(); err != nil {
				log.Error("update error: openid=%v,err=%v", openid, err)
			} else {
				log.Finest("update ok: openid=%v,id=%v", openid, regid)
			}
		}
	} else {
		log.Error("get user - not have user: openid=%v", openid)
	}
}

func DefaultEventHandler(ctx *core.Context) {
	msg := ctx.MixedMsg
	log.Finest("DefaultEventHandler: event_type=%v, event_key=%v", msg.EventType, msg.EventKey)
}

//11
func WeixinHander() *martini.ClassicMartini {

	//jssdk
	var TicketServer = jssdk.NewDefaultTicketServer(core.NewClient(accessTokenServer, nil))

	m := martini.Classic()
	m.Use(render.Renderer())

	key := config.Config.SessionKey
	host := config.Config.SessionStoreIP
	port := config.Config.SessionStorePort

	store, _ := redistore.NewRediStore(10, "tcp", host+":"+port, "", []byte(key))
	m.Use(sessions.Sessions("wx_session", store))
	//weixin

	log.Finest("run here.")

	m.Get("/ticket", func(r render.Render) {
		if ticket, err := TicketServer.Ticket(); err != nil {
			r.JSON(200, map[string]interface{}{"Ret": -1, "Msg": err.Error()})
		} else {
			r.JSON(200, map[string]interface{}{"Ret": 0, "Ticket": ticket})
		}

	})

	m.Post("/sign", func(req *http.Request, r render.Render) {
		url := req.FormValue("url")
		if ticket, err := TicketServer.Ticket(); err != nil {
			r.JSON(200, map[string]interface{}{"Ret": -1, "Msg": err.Error()})
		} else {
			nonceStr := New(10)
			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			fmt.Println("weixin sign:", ticket, nonceStr, timestamp, url)
			sign := jssdk.WXConfigSign(ticket, nonceStr, timestamp, url)
			r.JSON(200, map[string]interface{}{"Ret": 0, "AppId": appId, "Timestamp": timestamp, "NonceStr": nonceStr, "Signature": sign})
		}
	})

	m.Get("/config", func(req *http.Request, r render.Render) {
		if !validateUrl(r, req) {
			log.Error("Wechat Service: this http request is not from Wechat platform!")
			return
		}
		//wxCallbackHandler(w,req)

	})

	m.Get("/qr/:id", func(r render.Render, params martini.Params) {
		id := params["id"]
		accountClient := core.NewClient(accessTokenServer, nil)
		if qr, err := qrcode.CreateStrScenePermQrcode(accountClient, id); err != nil {
			r.JSON(200, map[string]interface{}{"Ret": -1, "Msg": err.Error()})
		} else {
			r.JSON(200, map[string]interface{}{"Ret": 0, "Url": "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + qr.Ticket})
		}
	})

	m.Post("/cachemedia", func(req *http.Request, r render.Render) {
		k := req.FormValue("mediaid")
		var b bytes.Buffer
		buf := bufio.NewWriter(&b)
		mediaClient := core.NewClient(accessTokenServer, nil)
		media.DownloadToWriter(mediaClient, k, buf)
		storage := &storages.Storage{Key: k, Value: b.Bytes()}
		if _, err := storages.SetStorage(storage); err != nil {
			r.JSON(200, &storages.ErrorResp{Ret: 504, Msg: err.Error()})
			return
		}
		r.JSON(200, &storages.StorageKeyResp{Ret: 0, Key: k})
	})

	m.Get("/access_token", wechatServerAuth)
	m.Get("/phc", wechatAuthCallback)

	m.Any("/:path", func(w http.ResponseWriter, r *http.Request, params martini.Params) {
		path := params["path"]
		log.Finest("path=%v,url=%v", path, mpoauth2.AuthCodeURL(appId, config.Config.WxRedirectUrl, config.Config.WxScope, ""))
		http.Redirect(w, r, mpoauth2.AuthCodeURL(appId, config.Config.WxRedirectUrl, config.Config.WxScope, ""), http.StatusFound)
	})

	//m.Any(".*", wechatServer.ServeHTTP)

	return m
}

func makeSignature(timestamp, nonce string) string {
	sl := []string{token, timestamp, nonce}
	sort.Strings(sl)
	s := sha1.New()
	io.WriteString(s, strings.Join(sl, ""))
	return fmt.Sprintf("%x", s.Sum(nil))
}

func validateUrl(r render.Render, req *http.Request) bool {
	timestamp := strings.Join([]string{req.URL.Query().Get("timestamp")}, "")
	nonce := strings.Join([]string{req.URL.Query().Get("nonce")}, "")
	signatureGen := makeSignature(timestamp, nonce)

	signatureIn := strings.Join([]string{req.URL.Query().Get("signature")}, "")
	if signatureGen != signatureIn {
		return false
	}
	echostr := strings.Join([]string{req.URL.Query().Get("echostr")}, "")

	//fmt.Fprintf(w, echostr)

	r.Text(200, echostr)

	return true
}

func loginByOpenid(token, openid string, appid int64) (userid, siteid int64) {

	url := config.Config.WxServiceUrl + "/user/loginbyweixin?token=" + token + "&appid=" + utils.ConvertToString(appid)

	params := make(map[string]interface{})
	params["weixinid"] = openid
	b, err := json.Marshal(&params)
	if err != nil {
		log.Error("json masharl错误")
		return 0, 0
	}

	log.Finest("(loginByOpenid): url=%s", url)

	vs, err := utils.ServicePost(url, string(b))
	if err != nil {
		log.Error("获取用户id出错: %v", err)
		return 0, 0
	}

	if vs != nil {
		v := vs.(map[string]interface{})
		userid = utils.Convert2Int64(v["UserId"])
		siteid = utils.Convert2Int64(v["SiteId"])
	}

	return
}
