package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image/color"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"unsafe"
)

type QueueItemData struct {
	Prefix string
	Content string
}

type QueueListView struct {
	mu					sync.Mutex
	List 				*widget.List
	BindData			*[]QueueItemData
	CurrentSelectedId  	widget.ListItemID
	HeaderView			*canvas.Text
	HeaderData			string
}

var listView = &QueueListView{
		BindData:&[]QueueItemData{},
		CurrentSelectedId: -1,
		HeaderData: "当前并无接待的人",
	}
var keyWindow *fyne.Window

var queueWindowPtr *fyne.Window
var queueWindowSetting = struct {
	IsFloating		bool
	IsHideDecorate	bool
}{}

func uiInit()  {
	a := app.NewWithID("com.queue.demo")
	a.Settings().SetTheme(&OpaqueTheme{})
	fmt.Println(fyne.CurrentApp().Storage().RootURI())

	queueWindow := queueWindow()
	mainWindow := fyneWindow(queueWindow,nil)
	codeWindow := codeInputWindow(mainWindow)

	(*codeWindow).ShowAndRun()

	//test code
	//mainWindow := fyneWindow(nil,nil)
	//codeWindow := codeInputWindow(mainWindow)
	//(*mainWindow).ShowAndRun()
}

func checkIfHaveFile(file string) bool {
	lists := fyne.CurrentApp().Storage().List()
	fmt.Println("FileList ",lists)
	for _, i2 := range lists {
		if i2 == file {
			return true
		}
	}
	return false
}

var isRememberCoder = false
func codeInputWindow(mainWindow *fyne.Window) *fyne.Window {
	w := fyne.CurrentApp().NewWindow("身份码")

	r, err := fyne.CurrentApp().Storage().Open("coderMember")
	if err == nil {
		str, _ := ioutil.ReadAll(r)
		defer r.Close()
		if string(str) == "1" {
			isRememberCoder = true
		}
	}else {
		fmt.Println("Read CoderMember Error ",err)
	}

	title := canvas.NewText("认证身份后可开启玩法", color.NRGBA{
		R: 255,
		G: 102,
		B: 153,
		A: 255,
	})
	title.TextSize = 18

	identifyLabel := canvas.NewText("身份码", color.Black)
	identifyInput := widget.NewEntry()
	identifyInput.SetPlaceHolder("请输入身份码")
	if isRememberCoder {
		coder, err := fyne.CurrentApp().Storage().Open("coder")
		if err == nil {
			str, _ := ioutil.ReadAll(coder)
			defer coder.Close()
			if len(string(str)) > 0 {
				identifyInput.SetText(string(str))
			}
		}else {
			fmt.Println("Open Coder Error ",err)
		}
	}

	hint := canvas.NewText("在获取推流地址后可获取身份码", color.NRGBA{
		R: 130,
		G: 130,
		B: 130,
		A: 255,
	})
	ul, _ := url.Parse("https://link.bilibili.com/p/center/index#/my-room/start-live")
	hintLink := widget.NewHyperlink("去获取", ul)

	startBtn := widget.NewButton("  开启玩法  ", func() {
		str := identifyInput.Text
		if len(str) <= 0 {
			dialog.ShowError(fmt.Errorf("请输入身份码"), w)
			return
		}
		err := RunBiliDanmu(str)
		if err != nil {
			dialog.ShowError(fmt.Errorf(err.Error()), w)
			return
		}

		if isRememberCoder {
			var coder fyne.URIWriteCloser
			var err error
			if checkIfHaveFile("coder") {
				coder, err = fyne.CurrentApp().Storage().Save("coder")
			}else {
				coder, err = fyne.CurrentApp().Storage().Create("coder")
			}
			if err == nil {
				coder.Write([]byte(identifyInput.Text))
				defer coder.Close()
			}else {
				fmt.Println("Write Coder Error ",err)
			}
		}

		w.Close()
		//打开界面
		(*mainWindow).Show()
	})

	rememberCheckBox := widget.NewCheck("记住身份码", func(b bool) {
		var rw fyne.URIWriteCloser
		var err error
		if checkIfHaveFile("coderMember") {
			fmt.Println("Have coder")
			rw, err = fyne.CurrentApp().Storage().Save("coderMember")
		}else {
			fmt.Println("Dont have coder")
			rw, err = fyne.CurrentApp().Storage().Create("coderMember")
		}
		if err != nil {
			dialog.ShowError(fmt.Errorf("无法访问记录状态"), w)
			return
		}
		if b {
			isRememberCoder = true
			rw.Write([]byte("1"))
		}else {
			isRememberCoder = false
			rw.Write([]byte("0"))
		}
		defer rw.Close()
	})

	if isRememberCoder {
		rememberCheckBox.SetChecked(true)
	}

	w.SetContent(container.NewVBox(
		container.New(layout.NewGridWrapLayout(fyne.NewSize(1,10))),
		container.NewCenter(title),
		container.New(layout.NewGridWrapLayout(fyne.NewSize(1,10))),
		container.NewBorder(nil,nil,identifyLabel,nil,identifyInput),
		container.New(layout.NewHBoxLayout(),hint, layout.NewSpacer(), hintLink),
		container.NewCenter(startBtn),
		container.NewCenter(rememberCheckBox),
		))

	return &w
}

