package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/jokin1999/certark/acme"
	"github.com/jokin1999/certark/acme/drivers"
	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// check if task profile exists
func checkTaskProfileExists(taskname string) bool {
	res := certark.FileOrDirExists(certark.TaskConfigDir + "/" + taskname)
	if res {
		ark.Debug().Msg("Task profile exists")
	} else {
		ark.Debug().Msg("Task profile does not exist")
	}
	return res
}

func init() {
	// task main command
	var taskCmd = cmdTask()

	// task ls
	taskCmd.AddCommand(cmdTaskLs())

	// task show
	taskCmd.AddCommand(cmdTaskShow())

	// task add
	taskCmd.AddCommand(cmdTaskAdd())

	// task append
	taskCmd.AddCommand(cmdTaskAppend())

	// task subtract
	taskCmd.AddCommand(cmdTaskSubtract())

	// task acme
	taskCmd.AddCommand(cmdTaskSetAcmeUser())

	// task set
	taskCmd.AddCommand(cmdTaskSet())

	// task run
	taskCmd.AddCommand(cmdTaskRun())

	//TODO - Delete task

	rootCmd.AddCommand(taskCmd)
}

// task command
func cmdTask() *cobra.Command {
	return &cobra.Command{
		Use:   "task",
		Short: "Task profiles management",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
}

// task ls
func cmdTaskLs() *cobra.Command {
	return &cobra.Command{
		Use:   "ls",
		Short: "List task profiles",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			listTasks()
		},
	}
}

// task show
func cmdTaskShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show [TASK]",
		Short: "Show a task profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				task := args[0]
				showTaskProfile(task)
			} else {
				cmd.Help()
			}
		},
	}
}

// task add
func cmdTaskAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add [TASK]",
		Short: "Add a task profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				task := args[0]
				addTaskProfile(task)
			} else {
				cmd.Help()
			}
		},
	}

	return c
}

// task append command
func cmdTaskAppend() *cobra.Command {
	c := &cobra.Command{
		Use:   "append [TASK] [DOAMIN]",
		Short: "Append domains in a task profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 1 {
				task := args[0]
				appendDomainTaskProfile(task, args[1:])
			} else {
				cmd.Help()
			}
		},
	}
	return c
}

// task subtract command
func cmdTaskSubtract() *cobra.Command {
	c := &cobra.Command{
		Use:   "sub [TASK] [DOAMIN]",
		Short: "Subtract a domain in a task profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 1 {
				task := args[0]
				domain := args[1]
				subtractDomainTaskProfile(task, domain)
			} else {
				cmd.Help()
			}
		},
	}
	return c
}

// task acme command
func cmdTaskSetAcmeUser() *cobra.Command {
	c := &cobra.Command{
		Use:   "acme [TASK] [ACME_ACCOUNT]",
		Short: "Set a acme user account in a task profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 1 {
				task := args[0]
				acme := args[1]
				setAcmeUserTaskProfile(task, acme)
			} else {
				cmd.Help()
			}
		},
	}
	return c
}

