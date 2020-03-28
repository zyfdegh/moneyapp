package app

import (
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
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

	sess *models.Session

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
	this.LaunchMainPanel()
}

func (this *App) LaunchMainPanel() {
	// TODO load session from prefrence
	main := this.app.NewWindow("Hello")
	main.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			this.app.Quit()
		})))
	main.SetMaster()
	main.Hide()
	this.LaunchLoginPanel(main)
}

func (this *App) LaunchLoginPanel(father fyne.Window) {
	win := this.app.NewWindow(consts.AppName)

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
			log.Println("[Login] Name:", username.Text)
			prog := dialog.NewProgressInfinite("登录中", "请稍后...", win)
			prog.Show()

			sess, err := this.userApi.Login(username.Text, password.Text)
			if err != nil {
				prog.Hide()
				dialog.ShowError(err, win)
				return
			}
			// TODO store session to prefrence
			this.sess = sess
			win.Close()
			father.Show()
		},
	}
	loginForm.Append("用户名：", username)
	loginForm.Append("密码：", password)

	win.SetContent(widget.NewVBox(
		widget.NewLabel("欢迎使用"+consts.AppName+"，请登录"),
		layout.NewSpacer(),
		loginForm,
		layout.NewSpacer(),
		widget.NewButton("退出", func() {
			this.app.Quit()
		}),
	))
	win.ShowAndRun()
}