func fyneWindow(qw*fyne.Window, listWindow *fyne.Window) *fyne.Window {

	w := fyne.CurrentApp().NewWindow("Control Panel")
	//glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)


	keyWindow = &w

	accept := widget.NewButton("接待一位", func() {
		if listView.List == nil {
			return
		}
		listView.mu.Lock()
		defer listView.mu.Unlock()
		if len(*listView.BindData) > 0 {
			//datalist, _ := listView.BindData.Get()
			datalist := *listView.BindData
			targetVal := datalist[0]
			targetList := datalist[1:]
			listView.List.Unselect(listView.CurrentSelectedId)
			listView.CurrentSelectedId = -1
			//listView.BindData.Set(targetList)
			listView.BindData = &targetList
			listView.List.Refresh()

			fmt.Println("接待贵宾 ", targetVal)
			listView.HeaderData = fmt.Sprintf("正在接待 ： %s%s",targetVal.Prefix,targetVal.Content)
			if listView.HeaderView != nil {
				listView.HeaderView.Text = listView.HeaderData
				listView.HeaderView.Refresh()
				fmt.Println("刷新头部")
			}
		}else {
			//d := dialog.NewInformation("提醒","暂时没有人排队中哦",*keyWindow)
			//d.Show()
			listView.HeaderData = fmt.Sprint("当前并无接待的人")
			if listView.HeaderView != nil {
				listView.HeaderView.Text = listView.HeaderData
				listView.HeaderView.Refresh()
			}
		}
	})


	groupLevelLimit := buildGroupLevelLimit()
	groupQueueCondition := buildGroupQueueCondition()
	rect := canvas.NewRectangle(color.White)
	rect.StrokeColor = color.White
	rect.StrokeWidth = 2
	//w.SetContent(container.NewBorder(nil, nil, nil, container.NewVBox(add, remove,move), list))
	layout.NewMaxLayout()
	w.SetContent(
		container.New(&FullLayout{},
			rect,
			container.NewVBox(
				groupLevelLimit,
				groupQueueCondition,
				buildQueueAction(),
				accept,
				buildQueueSettingPanel(),
			),
		),
	)

	w.SetMaster()
	return &w
}

