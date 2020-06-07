package noodle

import (
	"errors"
	"syscall/js"
)

type Image struct {
	data js.Value
}

//LoadImage loads a new image
func LoadImage(url string) (*Image, error) {
	ch := make(chan error, 1)
	img := js.Global().Get("Image").New()

	//Prepare the events
	loadEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() { ch <- nil }()
		return nil
	})
	defer loadEvent.Release()
	img.Call("addEventListener", "load", loadEvent)

	errorEvent := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() { ch <- errors.New("Failed to load image") }()
		return nil
	})
	defer errorEvent.Release()
	img.Call("addEventListener", "error", errorEvent)

	//Set the source
	img.Set("src", url)

	//Wait for the source to load
	err := <-ch
	if err != nil {
		return nil, err
	}

	//Finish
	return &Image{img}, nil
}

//Data gets the JS value
func (i *Image) Data() js.Value {
	return i.data
}

//Width gets the width in pixels
func (i *Image) Width() int {
	return i.data.Get("width").Int()
}

//Height gets the height in pixels
func (i *Image) Height() int {
	return i.data.Get("height").Int()
}
