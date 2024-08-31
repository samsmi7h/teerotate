package teerotate

type hooks struct {
	postRotation func()
}

func (r *RotatingLogger) WithPostRotationHook(hook func()) {
	r.hooks.postRotation = hook
}
