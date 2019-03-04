package apis

import (
	"fmt"
	"image/png"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/qpliu/qrencode-go/qrencode"
)

func QrHander() *martini.ClassicMartini {
	m := martini.Classic()
	m.Get("/:str", func(params martini.Params, w http.ResponseWriter) {
		text := params["str"]
		if text == "" {
			fmt.Fprintf(w, "form value text is required")
			return
		}

		w.Header().Add("Content-Type", "image/png")
		grid, err := qrencode.Encode(text, qrencode.ECLevelQ)
		if err != nil {
			fmt.Fprintf(w, "QRCode error:%s", err)
			return
		}
		png.Encode(w, grid.Image(8))
	})
	return m
}