func queueWindow() *fyne.Window {
	w := fyne.CurrentApp().NewWindow("Queue")
	glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
	if queueWindowSetting.IsFloating {
		glfw.WindowHint(glfw.Floating, glfw.True)
	}else {
		glfw.WindowHint(glfw.Floating, glfw.False)
	}

	if queueWindowSetting.IsHideDecorate {
		glfw.WindowHint(glfw.Decorated, glfw.False)
	}else {
		glfw.WindowHint(glfw.Decorated, glfw.True)
	}

	header := canvas.NewText(listView.HeaderData, QueueListItemColor.ContentColor)
	header.TextSize = 18
	header.TextStyle.Bold = true
	header.TextStyle.TabWidth = 3
	listView.HeaderView = header

	//data := binding.BindStringList(
	//	&[]string{},
	//)
	//data := []QueueItemData{
	//	{
	//		Prefix:  "1",
	//		Content: "测试",
	//	},
	//	{
	//		Prefix:  "2",
	//		Content: "测测试",
	//	},
	//}
	//
	//
	//listView.BindData = &data
	list:= widget.NewList(func() int {
		return len(*listView.BindData)
	}, func() fyne.CanvasObject {
		return NewQueueListItem("Text", "Next")
	}, func(id widget.ListItemID, object fyne.CanvasObject) {
		c, ok := object.(*fyne.Container)
		if ok {
			item := (*listView.BindData)[id]
			c0, ok := c.Objects[0].(*canvas.Text)
			if ok {
				c0.Text = item.Prefix
			}
			c1, ok := c.Objects[1].(*canvas.Text)
			if ok {
				c1.Text = item.Content
			}
			//fmt.Println("id ", id,"-", c.Objects)
		}
	})

	list.OnSelected = func(id widget.ListItemID) {
		fmt.Println("Selected ", id)
		listView.CurrentSelectedId = id
	}
	listView.List = list

	w.SetContent(
		container.NewBorder(container.NewCenter(header),nil,nil,nil,
			list,
		),
	)
	w.SetOnClosed(func() {
		queueWindowPtr = nil
		listView.HeaderView = nil
		listView.List = nil
	})
	//w.SetCloseIntercept(func() {
	//	(*queueWindowPtr).Hide()
	//})

	return &w
}

func GetUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

func buildGroupLevelLimit() *fyne.Container {
	title := canvas.NewText("仅允许以下用户排队", color.Black)
	title.TextSize = 16
	title.TextStyle.Bold = true

	td := widget.NewCheck("提督", func(b bool) {
		fmt.Println("提督 set to ", b)
		QueueFilterData.mu.Lock()
		defer QueueFilterData.mu.Unlock()
		if b {
			QueueFilterData.UserEnable |= UserTypeTiDu
		}else {
			QueueFilterData.UserEnable &= ^UserTypeTiDu
		}
	})
	jz := widget.NewCheck("舰长", func(b bool) {
		fmt.Println("舰长 set to ", b)
		QueueFilterData.mu.Lock()
		defer QueueFilterData.mu.Unlock()
		if b {
			QueueFilterData.UserEnable |= UserTypeJianZhang
		}else {
			QueueFilterData.UserEnable &= ^UserTypeJianZhang
		}
	})
	zd := widget.NewCheck("总督", func(b bool) {
		fmt.Println("总督 set to ", b)
		QueueFilterData.mu.Lock()
		defer QueueFilterData.mu.Unlock()
		if b {
			QueueFilterData.UserEnable |= UserTypeZongDu
		}else {
			QueueFilterData.UserEnable &= ^UserTypeZongDu
		}
	})
	fs := widget.NewCheck("粉丝等级 >= ", func(b bool) {
		fmt.Println("粉丝 set to ", b)
		QueueFilterData.mu.Lock()
		defer QueueFilterData.mu.Unlock()
		if b {
			QueueFilterData.UserEnable |= UserTypeFans
		}else {
			QueueFilterData.UserEnable &= ^UserTypeFans
		}
	})
	fsInput := widget.NewSelect([]string{"1","2","3","4","5","6","7","8","9","10"}, func(s string) {
		fmt.Println("Level Change to ", s)
		QueueFilterData.mu.Lock()
		defer QueueFilterData.mu.Unlock()
		v , _ := strconv.Atoi(s)
		QueueFilterData.UserEnableLevel = int64(v)
	})
	fsInput.SetSelectedIndex(0)

	return container.NewVBox(title, container.NewVBox(container.NewHBox(td, jz, zd), container.NewHBox(fs, fsInput)))
}

