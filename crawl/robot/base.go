package robot

import . "../base"

type RobotBase struct {
}

func (p *RobotBase) Can(id string, task int32) bool {
	return false
}

type RealtimeTickRes struct {
	Id string
	RealtimeTick
}

func (p *RobotBase) GetRealtimeTick(ids string) []RealtimeTickRes {
	return nil
}