// task set
func cmdTaskSet() *cobra.Command {
	var (
		domain                  string
		acmeuser                string
		enabled                 bool
		dns_profile             string
		dns_ttl                 int64
		dns_propagation_timeout int64
		dns_polling_interval    int64
		url_check_enable        bool
		url_check_target        string
		url_check_interval      int64
	)

	c := &cobra.Command{
		Use:   "set [TASK]",
		Short: "Set config values in a task profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) == 1 {
				task := args[0]
				// set domain
				if cmd.Flags().Lookup("domain").Changed {
					ok := setTaskProfile(task, "domain", domain)
					if !ok {
						ark.Error().Msg("Set domain failed")
					}
				}

				// set acme user
				if cmd.Flags().Lookup("user").Changed {
					ok := setTaskProfile(task, "acme_user", acmeuser)
					if !ok {
						ark.Error().Msg("Set acme user failed")
					}
				}

				// set task enable
				if cmd.Flags().Lookup("enable").Changed {
					ok := setTaskProfile(task, "enable", "true")
					if !ok {
						ark.Error().Msg("Enable task failed")
					}
				}

				// set task disable
				if cmd.Flags().Lookup("disable").Changed {
					ok := setTaskProfile(task, "enable", "false")
					if !ok {
						ark.Error().Msg("Enable task failed")
					}
				}

				// set dns profile
				if cmd.Flags().Lookup("profile").Changed {
					ok := setTaskProfile(task, "dns_profile", dns_profile)
					if !ok {
						ark.Error().Msg("Set dns profile failed")
					}
				}

				// set dns ttl
				if cmd.Flags().Lookup("ttl").Changed {
					ok := setTaskProfile(task, "dns_ttl", strconv.Itoa(int(dns_ttl)))
					if !ok {
						ark.Error().Msg("Set dns ttl failed")
					}
				}

				// set dns propagation timeout
				if cmd.Flags().Lookup("propagation").Changed {
					ok := setTaskProfile(task, "dns_propagation_timeout", strconv.Itoa(int(dns_propagation_timeout)))
					if !ok {
						ark.Error().Msg("Set dns propagation timeout failed")
					}
				}

				// set dns polling interval
				if cmd.Flags().Lookup("propagation").Changed {
					ok := setTaskProfile(task, "dns_polling_interval", strconv.Itoa(int(dns_polling_interval)))
					if !ok {
						ark.Error().Msg("Set dns polling interval failed")
					}
				}

				// set url check enable
				if cmd.Flags().Lookup("url_check_enable").Changed {
					ok := setTaskProfile(task, "url_check_enable", "true")
					if !ok {
						ark.Error().Msg("Enable url check failed")
					}
				}

				// set url check disable
				if cmd.Flags().Lookup("url_check_disable").Changed {
					ok := setTaskProfile(task, "url_check_enable", "false")
					if !ok {
						ark.Error().Msg("Enable url check failed")
					}
				}

				// set url check target
				if cmd.Flags().Lookup("url_check_target").Changed {
					ok := setTaskProfile(task, "url_check_target", url_check_target)
					if !ok {
						ark.Error().Msg("Set url check target failed")
					}
				}

				// set url check interval
				if cmd.Flags().Lookup("url_check_interval").Changed {
					ok := setTaskProfile(task, "url_check_interval", strconv.Itoa(int(url_check_interval)))
					if !ok {
						ark.Error().Msg("Set url check interval failed")
					}
				}

			} else {
				cmd.Help()
			}
		},
	}

	c.Flags().StringVarP(&domain, "domain", "d", "", "set domains, separated by commas")
	c.Flags().StringVarP(&acmeuser, "user", "u", certark.DefaultTaskProfile.AcmeUser, "set acme user")
	c.Flags().BoolVar(&enabled, "enable", true, "enable task")
	c.Flags().BoolVar(&enabled, "disable", false, "disable task")

	c.Flags().StringVarP(&dns_profile, "profile", "p", certark.DefaultTaskProfile.DNSProfile, "set dns profile")
	c.Flags().Int64VarP(&dns_ttl, "ttl", "t", certark.DefaultTaskProfile.DNSTTL, "set dns record ttl")
	c.Flags().Int64Var(&dns_propagation_timeout, "propagation", certark.DefaultTaskProfile.DNSPropagationTimeout, "set propagation timeout in seconds")
	c.Flags().Int64Var(&dns_polling_interval, "interval", certark.DefaultTaskProfile.DNSPollingInterval, "set polling interval in seconds")

	c.Flags().BoolVar(&url_check_enable, "url_check_enable", true, "enable url check")
	c.Flags().BoolVar(&url_check_enable, "url_check_disable", false, "disable url check")
	c.Flags().StringVar(&url_check_target, "url_check_target", certark.DefaultTaskProfile.UrlCheckTarget, "set url check target")
	c.Flags().Int64Var(&url_check_interval, "url_check_interval", certark.DefaultTaskProfile.UrlCheckInterval, "set url check interval in days")
	return c
}

