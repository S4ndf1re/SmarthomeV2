package main

import (
	"Smarthome/gui"
	"Smarthome/user"
	"fmt"
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

func newUserContainer() *gui.Container {

	data := gui.NewData("new_user_data", gui.NewAlert("new_user_alert", "", "info"))

	onInit := func(user string) {
	}

	onUnload := func(user string) {
		data.Update(gui.NewAlert("new_user_alert", "", "info"))
	}

	container := gui.NewContainer("default_new_user", "Create User", onInit, onUnload)

	container.Add(data)

	username := ""
	password := ""

	textUsername := gui.NewTextField("new_user_username", "Username", func(user string, text string) {
		username = text
	})
	container.Add(textUsername)

	textPassword := gui.NewTextField("new_user_password", "Password", func(user string, text string) {
		password = text
	})
	container.Add(textPassword)

	submit := gui.NewButton("new_user_submit", "Submit", func(userReq string) {
		if _, err := user.Load(username); err == nil {
			data.Update(gui.NewAlert("new_user_alert", "User already exists", "error"))
			return
		}

		userCreated := user.New(username, password)
		if err := userCreated.Store(); err != nil {
			data.Update(gui.NewAlert("new_user_alert", fmt.Sprintf("%s", err.Error()), "error"))
			return
		}

		data.Update(gui.NewAlert("new_user_alert", fmt.Sprintf("User %s successfully created", username), "success"))

	})
	container.Add(submit)

	return container

}
