package utils

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/chromedp/chromedp"
	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/gomail.v2"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//LockIeCurrency
/**
 * 接口入口
 */
type LockIeCurrency struct {
}

//HeaderPcArr
/**
 *  电脑端采集头
 */
var HeaderPcArr = []string{
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3809.71 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_0_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
	"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36 Edg/100.0.1185.50",
	"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) HeadlessChrome/99.0.4844.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.69",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.5414.74 Safari/537.36",
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.70",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 Edg/109.0.1518.55",
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
}

//HeaderModelArr
/**
 *  手机端采集头
 */
var HeaderModelArr = []string{
	"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 10; HarmonyOS; LYA-AL00; HMSCore 6.9.0.302; GMSCore 20.15.16) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.88 HuaweiBrowser/13.0.1.301 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.25 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 13; 22041211AC Build/TP1A.220624.014; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/86.0.4240.99 XWEB/4375 MMWEBSDK/20221206 Mobile Safari/537.36 MMWEBID/1017 MicroMessenger/8.0.32.2300(0x28002053) WeChat/arm64 Weixin NetType/WIFI Language/zh_CN ABI/arm64",
	"Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Mobile Safari/537.36 MicroMessenger/7.0.1",
	"Mozilla/5.0 (Linux; Android 8.0; Pixel 2 Build/OPD3.170816.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Mobile Safari/537.36",
}

//HeaderCaiJiArr
/**
 * 内容采集访问头
 */
var HeaderCaiJiArr = []string{
	"Mozilla/5.0 (compatible; Baiduspider/2.0;++http://www.baidu.com/search/spider.html)",                           //百度蜘蛛
	"Mozilla/5.0 (compatible; Googlebot/2.1;+http://www.google.com/bot.html)",                                       //谷歌蜘蛛
	"Mozilla/5.0 (compatible; BingPreview/2.0;+http://www.bing.com/bingbot.htm)",                                    //必应蜘蛛
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36", //360蜘蛛
}

//QueryDetail
/**
 * 使用异步chromedb抓取加载完成后的HTML页面 urlpath 访问路径，proxyName 代理地址
 */
func (w *LockIeCurrency) QueryDetail(urlPath, proxyName string) (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	//判断是否有代理
	if proxyName != "" {
		//设置代理
		o := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.ProxyServer(proxyName))
		o = append(chromedp.DefaultExecAllocatorOptions[:], chromedp.UserAgent(HeaderPcArr[rand.Intn(len(HeaderPcArr))]))
		ctx, cancel = chromedp.NewExecAllocator(context.Background(), o...)
		ctx, cancel = chromedp.NewContext(ctx)
	}
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	var example string
	err := chromedp.Run(ctx,
		// 导航到https://www.baidu.com/
		chromedp.Navigate(urlPath),
		// 等待body > footer元素渲染完成
		chromedp.WaitVisible(`body`),
		// 点击指定的元素
		//chromedp.Click(`/html/body/div[6]/div/div[1]/form/div[1]/h3`, chromedp.NodeVisible),
		//获取指定元素内容
		chromedp.OuterHTML(`document.querySelector("body")`, &example, chromedp.ByJSPath),
	)
	if err != nil {
		return "", fmt.Errorf("运行模拟器异常：%s", err)
	}
	//使用插件分析分解html标签
	return example, nil

}

//Resolve
/**
 * 读取指定标签内容 htmlStr html源码  labelStr 需要提取内容的标签名称如：p
 */
func (w *LockIeCurrency) Resolve(htmlStr, labelStr string) (string, error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		return "", fmt.Errorf("HTML分解异常：%s", err)
	}
	//获取p标签内容
	strStr := ""
	dom.Find(labelStr).Each(func(i int, selection *goquery.Selection) {
		strStr = fmt.Sprintf("%s<%s>%s</%s>", strStr, labelStr, selection.Text(), labelStr)
	})
	return strStr, nil
}

//ZzFind
/**
 * 使用正则表达式拿到一条数据
 * @regular 正则表达式，语法如：<p>(.*?)</p>
 * @str 需要检索的字符串
 * @return 匹配成功的字符串
 */
