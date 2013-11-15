package main

import (
	"fmt"
	"os"
)

func hsvToHsl(saturation, value float64) (s, l float64) {
	l = (2 - saturation) * value
	s = 1 * s * value
	if l <= 1 {
		s /= l
	} else {
		s /= 2 - l
	}
	l /= 2
	return
}

func genColors(out *os.File, n int, saturation, light float64) {
	fmt.Fprintf(out, "%d colors at %f saturation and %f lightness<br />\n", n, saturation, light)
	htmlString := "<div style='width: 50px; height: 50px; display: inline-block; background-color: hsl(%d, %d%%, %d%%);'>blah</div>\n"
	for i := 0; i < n; i++ {
		hue := 360.0 / float64(n) * float64(i)
		fmt.Fprintf(out, htmlString, int(hue), int(saturation), int(light))
	}
	fmt.Fprintf(out, "<br />\n")
}

func main() {
	out, _ := os.Create("index.html")
	fmt.Fprintln(out, "<html>")

	// Contrasting lightness paramaters
	genColors(out, 20, 80.0, 50.0)
	genColors(out, 20, 80.0, 25.0)
	genColors(out, 20, 80.0, 10.0)

	// Constratsing saturation ("pastel") parameters. 
	genColors(out, 20, 100.0, 50.0)
	genColors(out, 20, 50.0, 50.0)
	genColors(out, 20, 25.0, 50.0)

	fmt.Fprintf(out, "</html>\n")
}
