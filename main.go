package main

import (
	"fyne.io/fyne/v2"
	"strconv"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

// convertToDecimal 将不同进制的字符串转换为十进制数
func convertToDecimal(inputStr string, base int) (float64, error) {
	if base == 10 {
		return strconv.ParseFloat(inputStr, 64)
	}
	// 先转换为 int64 类型
	num, err := strconv.ParseInt(inputStr, base, 64)
	if err != nil {
		return 0, err
	}
	return float64(num), nil
}

// convertFromDecimal 将十进制数转换为指定进制的字符串
func convertFromDecimal(num float64, base int) string {
	switch base {
	case 2:
		return strconv.FormatInt(int64(num), 2)
	case 8:
		return strconv.FormatInt(int64(num), 8)
	case 10:
		return strconv.FormatFloat(num, 'f', -1, 64)
	case 16:
		return strconv.FormatInt(int64(num), 16)
	default:
		return ""
	}
}

func main() {
	// 创建 Fyne 应用和窗口
	a := app.New()
	w := a.NewWindow("Base Converter v0.1")

	// 进制转换区域
	// 创建输入框和绑定数据
	convInput := widget.NewEntry()
	convInput.SetPlaceHolder("输入数值用于进制转换")
	convInputData := binding.NewString()
	convInput.Bind(convInputData)

	// 创建源进制和目标进制选择下拉框
	sourceBase := widget.NewSelect([]string{"2", "8", "10", "16"}, nil)
	sourceBase.SetSelected("10")
	targetBase := widget.NewSelect([]string{"2", "8", "10", "16"}, nil)
	targetBase.SetSelected("10")

	// 创建转换按钮和结果显示区域
	convertButton := widget.NewButton("转换", nil)
	convResultData := binding.NewString()
	convResult := widget.NewLabelWithData(convResultData)

	// 转换按钮点击事件处理函数
	convertButton.OnTapped = func() {
		inputStr, _ := convInputData.Get()
		sBase, _ := strconv.Atoi(sourceBase.Selected)
		tBase, _ := strconv.Atoi(targetBase.Selected)

		// 将输入转换为十进制
		decimal, err := strconv.ParseInt(inputStr, sBase, 64)
		if err != nil {
			convResultData.Set("输入无效")
			return
		}

		// 将十进制转换为目标进制
		var output string
		switch tBase {
		case 2:
			output = strconv.FormatInt(decimal, 2)
		case 8:
			output = strconv.FormatInt(decimal, 8)
		case 10:
			output = strconv.FormatInt(decimal, 10)
		case 16:
			output = strconv.FormatInt(decimal, 16)
		}
		convResultData.Set(output)
	}

	conversionContent := container.NewVBox(
		widget.NewLabel("进制转换区域"),
		convInput,
		container.NewHBox(
			widget.NewLabel("源进制:"),
			sourceBase,
			widget.NewLabel("目标进制:"),
			targetBase,
		),
		convertButton,
		widget.NewLabel("转换结果:"),
		convResult,
	)

	// 运算区域
	// 创建第二个输入框和运算结果显示区域
	opInput1 := widget.NewEntry()
	opInput1.SetPlaceHolder("输入第一个数值用于运算")
	opInputData1 := binding.NewString()
	opInput1.Bind(opInputData1)

	opInput2 := widget.NewEntry()
	opInput2.SetPlaceHolder("输入第二个数值用于运算")
	opInputData2 := binding.NewString()
	opInput2.Bind(opInputData2)

	operationResultData := binding.NewString()
	operationResult := widget.NewLabelWithData(operationResultData)

	// 新增运算进制选择下拉框
	operationBase := widget.NewSelect([]string{"2", "8", "10", "16"}, nil)
	operationBase.SetSelected("10")

	// 四则运算按钮点击事件处理函数
	performOperation := func(operator string) {
		inputStr1, _ := opInputData1.Get()
		inputStr2, _ := opInputData2.Get()
		base, _ := strconv.Atoi(operationBase.Selected)

		num1, err1 := convertToDecimal(inputStr1, base)
		num2, err2 := convertToDecimal(inputStr2, base)

		if err1 != nil || err2 != nil {
			operationResultData.Set("输入无效")
			return
		}

		var result float64
		switch operator {
		case "+":
			result = num1 + num2
		case "-":
			result = num1 - num2
		case "*":
			result = num1 * num2
		case "/":
			if num2 == 0 {
				operationResultData.Set("除数不能为零")
				return
			}
			result = num1 / num2
		}
		operationResultData.Set(convertFromDecimal(result, base))
	}

	// 创建四则运算按钮
	addButton := widget.NewButton("+", func() { performOperation("+") })
	subtractButton := widget.NewButton("-", func() { performOperation("-") })
	multiplyButton := widget.NewButton("*", func() { performOperation("*") })
	divideButton := widget.NewButton("/", func() { performOperation("/") })

	operationContent := container.NewVBox(
		widget.NewLabel("四则运算区域"),
		opInput1,
		opInput2,
		container.NewHBox(
			widget.NewLabel("运算进制:"),
			operationBase,
		),
		container.NewHBox(
			addButton,
			subtractButton,
			multiplyButton,
			divideButton,
		),
		widget.NewLabel("运算结果:"),
		operationResult,
	)

	// 整体布局
	content := container.NewVBox(
		conversionContent,
		widget.NewSeparator(),
		operationContent,
	)

	// 设置窗口内容并显示
	w.SetContent(content)

	// 窗口大小设置
	w.Resize(fyne.NewSize(400, 400))
	// 窗口居中显示
	w.CenterOnScreen()

	// 显示窗口并运行应用
	w.ShowAndRun()
}