func (l *LockIeCurrency) ZzFind(regular, str string) string {
	resArr := l.ZzeArr(regular, str)
	if resArr != nil && resArr[0] != nil && resArr[0][1] != "" {
		return resArr[0][1]
	}
	return ""
}

//ZzeArr
/**
 * 正则表达式检索数据，并返回检索结果
 * @regular 正则表达式语法如：<p>(.*?)</p>
 * @str 需要检索的字符串
 * @return 结果集
 */
func (l *LockIeCurrency) ZzeArr(regular string, str string) [][]string {
	rp := regexp.MustCompile(regular)
	findArr := rp.FindAllStringSubmatch(str, -1)
	if len(findArr) < 1 {
		return nil
	}
	return findArr
}

//GetLocalIP
/**
 * 获取访问IP地址
 */
func (l *LockIeCurrency) GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

//GetIpArea
/**
 * 获取IP地址-综合下面1-2-3
 */
func (l *LockIeCurrency) GetIpArea(ipStr string) string {
	if ipStr == "" {
		ipStr = l.GetLocalIP()
	}
	if ipStr == "127.0.0.1" || ipStr == "localhost" {
		return "本地"
	}
	addTxt := l.GetIpAddress(ipStr)
	if addTxt != "未知" && addTxt != "" {
		return addTxt
	}
	addTxt = l.GetIpAddress2(ipStr)
	if addTxt != "未知" && addTxt != "" {
		return addTxt
	}
	addTxt = l.GetIpAddress3(ipStr)
	if addTxt != "未知" && addTxt != "" {
		return addTxt
	}
	addTxt = l.GetIpAddress4(ipStr)
	if addTxt != "未知" && addTxt != "" {
		return addTxt
	}
	return "未知"
}

//GetIpAddress
/**
 * 查询IP方式1
 */
func (l *LockIeCurrency) GetIpAddress(ip string) string {
	res, err := HttpGet(fmt.Sprintf("https://www.ip.cn/ip/%s.html", ip))
	if err != nil {
		return "未知错误"
	}
	return l.ZzFind(`<div id="tab0_address">(.*?)</div>`, res)
}

//GetIpAddress2
/**
 *  查询IP地址方式2
 */
func (l *LockIeCurrency) GetIpAddress2(ip string) string {
	var params struct {
		Status     string `json:"status"`
		Country    string `json:"country"`
		RegionName string `json:"regionName"`
		City       string `json:"city"`
	}
	res, err := HttpGet(fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ip))
	if err != nil {
		return ""
	}
	err = json.Unmarshal([]byte(res), &params)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s%s%s", params.Country, params.RegionName, params.City)
}

//GetIpAddress3
/**
 *  查询IP地址方式3
 */
func (l *LockIeCurrency) GetIpAddress3(ip string) string {
	var params struct {
		Status    string `json:"status"`
		Info      string `json:"info"`
		Infocode  string `json:"infocode"`
		Province  string `json:"province"`
		City      string `json:"city"`
		Adcode    string `json:"adcode"`
		Rectangle string `json:"rectangle"`
	}
	res, err := HttpGet(fmt.Sprintf("http://restapi.amap.com/v3/ip?key=a0e21986aa11f68ce6d1d5da00b1d423&ip=%s", ip))
	if err != nil {
		return ""
	}
	err = json.Unmarshal([]byte(res), &params)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s%s", params.Province, params.City)
}

//GetIpAddress4
/**
 * 查询IP地址方式4
 */
func (l *LockIeCurrency) GetIpAddress4(ipStr string) string {
	if ipStr == "" {
		return "无"
	}
	urlPath := fmt.Sprintf("https://www.ip138.com/iplookup.asp?ip=%s&action=2", ipStr)
	res := l.HttpGetPc(urlPath, "www.ip138.com", "https://")
	//currency :=new(middles.Currency)
	resStr := l.ConvertGBK2Str(string(res))
	par := `{"begin":.*?, "end":.*?, "ct":"(.*?)", "prov":"(.*?)", "city":"(.*?)", "area":"(.*?)", "idc":"", "yunyin":"(.*?)", "net":""}`
	resArr := l.ZzeArr(par, resStr)
	city := ""
	if resArr != nil {
		for _, v := range resArr {
			for k, vo := range v {
				if k > 0 {
					if k == 5 {
						vo = fmt.Sprintf("-%s", vo)
					}
					city = fmt.Sprintf("%s%s", city, vo)
				}
			}
		}
	}
	return city
}

