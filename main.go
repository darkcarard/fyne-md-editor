package main

import (
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type appConfig struct {
	editWidget    *widget.Entry
	previewWidget *widget.RichText
	currentFile   fyne.URI
	saveMenuItem  *fyne.MenuItem
}

var config appConfig

func main() {
	application := app.NewWithID("mdeditor-20250107")
	window := application.NewWindow("Markdown Editor")
	editWidget, previewWidget := config.makeUI()
	config.createMenuItems(window)
	window.SetContent(container.NewHSplit(editWidget, previewWidget))
	window.Resize(fyne.NewSize(800, 500))
	window.CenterOnScreen()

	window.ShowAndRun()

}

func (config *appConfig) makeUI() (*widget.Entry, *widget.RichText) {
	editWidget := widget.NewMultiLineEntry()
	previewWidget := widget.NewRichTextFromMarkdown("")
	config.editWidget = editWidget
	config.previewWidget = previewWidget
	editWidget.OnChanged = previewWidget.ParseMarkdown

	return editWidget, previewWidget
}

func (config *appConfig) createMenuItems(window fyne.Window) {
	openItem := fyne.NewMenuItem("Open...", config.openFunc(window))
	saveItem := fyne.NewMenuItem("Save", config.saveFunc(window))
	config.saveMenuItem = saveItem
	config.saveMenuItem.Disabled = true
	saveAsItem := fyne.NewMenuItem("Save as...", config.saveAsFunc(window))
	fileMenu := fyne.NewMenu("File", openItem, saveItem, saveAsItem)
	menu := fyne.NewMainMenu(fileMenu)

	window.SetMainMenu(menu)
}

func (config *appConfig) saveFunc(window fyne.Window) func() {
	return func() {
		if config.currentFile != nil {
			write, err := storage.Writer(config.currentFile)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			defer write.Close()

			write.Write([]byte(config.editWidget.Text))
		}
	}
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (config *appConfig) openFunc(window fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if read == nil {
				return
			}

			defer read.Close()

			data, err := io.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			config.editWidget.SetText(string(data))
			config.currentFile = read.URI()
			window.SetTitle(window.Title() + " - " + read.URI().Name())
			config.saveMenuItem.Disabled = false

		}, window)

		openDialog.SetFilter(filter)
		openDialog.Show()
	}
}

func (config *appConfig) saveAsFunc(window fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}

			if write == nil {
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Invalid file extension!", window)
				return
			}

			defer write.Close()

			write.Write([]byte(config.editWidget.Text))
			config.currentFile = write.URI()

			window.SetTitle(window.Title() + " - " + write.URI().Name())
			config.saveMenuItem.Disabled = false
		}, window)

		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)

		saveDialog.Show()
	}
}