func buildGroupQueueCondition() *fyne.Container  {

	title := canvas.NewText("怎么排队", color.Black)
	title.TextSize = 16
	title.TextStyle.Bold = true

	content := widget.NewEntry()
	content.OnChanged = func(s string) {
		fmt.Println("内容变更 ",s)
		QueueFilterData.mu.Lock()
		defer QueueFilterData.mu.Unlock()

		QueueFilterData.QueueContent = s
	}
	content.SetPlaceHolder("弹幕内容")

	gift := widget.NewRadioGroup([]string{"发送弹幕","送礼","指定礼物(暂不支持)"}, func(s string) {
		fmt.Println("Queue Condition change to ",s)
		QueueFilterData.mu.Lock()
		defer QueueFilterData.mu.Unlock()

		content.SetText("")
		content.Enable()
		if s == "发送弹幕" {
			content.SetPlaceHolder("弹幕内容")
			QueueFilterData.QueueKind = QueueTypeDanmu
		}else if s == "送礼" {
			QueueFilterData.QueueKind = QueueTypeGift
			content.SetPlaceHolder("礼物价值，以电池为单位，例如大于100电池礼物，这里填写 100 即可")
		}else {
			QueueFilterData.QueueKind = QueueTypeSpecialGift
			content.SetPlaceHolder("不需要填写")
			content.Disable()
		}
	})
	se := QueueFilterData.QueueKind
	if se == QueueTypeDanmu {
		gift.SetSelected("发送弹幕")
	}else if se == QueueTypeGift{
		gift.SetSelected("送礼")
	}

	return container.NewVBox(title, container.NewHBox(gift), content)
}

func buildGroupQueueRule() *fyne.Container {
	title := canvas.NewText("排队规则", color.Black)
	title.TextSize = 16
	title.TextStyle.Bold = true

	cd := widget.NewCheck("是否文明排队", func(b bool) {
		fmt.Println("文明排队 set to ", b)
	})
	cd.SetChecked(true)

	return container.NewVBox(title, cd)
}

type ListControlContainer struct {
	Buttons []*widget.Button
}