/********************************** http start **************************************/

//HttpGetMobile
/**
 * 模拟手机端启用get方式访问
 */
func (l *LockIeCurrency) HttpGetMobile(urlStr, domain, http string) []byte {
	//请求头
	agentStr := HeaderModelArr[l.RandRange(0, len(HeaderModelArr)-1)]
	var header = map[string]string{
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                agentStr,
		"Referer":                   http + domain,
	}
	return l.GetHtmlHeader(urlStr, &header)
}

//HttpGetPc
/**
 * 模拟电脑端使用浏览器访问
 */
func (l *LockIeCurrency) HttpGetPc(urlStr, domain, httpStr string) []byte {
	agentStr := HeaderPcArr[l.RandRange(0, len(HeaderPcArr)-1)]
	//随机生成
	var header = map[string]string{
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                agentStr,
		"Referer":                   httpStr + domain,
	}
	return l.GetHtmlHeader(urlStr, &header)
}

//HttpGather
/**
 * 采集数据时调用此方法
 */
func (l *LockIeCurrency) HttpGather(urlStr, domain, httpStr string) []byte {
	agentStr := HeaderCaiJiArr[l.RandRange(0, len(HeaderCaiJiArr)-1)]
	//随机生成
	var header = map[string]string{
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                agentStr,
		"Referer":                   httpStr + domain,
	}
	return l.GetHtmlHeader(urlStr, &header)
}

//GetHtmlHeader
/**
 * 传入URL地址进行远程访问并返回内容
 */
func (l *LockIeCurrency) GetHtmlHeader(urlPath string, header *map[string]string) []byte {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil || req == nil {
		fmt.Println("访问错误：", err)
		return nil
	}
	for key, value := range *header {
		req.Header.Add(key, value)
	}
	resp, err1 := client.Do(req)
	if err1 != nil {
		fmt.Println("访问错误2：", err1)
		return nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("访问错误3：", err2)
		return nil
	}
	return body
}

// ReStr
/**
 * 替换空和换行字符
 */
func (l *LockIeCurrency) ReStr(str string) string {
	htmlStr := strings.Replace(str, "\r\n", "", -1)
	htmlStr = strings.Replace(htmlStr, "\n", "", -1)
	htmlStr = strings.Replace(htmlStr, "\r", "", -1)
	htmlStr = strings.Replace(htmlStr, "&nbsp;", "", -1)
	htmlStr = strings.Replace(htmlStr, "\t", "", -1)
	htmlStr = strings.Replace(htmlStr, " ", "", -1)
	return htmlStr
}

//GetUrlImgToBase64
/**
 * 远程图片转base64
 */
func (l *LockIeCurrency) GetUrlImgToBase64(urlPath string) string {
	//获取远端图片
	res, err := http.Get(urlPath)
	if err != nil {
		return ""
	}
	defer func() {
		_ = res.Body.Close()
	}()
	// 读取获取的[]byte数据
	data, _ := ioutil.ReadAll(res.Body)

	imageBase64 := base64.StdEncoding.EncodeToString(data)
	return imageBase64
}

//RandomStr
/**
 * 传入字符数组，并随机打乱顺序
 */
func (l *LockIeCurrency) RandomStr(strings []string) string {
	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}
	str := ""
	for i := 0; i < len(strings); i++ {
		str += strings[i]
	}
	return str
}

//RandomRune
/**
 * 传入指定字符串生成随机字符串
 * letters := []rune("阿啊哀唉挨矮爱碍安岸按案暗昂袄傲奥八巴扒吧疤拔把坝爸罢霸白百柏摆败拜班
 * num 生成的长度
 */
