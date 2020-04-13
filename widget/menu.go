package widget

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/internal/widget"
	"fyne.io/fyne/theme"
)

// NewPopUpMenuAtPosition creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be positioned at the provided location and shown as an overlay on the specified canvas.
func NewPopUpMenuAtPosition(menu *fyne.Menu, c fyne.Canvas, pos fyne.Position) *PopUp {
	options := NewVBox()
	for _, item := range menu.Items {
		if item.IsSeparator {
			options.Append(newSeparator())
		} else {
			options.Append(newMenuItemWidget(item))
		}
	}
	pop := newPopUp(options, c)
	pop.NotPadded = true
	focused := c.Focused()
	for _, o := range options.Children {
		if item, ok := o.(*menuItemWidget); ok {
			item.DismissAction = func() {
				if c.Focused() == nil {
					c.Focus(focused)
				}
				pop.Hide()
			}
		}
	}
	pop.ShowAtPosition(pos)
	return pop
}

// NewPopUpMenu creates a PopUp widget populated with menu items from the passed menu structure.
// It will automatically be shown as an overlay on the specified canvas.
func NewPopUpMenu(menu *fyne.Menu, c fyne.Canvas) *PopUp {
	return NewPopUpMenuAtPosition(menu, c, fyne.NewPos(0, 0))
}

type menuItemWidget struct {
	BaseWidget
	DismissAction func()
	Item          *fyne.MenuItem

	hovered bool
}

func (t *menuItemWidget) Tapped(*fyne.PointEvent) {
	t.Item.Action()
	if t.DismissAction != nil {
		t.DismissAction()
	}
}

func (t *menuItemWidget) CreateRenderer() fyne.WidgetRenderer {
	text := canvas.NewText(t.Item.Label, theme.TextColor())
	return &menuItemWidgetRenderer{widget.NewBaseRenderer([]fyne.CanvasObject{text}), text, t}
}

// MouseIn is called when a desktop pointer enters the widget
func (t *menuItemWidget) MouseIn(*desktop.MouseEvent) {
	t.hovered = true
	t.Refresh()
}

// MouseOut is called when a desktop pointer exits the widget
func (t *menuItemWidget) MouseOut() {
	t.hovered = false
	t.Refresh()
}

// MouseMoved is called when a desktop pointer hovers over the widget
func (t *menuItemWidget) MouseMoved(*desktop.MouseEvent) {
}

func newMenuItemWidget(item *fyne.MenuItem) *menuItemWidget {
	ret := &menuItemWidget{Item: item}
	ret.ExtendBaseWidget(ret)
	return ret
}

type menuItemWidgetRenderer struct {
	widget.BaseRenderer
	text *canvas.Text
	w    *menuItemWidget
}

func (r *menuItemWidgetRenderer) Layout(size fyne.Size) {
	padding := r.padding()
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(padding.Width/2, padding.Height/2))
}

func (r *menuItemWidgetRenderer) MinSize() fyne.Size {
	return r.text.MinSize().Add(r.padding())
}

func (r *menuItemWidgetRenderer) Refresh() {
	if r.text.TextSize != theme.TextSize() {
		defer r.Layout(r.w.Size())
	}
	r.text.TextSize = theme.TextSize()
	r.text.Color = theme.TextColor()
	canvas.Refresh(r.text)
}

func (r *menuItemWidgetRenderer) BackgroundColor() color.Color {
	if r.w.hovered {
		return theme.HoverColor()
	}

	return color.Transparent
}

func (r *menuItemWidgetRenderer) padding() fyne.Size {
	return fyne.NewSize(theme.Padding()*4, theme.Padding()*2)
}

func newSeparator() fyne.CanvasObject {
	s := canvas.NewRectangle(theme.DisabledTextColor())
	s.SetMinSize(fyne.NewSize(1, 2))
	return s
}
