package gradient

import "nyiyui.ca/halation/aiz"

func init() {
	aiz.GradientTypes["nyiyui.ca/halation/gradient.LinearGradient"] = func() aiz.Gradient { return new(LinearGradient) }
}

type LinearGradient struct {
	Duration_            int `json:"duration"`
	PreferredResolution_ int `json:"preferredResolution"`
}

func (l *LinearGradient) Duration() int { return l.Duration_ }

func (l *LinearGradient) PreferredResolution() int { return l.PreferredResolution_ }

func (l *LinearGradient) ValueAt(i int) float32 {
	return float32(i) / float32(l.Duration_)
}

func (l *LinearGradient) Values(resolution int) []float32 {
	values := make([]float32, 0, l.Duration_/resolution)
	for i := 0; i < l.Duration_; i += resolution {
		values = append(values, l.ValueAt(i))
	}
	return values
}

func (l *LinearGradient) TypeName() string { return "nyiyui.ca/halation/gradient.LinearGradient" }
