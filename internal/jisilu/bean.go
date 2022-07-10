package jisilu

import "net/http"

// 账号信息
const (
	loginProcess = "https://www.jisilu.cn/webapi/account/login_process/"
	returnUrl    = "https://www.jisilu.cn/"

	autoLogin = "1"
	aes       = "1"
)

var (
	header = http.Header{
		"Content-Type": []string{"application/x-www-form-urlencoded; charset=UTF-8"},
		"User-Agent":   []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.83 Safari/537.36"},
		"Host":         []string{"www.jisilu.cn"},
		"Origin":       []string{"https://www.jisilu.cn"},
		"Referer":      []string{"https://www.jisilu.cn/account/login/"},
	}
)

type Data struct {
	Page int   `json:"page"`
	Rows []Row `json:"rows"`
}

type Row struct {
	ID   string                 `json:"id"`
	Cell map[string]interface{} `json:"cell"`
}

// 转债详细信息
type CB struct {
	BondID       string  `json:"bond_id" table:"代码"` // 代码
	BondNm       string  `json:"bond_nm" table:"名称"` // 名称
	RatingCd     string  `json:"rating_cd" table:"评级"`
	CurrIssAmt   float64 `json:"curr_iss_amt" table:"规模（亿）"`
	ConvertValue float64 `json:"convert_value" table:"转股价值"`
	PremiumRt    float64 `json:"premium_rt" table:"转股溢价率"`
	ListDt       string  `json:"list_dt" table:"上市日期"`
	StockID      string  `json:"stock_id" table:"正股代码"` // 正股名称
	StockNm      string  `json:"stock_nm" table:"正股名称"`
	Sprice       float64 `json:"sprice" table:"正股价"`
	SwCd         string  `json:"sw_cd" table:"行业代码"`
	PriceTips    string  `json:"price_tips" table:"提示"`
}

// 待发转债
type PreCB struct {
	BondID    string `json:"bond_id" table:"代码"` // 代码
	BondNm    string `json:"bond_nm" table:"名称"` // 名称
	ListDt    string `json:"list_date" table:"上市日期"`
	ApplyDate string `json:"apply_date" table:"申购日期"`
}

type Detail struct {
	RedeemPrice float64   `json:"redeem_price" table:"到期赎回价"` // 到期赎回价
	CpnDesc     []float64 `json:"cpn_desc" table:"利率"`        // 利率
}

type CBAdvice struct {
	BondID       string  `json:"bond_id" table:"代码"`
	BondNm       string  `json:"bond_nm" table:"名称"`
	RatingCd     string  `json:"rating_cd" table:"评级"`
	CurrIssAmt   float64 `json:"curr_iss_amt" table:"规模（亿）"`
	ConvertValue float64 `json:"convert_value" table:"转股价值"`
	PremiumRt    float64 `json:"premium_rt" table:"转股溢价率"`
	ExpiryValue  float64 `json:"expiry_value" table:"到期价值"`
	DebtValue    float64 `json:"debt_value" table:"纯债现值"`
	ListForecast int     `json:"list_forecast" table:"上市预测"`
	ApplyAdvice  string  `json:"apply_advice" table:"申购建议"`
	ListDt       string  `json:"list_dt" table:"上市日期"`
	ApplyDate    string  `json:"apply_date" table:"申购日期"`
}
