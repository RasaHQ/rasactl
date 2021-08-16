package rasax

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

func (r *RasaX) progressBarBytes(maxBytes int64, description ...string) *progressbar.ProgressBar {
	desc := ""
	if len(description) > 0 {
		desc = description[0]
	}
	bar := progressbar.NewOptions64(
		maxBytes,
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(52),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionClearOnFinish(),
	)
	if err := bar.RenderBlank(); err != nil {
		r.Log.V(1).Error(err, "Can't render a progress bar")
		return nil
	}
	return bar
}
