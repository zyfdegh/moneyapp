package app

import (
	"log"
	"net/url"
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
	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon("主页", theme.HomeIcon(), this.welcomeTab(main)),
		widget.NewTabItemWithIcon("记账", theme.DocumentCreateIcon(), this.writeTab(main)),
		widget.NewTabItemWithIcon("查账", theme.SearchIcon(), this.queryTab(main)),
		widget.NewTabItemWithIcon("我的", theme.HomeIcon(), this.myTab(main)),
		widget.NewTabItemWithIcon("关于", theme.InfoIcon(), this.aboutTab(main)),
	)
	tabs.SetTabLocation(widget.TabLocationLeading)
	main.SetContent(tabs)
	main.CenterOnScreen()
	main.SetMaster()
	// TODO load session from prefrence
	if this.session == nil {
		// 登录
		main.Hide()
		this.launchLoginPanel(main)
	}
	main.ShowAndRun()
}

func (this *App) launchLoginPanel(father fyne.Window) {
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
			this.session = sess
			win.Close()
			father.Show()
		},
	}
	loginForm.Append("用户名：", username)
	loginForm.Append("密码：", password)

	win.CenterOnScreen()
	win.SetContent(widget.NewVBox(
		widget.NewLabel("欢迎使用"+consts.AppName+"，请登录"),
		layout.NewSpacer(),
		loginForm,
		layout.NewSpacer(),
		widget.NewButton("退出", func() {
			this.app.Quit()
		}),
	))
	win.Show()
}

func (this *App) welcomeTab(father fyne.Window) fyne.CanvasObject {
	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return widget.NewVBox(
		widget.NewLabelWithStyle("Welcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewHyperlinkWithStyle("fyne.io", link, fyne.TextAlignCenter, fyne.TextStyle{}),
		layout.NewSpacer(),

		widget.NewGroup("Theme",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewButton("Dark", func() {
					this.app.Settings().SetTheme(theme.DarkTheme())
				}),
				widget.NewButton("Light", func() {
					this.app.Settings().SetTheme(theme.LightTheme())
				}),
			),
		),
	)
}

func (this *App) writeTab(father fyne.Window) fyne.CanvasObject {
	return widget.NewVBox()
}
func (this *App) queryTab(father fyne.Window) fyne.CanvasObject {
	return widget.NewVBox()
}

func (this *App) myTab(father fyne.Window) fyne.CanvasObject {
	return widget.NewVBox()
}
func (this *App) aboutTab(father fyne.Window) fyne.CanvasObject {
	return widget.NewVBox()
}