func (l *ListControlContainer)enable(enable bool)  {
	for _, btn := range l.Buttons {
		if enable {
			btn.Enable()
		}else {
			btn.Disable()
		}

	}
}
var queueActionControl = &ListControlContainer{}
func buildQueueAction() *fyne.Container  {
	title := canvas.NewText("队列控制", color.Black)
	title.TextSize = 16
	title.TextStyle.Bold = true

	upTop := widget.NewButton("移到队头", func() {
		if listView.List == nil {
			return
		}

		if listView.CurrentSelectedId == 0 {
			return
		}

		if listView.CurrentSelectedId > 0 && listView.CurrentSelectedId < listView.List.Length() {
			listView.mu.Lock()
			defer listView.mu.Unlock()

			datalist := *listView.BindData
			curVal := datalist[listView.CurrentSelectedId]

			targetList := append([]QueueItemData{curVal}, datalist[0:listView.CurrentSelectedId]...)
			if listView.CurrentSelectedId != listView.List.Length() - 1 {
				targetList = append(targetList, datalist[listView.CurrentSelectedId + 1:]...)
			}

			listView.BindData = &targetList
			listView.List.Select(0)
		}else {
			//d := dialog.NewInformation("提醒","需要先选中一个人哦",*keyWindow)
			//d.Show()
		}
	})
	upOne := widget.NewButton("往上移动一位", func() {

		if listView.List == nil {
			return
		}

		if listView.CurrentSelectedId == 0 {
			return
		}
		if listView.CurrentSelectedId > 0 && listView.CurrentSelectedId < listView.List.Length() {
			listView.mu.Lock()
			defer listView.mu.Unlock()

			targetVal := (*listView.BindData)[listView.CurrentSelectedId - 1]
			curVal := (*listView.BindData)[listView.CurrentSelectedId]
			(*listView.BindData)[listView.CurrentSelectedId - 1] = curVal
			(*listView.BindData)[listView.CurrentSelectedId] = targetVal
			listView.List.Select(listView.CurrentSelectedId - 1)
		}else {
			//d := dialog.NewInformation("提醒","需要先选中一个人哦",*keyWindow)
			//d.Show()
		}
	})
	remove := widget.NewButton("移出队列", func() {
		if listView.List == nil {
			return
		}

		if listView.CurrentSelectedId >= 0 && listView.CurrentSelectedId < listView.List.Length() {
			listView.mu.Lock()
			defer listView.mu.Unlock()

			datalist := *listView.BindData

			targetList := append(datalist[0:listView.CurrentSelectedId], datalist[listView.CurrentSelectedId + 1:]...)
			listView.List.Unselect(listView.CurrentSelectedId)
			listView.CurrentSelectedId = -1
			listView.BindData = &targetList

		}else {
			//d := dialog.NewInformation("提醒","需要先选中一个人哦",*keyWindow)
			//d.Show()
		}
	})
	downOne := widget.NewButton("往下移动一位", func() {
		if listView.List == nil {
			return
		}

		lastIndex := listView.List.Length() - 1
		if listView.CurrentSelectedId == lastIndex {
			return
		}
		if listView.CurrentSelectedId >= 0 && listView.CurrentSelectedId < listView.List.Length() {
			listView.mu.Lock()
			defer listView.mu.Unlock()

			targetVal := (*listView.BindData)[listView.CurrentSelectedId + 1]
			curVal := (*listView.BindData)[listView.CurrentSelectedId]
			(*listView.BindData)[listView.CurrentSelectedId + 1] = curVal
			(*listView.BindData)[listView.CurrentSelectedId] = targetVal
			listView.List.Select(listView.CurrentSelectedId + 1)
		}else {
			//d := dialog.NewInformation("提醒","需要先选中一个人哦",*keyWindow)
			//d.Show()
		}
	})
	downBottom := widget.NewButton("移到队尾", func() {
		if listView.List == nil {
			return
		}

		lastIndex := listView.List.Length() - 1
		if listView.CurrentSelectedId == lastIndex {
			return
		}
		if listView.CurrentSelectedId >= 0 && listView.CurrentSelectedId < listView.List.Length() {
			listView.mu.Lock()
			defer listView.mu.Unlock()

			datalist := *listView.BindData
			curVal := datalist[listView.CurrentSelectedId]

			targetList := append(datalist[0:listView.CurrentSelectedId], datalist[listView.CurrentSelectedId + 1:]...)
			targetList = append(targetList, curVal)

			listView.BindData = &targetList
			listView.List.Select(lastIndex)
		}else {
			//d := dialog.NewInformation("提醒","需要先选中一个人哦",*keyWindow)
			//d.Show()
		}
	})

	queueActionControl.Buttons = []*widget.Button{upTop,upOne,remove,downOne,downBottom}

	return container.NewVBox(title,upTop, upOne, remove, downOne, downBottom)
}

