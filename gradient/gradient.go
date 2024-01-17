package gradient

import (
	"time"

	"nyiyui.ca/halation/aiz"
	"nyiyui.ca/halation/timeutil"
)

func init() {
	aiz.GradientTypes["nyiyui.ca/halation/gradient.LinearGradient"] = func() aiz.Gradient { return new(LinearGradient) }
}

type LinearGradient struct {
	Duration_            timeutil.Duration `json:"duration"`
	PreferredResolution_ timeutil.Duration `json:"preferredResolution"`
}

func (l *LinearGradient) Duration() time.Duration { return time.Duration(l.Duration_) }

func (l *LinearGradient) PreferredResolution() time.Duration {
	return time.Duration(l.PreferredResolution_)
}

func (l *LinearGradient) ValueAt(i time.Duration) float32 {
	return float32(i) / float32(l.Duration_)
}

func (l *LinearGradient) Values(resolution time.Duration) []float32 {
	values := make([]float32, 0, time.Duration(l.Duration_)/resolution)
	for i := time.Duration(0); i <= time.Duration(l.Duration_); i += resolution {
		values = append(values, l.ValueAt(i))
	}
	return values
}

func (l *LinearGradient) TypeName() string { return "nyiyui.ca/halation/gradient.LinearGradient" }
