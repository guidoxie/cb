package cmd

import (
	"bytes"
	"fmt"
	"github.com/guidoxie/cb/internal/jisilu"
	"github.com/modood/table"
	"github.com/spf13/cobra"
	"html/template"
	"log"
	"time"
)

const (
	OutputTypeStdout = "stdout"
	OutputTypeHtml   = "html"
)

// 根命令
var (
	date   string
	output string
	root   = &cobra.Command{
		Use: "cb",
		Long: `可转债上市预测/申购建议

计算公式：
转股价值 = 正股价格/转股价*100

到期价值 = 票面利率+赎回价

纯债现值 =  到期赎回价/(1+中债企业债收益率)^到期年数

上市预测 = 转股价值*(1+参考同行业转债溢价率)

转股溢价率 = 转债价格/转股价值-100%

申购建议参考指标：
一，评级不低于A-

二，转股价值是否高于100

三，转股溢价率是否为负

四，纯债现值是否高于90
`,
		// 隐藏completion命令
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
		Run: func(cmd *cobra.Command, args []string) {
			client, err := jisilu.NewClient()
			if err != nil {
				log.Fatal(err)
			}
			cb, err := client.PreCBList(date)
			if err != nil {
				log.Fatal(err)
			}
			switch output {
			case OutputTypeStdout:
				table.Output(cb)
			case OutputTypeHtml:
				content, err := outputHtml(cb)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(content)
			}

		},
	}
)

func Execute() error {
	return root.Execute()
}

func init() {
	root.Flags().StringVarP(&date, "date", "d", time.Now().Format("2006-01-02"), "上市/申购日期")
	root.Flags().StringVarP(&output, "output", "o", "stdout", "输出类型，stdout：终端输出 html：输出HTML文件")
}

func outputHtml(cb []jisilu.CBAdvice) (string, error) {
	//读取模版文件
	var (
		buf   = bytes.NewBuffer(nil)
		list  = make([]jisilu.CBAdvice, 0)
		apply = make([]jisilu.CBAdvice, 0)
	)
	for i := range cb {
		if len(cb[i].ListDt) != 0 {
			list = append(list, cb[i])
		} else {
			apply = append(apply, cb[i])
		}
	}

	temp := template.Must(template.ParseFiles("config/template.html"))
	if err := temp.Execute(buf, map[string][]jisilu.CBAdvice{
		"list":  list,
		"apply": apply,
	}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