func buildQueueSettingPanel() *fyne.Container {
	title := canvas.NewText("队列窗口设置", color.Black)
	title.TextSize = 16
	title.TextStyle.Bold = true

	colorRec := canvas.NewRectangle(themeBGColor)
	colorRec.SetMinSize(fyne.NewSize(25, 25))
	colorRec.StrokeWidth = 2
	colorRec.StrokeColor = color.Gray{Y: 188}
	colorBox := container.NewMax(
		widget.NewButton("", func() {
			d:= dialog.NewColorPicker("","", func(c color.Color) {
				colorRec.FillColor = c
				themeBGColor = c
				colorRec.Refresh()
				resizeUpdateQueueWindow()
				fmt.Println("change color")
			},*keyWindow)
			d.Advanced = true
			d.Show()
		}),
		colorRec,
		)
	bgColor := container.NewHBox(layout.NewSpacer() ,colorBox, widget.NewLabel("背景颜色"))

	colorPreRec := canvas.NewRectangle(themeBGColor)
	colorPreRec.SetMinSize(fyne.NewSize(25, 25))
	colorPreRec.StrokeWidth = 2
	colorPreRec.StrokeColor = color.Gray{Y: 188}
	colorPreBox := container.NewMax(
		widget.NewButton("", func() {
			d := dialog.NewColorPicker("","", func(c color.Color) {
				colorPreRec.FillColor = c
				colorPreRec.Refresh()

				QueueListItemColor.PrefixColor = c
				if (*queueWindowPtr != nil) {
					var px int
					var py int

					y := GetUnexportedField(reflect.ValueOf(*queueWindowPtr).Elem().FieldByName("viewport"))
					preW, ok := y.(*glfw.Window)
					if ok && preW != nil {
						px , py = preW.GetPos()
					}
					recordSize := (*queueWindowPtr).Content().Size()
					fmt.Println("Record Position ", px, " - ",py, " Size ", recordSize)
					(*queueWindowPtr).Close()
					showQueueWindow()
					go func() {
						(*queueWindowPtr).Resize(recordSize)
					}()
				}
			},*keyWindow)
			d.Advanced = true
			d.Show()
		}),
		colorPreRec,
	)
	preColor := container.NewHBox(layout.NewSpacer() ,colorPreBox, widget.NewLabel("称号颜色"))


	colorTxtRec := canvas.NewRectangle(themeBGColor)
	colorTxtRec.SetMinSize(fyne.NewSize(25, 25))
	colorTxtRec.StrokeWidth = 2
	colorTxtRec.StrokeColor = color.Gray{Y: 188}
	colorTxtBox := container.NewMax(
		widget.NewButton("", func() {
			d := dialog.NewColorPicker("","", func(c color.Color) {
				colorTxtRec.FillColor = c
				colorTxtRec.Refresh()

				QueueListItemColor.ContentColor = c
				if (*queueWindowPtr != nil) {
					var px int
					var py int

					y := GetUnexportedField(reflect.ValueOf(*queueWindowPtr).Elem().FieldByName("viewport"))
					preW, ok := y.(*glfw.Window)
					if ok && preW != nil {
						px , py = preW.GetPos()
					}
					recordSize := (*queueWindowPtr).Content().Size()
					fmt.Println("Record Position ", px, " - ",py, " Size ", recordSize)
					(*queueWindowPtr).Close()
					showQueueWindow()
					go func() {
						(*queueWindowPtr).Resize(recordSize)
					}()
				}
			},*keyWindow)
			d.Advanced = true
			d.Show()
		}),
		colorTxtRec,
	)
	txtColor := container.NewHBox(layout.NewSpacer() ,colorTxtBox, widget.NewLabel("文字颜色"))

	changeQueueWindowOpaque := widget.NewCheck("是否透明", func(b bool) {
		ifOpenOpaque = b
		if b {
			fmt.Println("打开透明")
		}else {
			fmt.Println("关闭透明")
		}
		resizeUpdateQueueWindow()
	})


	showQueue := widget.NewButton("显示队列页面", showQueueWindow)

	//test code
	//addNum := 0
	//testAddOne := widget.NewButton("Add one", func() {
	//	addToQueue(QueueItemData{fmt.Sprint(addNum),"asdf"})
	//	//addNum ++
	//})

	return container.NewVBox(title,container.NewHBox(bgColor, preColor, txtColor),
		showQueue,
		changeQueueWindowOpaque,
		//testAddOne,
		)
}


func resizeUpdateQueueWindow()  {
	if queueWindowPtr != nil{
		size := (*queueWindowPtr).Canvas().Size()
		(*queueWindowPtr).Resize(fyne.NewSize(size.Width + 1, size.Height + 1))
		(*queueWindowPtr).Resize(size)
	}
}

func showQueueWindow()  {
	if queueWindowPtr == nil {
		queueWindowPtr = queueWindow()
		(*queueWindowPtr).Show()
	}
}