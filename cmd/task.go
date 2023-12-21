package cmd

import (
	"fmt"
	"strconv"

	"github.com/josark2005/certark/ark"
	"github.com/josark2005/certark/certark"
	"github.com/spf13/cobra"
)

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
		Use:     "show [TASK]",
		Short:   "Show / inspec a task profile",
		Aliases: []string{"inspec"},
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
		domain      string
		acmeuser    string
		enabled     bool
		dns_profile string
		// dns_ttl                 int64
		// dns_propagation_timeout int64
		// dns_polling_interval    int64
		url_check_enable   bool
		url_check_target   string
		url_check_interval int64
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
					setTaskProfile(task, "domain", domain)
				}

				// set acme user
				if cmd.Flags().Lookup("user").Changed {
					setTaskProfile(task, "acme_user", acmeuser)
				}

				// set task enable
				if cmd.Flags().Lookup("enable").Changed {
					setTaskProfile(task, "enable", "true")
				}

				// set task disable
				if cmd.Flags().Lookup("disable").Changed {
					setTaskProfile(task, "enable", "false")
				}

				// set dns profile
				if cmd.Flags().Lookup("profile").Changed {
					setTaskProfile(task, "dns_profile", dns_profile)
				}

				// set url check enable
				if cmd.Flags().Lookup("url_check_enable").Changed {
					setTaskProfile(task, "url_check_enable", "true")
				}

				// set url check disable
				if cmd.Flags().Lookup("url_check_disable").Changed {
					setTaskProfile(task, "url_check_enable", "false")
				}

				// set url check target
				if cmd.Flags().Lookup("url_check_target").Changed {
					setTaskProfile(task, "url_check_target", url_check_target)
				}

				// set url check interval
				if cmd.Flags().Lookup("url_check_interval").Changed {
					setTaskProfile(task, "url_check_interval", strconv.Itoa(int(url_check_interval)))
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

	c.Flags().StringVarP(&dns_profile, "profile", "p", certark.DefaultTaskProfile.DnsProfile, "set dns profile")

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
				err := runTask(task)
				if err != nil {
					ark.Error().Err(err).Msg("Failed to run task")
				} else {
					ark.Info().Str("task", task).Msg("Task executed")
				}
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
func runTask(name string) error {
	return certark.RunTaskIndependently(name)
}
