package certark

import "github.com/josark2005/certark/ark"

// check acme user validity
func checkAcmeUserValidity() {
	ark.Debug().Msg("Prepare to check acme user validity")
	for acmeUser, acme := range AcmeUsers {
		ark.Debug().Str("acme", acmeUser).Msg("Checking acme user validity")

		// check enable
		if !acme.Enabled {
			ark.Debug().Str("acme", acmeUser).Msg("Acme user is disabled")
			AcmeUsersNotValid[acmeUser] = acme
			delete(AcmeUsers, acmeUser)
			continue
		}

		// check account and key
		if acme.Email == "" {
			ark.Warn().Str("acme", acmeUser).Str("reason", "acme email is empty").Msg("Acme user is invalid")
			AcmeUsersNotValid[acmeUser] = acme
			delete(AcmeUsers, acmeUser)
			continue
		}
		// check private key
		if acme.PrivateKey == "" {
			ark.Warn().Str("acme", acmeUser).Str("reason", "acme account privatekey is empty").Msg("Acme user is invalid")
			AcmeUsersNotValid[acmeUser] = acme
			delete(AcmeUsers, acmeUser)
			continue
		}
	}
}

// check task validity
func checkTaskValidity() {
	ark.Debug().Msg("Prepare to check task validity")
	for taskName, task := range Tasks {
		ark.Debug().Str("task", taskName).Msg("Checking task validity")

		// check acme user in tasks
		_, acmeUserExists := AcmeUsers[task.AcmeUser]
		if task.AcmeUser == "" || !acmeUserExists || AcmeUsers[task.AcmeUser].PrivateKey == "" {
			ark.Warn().Str("task", taskName).Str("acme", task.AcmeUser).Str("reason", "acme user is invalid").Msg("Task is invalid")
			TaskNotValid[taskName] = task
			delete(Tasks, taskName)
			continue
		}

		// check dns profile
		if task.DnsProfile == "" {
			ark.Warn().Str("task", taskName).Str("reason", "empty dns profile").Msg("Task is invalid")
			TaskNotValid[taskName] = task
			delete(Tasks, taskName)
			continue
		}
	}
}
