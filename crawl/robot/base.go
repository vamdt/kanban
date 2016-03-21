package robot

type RobotBase struct {
}

func (p *RobotBase) Can(id string, task int32) bool {
	return false
}

func (p *RobotBase) RealtimeTick(ids string) {
}
