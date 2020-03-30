package app

import (
	"errors"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/zyfdegh/moneyapp/consts"
	"github.com/zyfdegh/moneyapp/models"
	"github.com/zyfdegh/moneyapp/services"
)

func init() {
	os.Setenv("FYNE_FONT", "assets/fangzheng.ttf")
}

type App struct {
	app fyne.App

	session *models.Session

	userApi UserAPI
	monyApi MoneyAPI
}

func New() *App {
	return &App{
		app:     app.NewWithID("com.zyfdegh.moneyapp"),
		userApi: services.NewFakeUserAPI(),
	}
}

func (this *App) Run() {
	this.launchMainPanel()
}

func (this *App) launchMainPanel() {
	main := this.app.NewWindow(consts.AppName)
	welcomeTab := this.welcomeTab(main)
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("主页", theme.HomeIcon(), welcomeTab),
		widget.NewTabItemWithIcon("记账", theme.DocumentCreateIcon(), this.writeTab(main)),
		widget.NewTabItemWithIcon("查账", theme.SearchIcon(), this.queryTab(main)),
		widget.NewTabItemWithIcon("我的", theme.HomeIcon(), this.myTab(main)),
		widget.NewTabItemWithIcon("关于", theme.InfoIcon(), this.aboutTab(main)),
	)
	tabs.SetTabLocation(widget.TabLocationLeading)
	main.SetContent(tabs)
	main.CenterOnScreen()
	main.SetMaster()
	main.Hide()
	// TODO load session from prefrence
	if this.session == nil {
		// 登录
		go func() {
			this.launchSigninPanel(main)
			welcomeTab.Refresh()
		}()
	}
	main.ShowAndRun()
}

func (this *App) launchSigninPanel(parent fyne.Window) {
	win := this.app.NewWindow("登录")

	username := widget.NewEntry()
	username.SetPlaceHolder("输入账号")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("输入密码")

	loginForm := &widget.Form{
		OnCancel: func() {
			username.SetText("")
			password.SetText("")
		},
		OnSubmit: func() {
			if len(username.Text) == 0 || len(password.Text) == 0 {
				return
			}
			log.Println("[Signin] Name:", username.Text)
			prog := dialog.NewProgressInfinite("登录中", "请稍后...", win)
			prog.Show()

			sess, err := this.userApi.Login(username.Text, password.Text)
			if err != nil {
				prog.Hide()
				dialog.ShowError(err, win)
				return
			}
			// TODO store session to prefrence
			this.session = sess
			win.Close()
			parent.Show()
		},
	}
	loginForm.Append("用户名：", username)
	loginForm.Append("密码：", password)

	win.CenterOnScreen()
	win.SetContent(widget.NewVBox(
		widget.NewLabel("你好，请登录"),
		layout.NewSpacer(),
		loginForm,
		widget.NewHBox(
			widget.NewLabel("尚无账号?点此"),
			widget.NewButton("注册", func() {
				win.Hide()
				this.launchSignupPanel(win)
			}),
		),
		layout.NewSpacer(),
		widget.NewButton("退出", func() {
			this.app.Quit()
		}),
	))
	win.Show()
}

func (this *App) launchSignupPanel(parent fyne.Window) {
	win := this.app.NewWindow("注册")

	username := widget.NewEntry()
	username.SetPlaceHolder("数字、字母、下划线")
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("数字、字母，符号，6位以上")
	passwordTwice := widget.NewPasswordEntry()
	passwordTwice.SetPlaceHolder("再输入一次")
	realname := widget.NewEntry()
	realname.SetPlaceHolder("张三")
	cellphone := widget.NewEntry()
	cellphone.SetPlaceHolder("手机号")

	signupForm := &widget.Form{
		OnCancel: func() {
			username.SetText("")
			password.SetText("")
			passwordTwice.SetText("")
			realname.SetText("")
			cellphone.SetText("")
		},
		OnSubmit: func() {
			if len(username.Text) == 0 || len(password.Text) == 0 {
				return
			}
			if password.Text != passwordTwice.Text {
				err := errors.New("两次密码不一样")
				dialog.ShowError(err, win)
				return
			}
			log.Println("[Signup] Name:", username.Text)
			prog := dialog.NewProgressInfinite("注册中", "请稍后...", win)
			prog.Show()

			err := this.userApi.Register(&models.User{
				Username: username.Text,
				Password: password.Text,
				Realname: realname.Text,
				Cellphone: cellphone.Text,
			})
			if err != nil {
				prog.Hide()
				dialog.ShowError(err, win)
				return
			}
			win.Close()
			parent.Show()
		},
	}
	signupForm.Append("用户名：", username)
	signupForm.Append("密码：", password)
	signupForm.Append("密码再输一次：", passwordTwice)
	signupForm.Append("真实姓名：", realname)
	signupForm.Append("手机号：", cellphone)

	win.CenterOnScreen()
	win.SetContent(widget.NewVBox(
		widget.NewLabel("你好，注册新帐号"),
		layout.NewSpacer(),
		signupForm,
	))
	win.Show()
}

func (this *App) welcomeTab(parent fyne.Window) fyne.CanvasObject {
	username := "游客"
	if this.session != nil && len(this.session.Username) > 0 {
		username = this.session.Username
	}
	return widget.NewVBox(
		widget.NewLabelWithStyle(username+", 欢迎使用"+consts.AppName, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
	)
}

func (this *App) writeTab(parent fyne.Window) fyne.CanvasObject {
	return widget.NewVBox()
}
func (this *App) queryTab(parent fyne.Window) fyne.CanvasObject {
	return widget.NewVBox()
}

func (this *App) myTab(parent fyne.Window) fyne.CanvasObject {
	return widget.NewVBox()
}
func (this *App) aboutTab(parent fyne.Window) fyne.CanvasObject {
	return widget.NewVBox()
}
