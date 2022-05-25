package handler

func Init() error {
	err := RegisterPlayerHandler()
	if err != nil {
		return err
	}
	err = RegisterRoomHandler()
	if err != nil {
		return err
	}

	return nil
}
