package main

import (
	"Smarthome/gui"
	"Smarthome/user"
)

func passwordChangeContainer() *gui.Container {
	oldPW := ""
	newPW := ""
	repeatPW := ""

	textOldPW := gui.NewTextField("pw_change_old_pw", "Old Password", func(username string, text string) {
		oldPW = text
	})
	textNewPw1 := gui.NewTextField("pw_change_new_pw_1", "New Password", func(username string, text string) {
		newPW = text
	})
	textNewPw2 := gui.NewTextField("pw_change_new_pw_2", "Repeat", func(username string, text string) {
		repeatPW = text
	})
	data := gui.NewData("pw_change_response", gui.NewAlert("pw_change_alert", "", "info"))

	submit := gui.NewButton("pw_change_submit", "Submit", func(username string) {
		if newPW != repeatPW {
			data.Update(gui.NewAlert("pw_change_alert", "New Password and Repeat must be equal", "error"))
			return
		}

		refUser, err := user.Load(username)
		if err != nil {
			data.Update(gui.NewAlert("pw_change_alert", err.Error(), "error"))
			return
		}

		newPwUser := user.New(username, oldPW)
		if !refUser.Equals(newPwUser) {
			data.Update(gui.NewAlert("pw_change_alert", "Old Password is not correct", "error"))
			return
		}

		newPwUser = user.New(username, newPW)
		if refUser.Remove() != nil {
			data.Update(gui.NewAlert("pw_change_alert", "Internal error. Could not change password", "error"))
			return
		}

		if newPwUser.Store() != nil {
			data.Update(gui.NewAlert("pw_change_alert", "Internal error. Could not change password", "error"))
			return
		}

		data.Update(gui.NewAlert("pw_change_alert", "Changed Password successfully", "success"))
	})

	container := gui.NewContainer("pw_change", "Change Password", func(s string) {
	}, func(s string) {
		data.Update(gui.NewAlert("pw_change_alert", "", "info"))
	})
	container.Add(textOldPW)
	container.Add(textNewPw1)
	container.Add(textNewPw2)
	container.Add(submit)
	container.Add(data)

	return container
}

// TODO (Jan): defaultNewUserContainer to create a new user. Report errors via data if any occur
