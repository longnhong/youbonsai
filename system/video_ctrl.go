package system

import (
	ov "LongTM/basic/o/video"
	"LongTM/basic/x/utils"
)

func (c *VideoWorker) VideoWorking(action *VideoAction) error {
	if action == nil {
		return nil
	}
	defer utils.Recover()
	defer action.Done()
	var video *ov.Video
	action.handlerAction(video)
	var err = action.GetError()
	return err
}