func (l *LockIeCurrency) RandomRune(letters []rune, num int) string {
	b := make([]rune, num)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//RandRange
/**
 * 生成范围随机数
 */
func (l *LockIeCurrency) RandRange(min, max int) int {
	z := rand.Intn(max - min)
	return z + min
}

//GetLcTimeInt64
/**
 * 获取当日凌晨时间戳
 */
func (l *LockIeCurrency) GetLcTimeInt64() int64 {
	timeStr := time.Now().Format("2006-01-02")
	ti, _ := time.Parse("2006-01-02", timeStr)
	exitTime := int64(8 * 60 * 60)
	return ti.Unix() - exitTime
}

//SendMail
/**
 * 发送邮件
 * @mailTo 需要接收邮件的邮箱地址，数组格式
 * @conf 服务邮箱配置信息
 * @alias 发送邮箱别名
 * @subject 发送主题
 * @body 发送内容
 * @filePath 附件物理地址，没有时为空
 * @fileName 附件名称，没有时为空
 */
func (l LockIeCurrency) SendMail(mailTo, conf []string, alias, subject, body, filePath, fileName string) error {
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
	//mailConn := map[string]string{
	//  "user": "xxx@163.com",
	//  "pass": "your password",
	//  "host": "smtp.163.com",
	//  "port": "465",
	//}
	mailConn := map[string]string{
		"user": conf[0], //conf[0] 需要发送邮件出去的邮箱地址
		"pass": conf[1], //conf[1]邮箱密码
		"host": conf[2], //conf[2] 邮箱服务地址，如：smtp.qq.com
		"port": conf[3], //conf[3] 邮箱发送端口号，如：465
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(mailConn["user"], alias)) //这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo...)    //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	//判断是否携带附件
	if filePath != "" {
		m.Attach(filePath, gomail.Rename(fileName), gomail.SetHeader(map[string][]string{
			"Content-Dispostition": {
				fmt.Sprintf("attachment;filename=\"%s\"", mime.QEncoding.Encode("UTF-8", fileName)),
			},
		}))
	}
	m.SetBody("text/html", body) //设置邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	return err
}

//Zip
/**
 * 生成ZIP文件
 * srcDir 要生成的所在文件目录
 * zipFileName 生成后的文件名，含路径一起
 */
func (p *LockIeCurrency) Zip(srcDir string, zipFileName string) error {

	dir, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	if len(dir) == 0 {
		return fmt.Errorf("空")
	}
	// 预防：旧文件无法覆盖  删除路径下所有的文件
	err = os.RemoveAll(zipFileName)

	// 创建：zip文件
	zipfile, _ := os.Create(zipFileName)
	defer func() {
		_ = zipfile.Close()
	}()
	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer func() {
		_ = archive.Close()
	}()
	// 遍历路径信息
	_ = filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {
		// 如果是源路径，提前进行下一个遍历
		if path == srcDir {
			return nil
		}
		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, srcDir+`\`)
		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}
		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer func() {
				_ = file.Close()
			}()
			_, _ = io.Copy(writer, file)
		}

		return nil
	})
	return nil
}

//ConvertGBK2Str
/**
 * 字符串转换成UTF8
 * @gbkStr 需要转换的字符
 */
func (l LockIeCurrency) ConvertGBK2Str(gbkStr string) string {
	//将GBK编码的字符串转换为utf-8编码
	ret, err := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	if err != nil {
		return err.Error() //如果转换失败返回空字符串
	}
	return ret
}

//UngZip
/**
 * 解压GZIP文件
 */
func (l *LockIeCurrency) UngZip(data []byte) (string, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer func() {
		_ = reader.Close()
	}()
	data, err = ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(data), nil

}

//StrToStrToTime
/**
 * 字符串格式转时间戳
 */
func (l *LockIeCurrency) StrToStrToTime(str string) int64 {
	if str == "" {
		return 0
	}
	var time64 int64 = 0
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	if str != "" && strings.Contains(str, ":") {
		//判断：出现的次数
		if len(str) == 16 {
			str = fmt.Sprintf("%s:00", str)
		}
		tm2, _ := time.ParseInLocation("2006-01-02 15:04:05", str, Loc)
		time64 = tm2.Unix()
	} else {
		tm2, _ := time.ParseInLocation("2006-01-02", str, Loc)
		time64 = tm2.Unix()
	}
	return time64
}

//Int64ToTimeToStr
/**
 * 时间戳转日期
 * @timeInt 时间戳
 * @typeFormat 需要转换的格式
 */
func (l *LockIeCurrency) Int64ToTimeToStr(timeInt int64, typeFormat string) string {
	ti := time.Unix(timeInt, 0)
	//typeFormat如："2006-01-02 15:04:15"
	return ti.Format(typeFormat)
}

//Encrypt
/**
 * 加密  aes-128-cbc
 */
func (l *LockIeCurrency) Encrypt(key string, iv string, str string) (string, error) {
	if str == "" {
		return "", nil
	}
	encodeKey := []byte(key)
	encodeIv := []byte(iv)
	encodeStr := []byte(str)
	//获取block块
	block, err := aes.NewCipher(encodeKey)
	if err != nil {
		return "", err
	}
	//补码
	encodeStr = l.PKCS7Padding(encodeStr, block.BlockSize())
	//加密模式
	blockMode := cipher.NewCBCEncrypter(block, encodeIv)
	//创建明文长度的数组
	dCrypts := make([]byte, len(encodeStr))
	//加密明文
	blockMode.CryptBlocks(dCrypts, encodeStr)
	return base64.StdEncoding.EncodeToString(dCrypts), nil
}

//PKCS7Padding
/**
 * 补码
 */
func (l *LockIeCurrency) PKCS7Padding(origData []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	//在切片后面追加char数量的byte(char)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}

//Decrypt
/**
 * 解密
 */
func (l *LockIeCurrency) Decrypt(key string, iv string, str string) (string, error) {
	if str == "" {
		return "", nil
	}
	decodeKey := []byte(key)
	decodeIv := []byte(iv)
	strByte, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		//判断是否最后没有等号，加上等号再试
		if !strings.Contains(str, "==") {
			str = fmt.Sprintf("%s==", str)
			strByte, err = base64.StdEncoding.DecodeString(str)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}

	}
	//获取block块
	block, err := aes.NewCipher(decodeKey)
	if err != nil {
		return "", err
	}
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, decodeIv)
	//创建明文长度的数组
	decrypted := make([]byte, len(strByte))
	//加密明文
	blockMode.CryptBlocks(decrypted, strByte)
	//去掉字符
	decrypted = l.PKCS7upPadding(decrypted)
	return string(decrypted), nil
}

//PKCS7upPadding
/**
 * 去掉字符
 */
func (l *LockIeCurrency) PKCS7upPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//ReckonLatLon
/**
 * 计算两地经纬度相差距离
 * @lat1 维度1
 * @lat2 维度2
 * @lon1 经度1
 * @lon2 经度2
 */
func (l *LockIeCurrency) ReckonLatLon(lat1, lat2, lon1, lon2 float64) float64 {
	earthRadius := 6378.137
	pi := 3.1415926535898
	radLat1 := lat1 * pi / 180
	radLat2 := lat2 * pi / 180
	a := radLat1 - radLat2
	b := (lon1 * pi / 180) - (lon2 * pi / 180)
	s := 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * earthRadius
	return math.Round(s*10000) / 10000
}

//SendAliSms
/**
 * 发送阿里云短信
 * @conf 发送配置信息（三个参数：第一个参数是发送区域，第二个参数是APPKey，第三个参数是Secret）
 * @phone 发送手机号
 * @sign 短信签名
 * @tempId 模板ID
 * @data 发送内容
 */
func (l *LockIeCurrency) SendAliSms(conf []string, phone, sign, tempId, data string) (bool, string) {
	//conf 第一个参数是发送区域，第二个参数是APPKey，第三个参数是Secret
	client, err := dysmsapi.NewClientWithAccessKey(conf[0], conf[1], conf[2])
	/* use STS Token
	client, err := dysmsapi.NewClientWithStsToken("cn-qingdao", "<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")
	*/
	if sign == "" {
		sign = "Lockie"
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"

	request.PhoneNumbers = phone  //接收短信的手机号码
	request.SignName = sign       //短信签名名称
	request.TemplateCode = tempId //短信模板ID
	if data != "" {
		request.TemplateParam = data
	}
	response, err := client.SendSms(request)
	if err != nil {
		return false, fmt.Sprintf("发送短信失败：%s", err)
	}
	return true, fmt.Sprintf("发送短信返回：%s", response)
}
