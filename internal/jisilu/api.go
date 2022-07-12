package jisilu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/guidoxie/cb/global"
	"github.com/guidoxie/cb/internal/chinabond"
	"github.com/spf13/cast"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type jiSiLu struct {
	industry *Industry
	list     []*CB
	client   *http.Client
}

func NewClient(isLogin ...bool) (*jiSiLu, error) {
	var res = &jiSiLu{}
	res.client = http.DefaultClient
	res.client.Transport = &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	res.client.Jar, _ = cookiejar.New(nil)
	if len(isLogin) > 0 && !isLogin[0] {
		return res, nil
	}
	form := make(url.Values)
	form.Add("return_url", returnUrl)
	form.Add("user_name", global.JiSiLu.UserName)
	form.Add("password", global.JiSiLu.Password)
	form.Add("auto_login", autoLogin)
	form.Add("aes", aes)
	req, err := http.NewRequest("POST", loginProcess, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header = header
	resp, err := res.client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respData := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}
	if err := json.Unmarshal(body, &respData); err != nil {
		return nil, err
	}
	if len(respData.Msg) != 0 {
		return nil, errors.New(respData.Msg)
	}
	return res, err
}

// 所有债券信息
func (j *jiSiLu) cbList() ([]*CB, error) {
	if len(j.list) > 0 {
		return j.list, nil
	}
	var addr = fmt.Sprintf("https://www.jisilu.cn/data/cbnew/cb_list_new/?___jsl=LST___t=%d", time.Now().UnixNano())

	req, err := http.NewRequest("POST", addr, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &Data{}
	if err := json.Unmarshal(bytes, data); err != nil {
		return nil, err
	}

	cb := make([]*CB, 0)

	for _, r := range data.Rows {
		cell := &CB{}
		if err := MarshalAndUnmarshal(cell, &r.Cell); err != nil {
			return nil, err
		}
		cb = append(cb, cell)
	}
	j.list = cb
	return cb, nil
}

// 待发转债
func (j *jiSiLu) preList() ([]*PreCB, error) {
	var addr = fmt.Sprintf("https://www.jisilu.cn/data/cbnew/pre_list/?___jsl=LST___t=%d", time.Now().UnixNano())

	req, err := http.NewRequest("POST", addr, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	data := &Data{}
	if err := json.Unmarshal(bytes, data); err != nil {
		return nil, err
	}

	cb := make([]*PreCB, 0)

	for _, r := range data.Rows {
		cell := &PreCB{}
		if err := MarshalAndUnmarshal(cell, &r.Cell); err != nil {
			return nil, err
		}
		// 过滤尚未批准的
		if len(cell.BondID) > 0 {
			cb = append(cb, cell)
		}
	}
	return cb, nil
}

func (j *jiSiLu) RedeemPriceANDCpnDesc(bondID string) (redeemPrice float64, cpnDesc []float64, err error) {
	var addr = fmt.Sprintf("https://www.jisilu.cn/data/convert_bond_detail/%s", bondID)
	req, err := http.NewRequest("POST", addr, nil)
	if err != nil {
		return 0, nil, err
	}
	req.Header = header
	resp, err := j.client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}
	body := string(bytes)

	// 提取到期赎回价格
	redeemPrice = cast.ToFloat64(FindFirstSubMatch(`<td id="redeem_price">\s(.*?)\s</td>`, body))

	// 提取利率
	cpnDesc = make([]float64, 0)
	// <td id="cpn_desc" colspan="7">第一年0.30%、第二年0.50%、第三年1.00%、第四年1.50%、第五年1.80%、第六年2.00%</td>
	cd := regexp.MustCompile("[0-9]+(\\.[0-9]*)?").FindAllString(FindFirstSubMatch(`<td id="cpn_desc" colspan="7">(.*)</td>`, body), -1)
	for _, c := range cd {
		cpnDesc = append(cpnDesc, cast.ToFloat64(c))
	}
	// 提取行业
	// <a href="/data/cbnew/industry-630701#cb" target="_cblist">电力设备-电池-锂电池</a>
	//res.Industry = FindFirstSubMatch(`target="_cblist">(.*?)</a>`, body)

	return
}

// 待发转债
func (j *jiSiLu) PreCBList(date ...string) ([]CBAdvice, error) {
	var (
		res        = make([]CBAdvice, 0)
		filterDate = time.Now().Format("2006-01-02")
	)

	if len(date) > 0 {
		filterDate = date[0]
	}

	// 获取待发转债
	preList, err := j.preList()
	if err != nil {
		return nil, err
	}

	for _, p := range preList {
		if p.ApplyDate == filterDate || p.ListDt == filterDate {
			cb, err := j.findOne(map[string]interface{}{
				"bond_id": p.BondID,
			})
			if err != nil {
				return nil, err
			}
			sub := CBAdvice{
				Market:       GetMarket(cb.BondID),
				BondID:       cb.BondID,
				BondNm:       cb.BondNm,
				RatingCd:     cb.RatingCd,
				CurrIssAmt:   cb.CurrIssAmt,
				ConvertValue: cb.ConvertValue,
				PremiumRt:    cb.PremiumRt,
				ListDt:       ParseDate(p.ListDt),
				ApplyDate:    ParseDate(p.ApplyDate),
			}
			// ExpiryValue 到期价值=票面利率+赎回价
			var year int
			year, sub.ExpiryValue, err = j.ExpiryValue(cb.BondID)
			if err != nil {
				return nil, err
			}

			// DebtValue 纯债价值
			sub.DebtValue, err = j.DebtValue(sub.ExpiryValue, year, sub.RatingCd, filterDate)
			if err != nil {
				return nil, err
			}
			// ListForecast 上市预测
			sub.ListForecast, err = j.ListForecast(cb)

			// ApplyAdvice 申购建议
			j.ApplyAdvice(&sub)
			res = append(res, sub)
		}
	}
	return res, nil
}

func (j *jiSiLu) findOne(condition map[string]interface{}) (*CB, error) {
	res, err := j.find(condition)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return nil, nil
}

func (j *jiSiLu) find(condition map[string]interface{}) (res []*CB, err error) {
	if len(condition) == 0 {
		return nil, nil
	}

	list, err := j.cbList()
	if err != nil {
		return nil, err
	}
	for k, v := range condition {
		for i := range list {
			switch k {
			case "bond_id":
				if list[i].BondID == cast.ToString(v) {
					res = append(res, list[i])
				}
			case "sw_cd":
				if list[i].SwCd == cast.ToString(v) {
					res = append(res, list[i])
				}
			}
		}
	}
	return
}

// 上市预测
func (j *jiSiLu) ListForecast(cb *CB) (int, error) {
	// ListForecast 上市预测 = 转股价值 * (1+参考溢价率)

	// 找出同行业已上市的
	list := make([]*CB, 0)
	industry, err := j.industryList()
	if err != nil {
		return 0, err
	}
	for _, swCb := range industry.SubLevel(industry.TopLevel(cb.SwCd).Number) {

		l, err := j.find(map[string]interface{}{"sw_cd": swCb.Number})
		if err != nil {
			return 0, err
		}
		for i := range l {
			if l[i].BondID != cb.BondID && len(l[i].ListDt) != 0 {
				list = append(list, l[i])
			}
		}
	}

	// 选出转股价值比较相近的
	var c *CB
	var abs = math.MaxFloat64
	for _, l := range list {
		if a := math.Abs(l.ConvertValue - cb.ConvertValue); a < abs {
			c = l
			abs = a
		}
	}
	if c != nil && abs != math.MaxFloat64 {
		var premiumRt = c.PremiumRt
		if c.ConvertValue < cb.ConvertValue {
			premiumRt += (c.PremiumRt/c.ConvertValue)*(cb.ConvertValue-c.ConvertValue) + 1
		} else if c.ConvertValue > cb.ConvertValue {
			premiumRt -= (c.PremiumRt/c.ConvertValue)*(c.ConvertValue-cb.ConvertValue) - 1
		}
		// TODO 暂不考虑评级问题
		return int(cb.ConvertValue * (1 + premiumRt/100)), nil
	}
	return 0, err
}

// ExpiryValue 到期价值=票面利率+赎回价
func (j *jiSiLu) ExpiryValue(bondID string) (int, float64, error) {
	r, cpnDesc, err := j.RedeemPriceANDCpnDesc(bondID)
	if err != nil {
		return 0, 0, err
	}

	var expiryValue = r
	for _, c := range cpnDesc[:len(cpnDesc)-1] {
		expiryValue += c
	}
	return len(cpnDesc), expiryValue, nil
}

// DebtValue 纯债价值 = 到期赎回价/(1+中债企业债收益率)^到期年数
func (j *jiSiLu) DebtValue(expiryValue float64, year int, ratingCd, date string) (float64, error) {
	yc, err := chinabond.YcDetail(year, ratingCd, date)
	if err != nil {
		return 0, err
	}
	return Round(expiryValue/math.Pow(1+yc/100, float64(year)), 2), nil
}

/*
第一，评级不低于A+或者AA-，看个人对于“破发”（价格跌破发行价）风险的承受力；
第二，看转股价值是否高于100；
第三，看转股溢价率是否为负；
第四，看纯债现值是否高于90；
*/
func (j *jiSiLu) ApplyAdvice(cb *CBAdvice) {
	advice := make([]string, 0)
	if cb.PremiumRt <= 0 {
		advice = append(advice, "★")
	}
	if cb.ConvertValue > 100 {
		advice = append(advice, "★")
	}
	if strings.Contains(cb.RatingCd, "A") {
		advice = append(advice, "★")
	}
	if cb.DebtValue > 90 {
		advice = append(advice, "★")
	}

	for i := 4 - len(advice); i > 0; i-- {
		advice = append(advice, "☆")
	}
	cb.ApplyAdvice = strings.Join(advice, "")
}

func (j *jiSiLu) industryList() (*Industry, error) {
	if j.industry != nil {
		return j.industry, nil
	}
	var addr = "https://www.jisilu.cn/data/cbnew/"
	req, err := http.NewRequest("GET", addr, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header
	resp, err := j.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	body := string(bytes)

	list := regexp.MustCompile(`<option data-pinyin=".*"  data-level="(.*)" value="(.*)" >(.*)</option>`).FindAllStringSubmatch(body, -1)

	res := &Industry{}
	for _, l := range list {
		res.Add(cast.ToInt(l[1]), l[2], l[3])
	}
	j.industry = res
	return res, nil
}