// task run command
func cmdTaskRun() *cobra.Command {
	return &cobra.Command{
		Use:   "run [TASK]",
		Short: "Run a task",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Error().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) > 0 {
				task := args[0]
				runTask(task)
			} else {
				cmd.Help()
			}
		},
	}
}

// list task profiles
func listTasks() {
	tasks, err := certark.ListTasks()
	if err != nil {
		ark.Error().Err(err).Msg("Failed to list tasks")
		return
	}
	for _, v := range tasks {
		fmt.Println(v)
	}
}

// show task profile
func showTaskProfile(task string) {
	profile, err := certark.GetTaskJsonPretty(task)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to show task")
		return
	}
	fmt.Println(profile)
}

// add task profile
func addTaskProfile(task string) {
	err := certark.AddTask(task)
	if err != nil {
		ark.Error().Msg("Failed to create task " + task)
	} else {
		ark.Info().Msg("Task " + task + " added")
	}
}

// set task profile
func setTaskProfile(task, key, value string) bool {
	err := certark.SetTaskProfile(task, key, value)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to set task profile")
		return false
	}
	return true
}

// Append domains in a task profile
func appendDomainTaskProfile(task string, domains []string) {
	err := certark.AppendDomainTaskProfile(task, domains)
	if err != nil {
		ark.Error().Msg("Failed to change task " + task)
	} else {
		ark.Info().Msg("Task " + task + " changed")
	}
}

// Remove domains in a task profile
func subtractDomainTaskProfile(task string, domain string) {
	err := certark.SubtractDomainTaskProfile(task, domain)
	if err != nil {
		ark.Error().Msg("Failed to change task " + task)
	} else {
		ark.Info().Msg("Task " + task + " changed")
	}
}

// set acme user in a task profile
func setAcmeUserTaskProfile(task string, acme string) {
	err := certark.SetAcmeUserTaskProfile(task, acme)
	if err != nil {
		ark.Error().Msg("Failed to change task " + task)
	} else {
		ark.Info().Msg("Task " + task + " changed")
	}
}

// run task
func runTask(task string) {
	if !checkTaskProfileExists(task) {
		err := errors.New("task does not existed")
		ark.Error().Err(err).Msg("Failed to run task")
		return
	}

	// read profile
	profileContent, err := os.ReadFile(certark.TaskConfigDir + "/" + task)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to run task")
		return
	}
	ark.Debug().Str("content", string(profileContent)).Msg("Read task profile")

	profile := string(profileContent)
	acmeUser := gjson.Get(profile, "acme_user").String()

	// check if acme user exists
	if !certark.CheckAcmeUserExists(acmeUser) {
		err := errors.New("acme user does not existed")
		ark.Error().Err(err).Str("task", task).Msg("Failed to found acme user in task profile")
		return
	}

	// read acme user profile
	au, err := certark.GetAcmeUser(acmeUser)
	if err != nil {
		ark.Error().Err(err).Str("task", task).Msg("Read acme user failed")
		return
	}
	config := lego.NewConfig(&acme.AcmeUser{
		Email: au.Email,
		Key:   acme.PrivateKeyDecode(au.PrivateKey),
	})
	config.CADirURL = lego.LEDirectoryStaging
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		panic(err)
	}

	// println(client)

	//TODO -
	// new provider
	drivers.ImportDrivers()
	provider_name := gjson.Get(profile, "dns_provider").String()
	driverCons, err := acme.GetDriver(provider_name)
	if err != nil {
		ark.Error().Err(err).Msg("Init dns driver failed")
		return
	}
	driver := driverCons()

	res, err := driver.RequestCertificate(client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
