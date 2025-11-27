package command

import (
	"os"
)

type trapService struct {
	managerTrap      chan os.Signal
	chatTrap         chan os.Signal
	notificationTrap chan os.Signal
	adminTrap        chan os.Signal
}

func setupTrapService() trapService {
	return trapService{
		managerTrap:      make(chan os.Signal, 1),
		chatTrap:         make(chan os.Signal, 1),
		notificationTrap: make(chan os.Signal, 1),
		adminTrap:        make(chan os.Signal, 1),
	}
}

func (t trapService) sendSignal(sig os.Signal) {
	t.managerTrap <- sig
	t.chatTrap <- sig
	t.notificationTrap <- sig
	t.adminTrap <- sig
}
