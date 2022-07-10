package chinabond

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var ycDefId = map[string]string{
	"AAA":  "2c9081e50a2f9606010a309f4af50111",
	"AAA-": "8a8b2ca045e879bf014607ebef677f8e",
	"AA+":  "2c908188138b62cd01139a2ee6b51e25",
	"AA":   "2c90818812b319130112c279222836c3",
	"AA-":  "8a8b2ca045e879bf014607f9982c7fc0",
	"A+":   "2c9081e91b55cc84011be40946ca0925",
	"A":    "2c9081e91e6a3313011e6d438a58000d",
	"A-":   "8a8b2ca04142df6a014148ca880f3046",
	"BBB+": "2c9081e91ea160e5011eab1f116c1a59",
	"BBB":  "8a8b2ca0455847ac0145650780ad68fb",
	"BB":   "8a8b2ca0455847ac0145650ba23b68ff",
	"B":    "8a8b2ca0455847ac0145650c3d726901",
	"CCC":  "8a8b2ca0455847ac0145650d03d26903",
	"CC":   "8a8b2ca0447ffc96014491641747535e",
}

// 中债企业债收益率
func YcDetail(year int, ratingCd string, date string) (float64, error) {
	id, ok := ycDefId[ratingCd]
	if !ok {
		return 0, errors.New("id not exist")
	}
	url := fmt.Sprintf("https://yield.chinabond.com.cn/cbweb-mn/yc/searchYc?xyzSelect=txy&&workTimes=%s&&dxbj=0&&qxll=0,&&yqqxN=N&&yqqxK=K&&ycDefIds=%s,&&wrjxCBFlag=0&&locale=zh_CN",
		date, id)
	resp, err := http.Post(url, "application/json;charset=UTF-8", nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	data := make([]struct {
		SeriesData [][]float64 `json:"seriesData"`
	}, 0)

	if err := json.Unmarshal(body, &data); err != nil {
		return 0, err
	}
	if len(data) == 0 {
		return 0, nil
	}

	for _, s := range data[0].SeriesData {
		if len(s) > 1 && int(s[0]) == year {
			return s[1], nil
		}

	}
	return 0, nil
}
