package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// check if task profile exists
func checkTaskProfileExists(taskname string) bool {
	res := certark.FileOrDirExists(taskConfigDir + "/" + taskname)
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

	rootCmd.AddCommand(taskCmd)
}

// task command
func cmdTask() *cobra.Command {
	return &cobra.Command{
		Use:   "task",
		Short: "Task configurations",
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

// task add command
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
		dns_provider            string
		dns_authuser            string
		dns_authkey             string
		dns_authtoken           string
		dns_zonetoken           string
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
					ok := setTaskProfile(task, "acme_user", domain)
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

				// set dns provider
				if cmd.Flags().Lookup("provider").Changed {
					ok := setTaskProfile(task, "dns_provider", dns_provider)
					if !ok {
						ark.Error().Msg("Set dns provider failed")
					}
				}

				// set dns auth user
				if cmd.Flags().Lookup("authuser").Changed {
					ok := setTaskProfile(task, "dns_authuser", dns_authuser)
					if !ok {
						ark.Error().Msg("Set dns auth user failed")
					}
				}

				// set dns auth key
				if cmd.Flags().Lookup("authkey").Changed {
					ok := setTaskProfile(task, "dns_authkey", dns_authkey)
					if !ok {
						ark.Error().Msg("Set dns authkey failed")
					}
				}

				// set dns auth token
				if cmd.Flags().Lookup("authtoken").Changed {
					ok := setTaskProfile(task, "dns_authtoken", dns_authtoken)
					if !ok {
						ark.Error().Msg("Set dns authtoken failed")
					}
				}

				// set dns zone token
				if cmd.Flags().Lookup("zonetoken").Changed {
					ok := setTaskProfile(task, "dns_zonetoken", dns_zonetoken)
					if !ok {
						ark.Error().Msg("Set dns zonetoken failed")
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
					ok := setTaskProfile(task, "url_check_target", dns_zonetoken)
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

	c.Flags().StringVar(&dns_provider, "provider", certark.DefaultTaskProfile.DNSProvider, "set dns provider")
	c.Flags().StringVar(&dns_authuser, "authuser", certark.DefaultTaskProfile.DNSAuthUser, "set dns auth user or email")
	c.Flags().StringVar(&dns_authkey, "authkey", certark.DefaultTaskProfile.DNSAuthKey, "set dns auth key")
	c.Flags().StringVar(&dns_authtoken, "authtoken", certark.DefaultTaskProfile.DNSAuthToken, "set dns auth token")
	c.Flags().StringVar(&dns_zonetoken, "zonetoken", certark.DefaultTaskProfile.DNSZoneToken, "set dns zone token")
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
	err := filepath.Walk(taskConfigDir, func(path string, info os.FileInfo, err error) error {
		if path == taskConfigDir {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println(path[len(taskConfigDir)+1:])
		return nil
	})
	if err != nil {
		ark.Error().Err(err).Msg("Failed to list task profiles")
		return
	}
}

// show task profile
func showTaskProfile(task string) {
	profile := taskConfigDir + "/" + task
	if !certark.FileOrDirExists(profile) || !certark.IsFile(profile) {
		err := errors.New("task " + task + " does not exist")
		ark.Error().Err(err).Msg("Failed to show acme user")
		return
	}

	// read file
	profileContent, err := os.ReadFile(profile)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to show task profile")
		return
	}

	var jsonBuff bytes.Buffer
	if err = json.Indent(&jsonBuff, profileContent, "", ""); err != nil {
		ark.Error().Err(err).Msg("Failed to show task profile")
		return
	}

	fmt.Println(jsonBuff.String())
}

// add task profile
func addTaskProfile(task string) {
	if checkTaskProfileExists(task) {
		err := errors.New("task existed")
		ark.Error().Err(err).Msg("Failed to create user profile")
		return
	}

	// create profile
	fp, err := os.OpenFile(taskConfigDir+"/"+task, os.O_CREATE|os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to create task profile")
		return
	}
	defer fp.Close()

	profile := certark.DefaultTaskProfile
	profile.TaskName = task
	profileJson, _ := json.Marshal(profile)

	// write profile to file
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Msg("Failed to add task " + task)
	} else {
		ark.Info().Msg("Task " + task + " added")
	}
}

// set task profile
func setTaskProfile(task, key, value string) bool {
	if !checkTaskProfileExists(task) {
		err := errors.New("task does not existed")
		ark.Error().Err(err).Msg("Failed to append domains to task profile")
		return false
	}

	supportedKey := []string{
		"domain",
		"acme_user",
		"enabled",
		"dns_provider",
		"dns_authuser",
		"dns_authkey",
		"dns_authtoken",
		"dns_zonetoken",
		"dns_ttl",
		"dns_propagation_timeout",
		"dns_polling_interval",
		"url_check_enable",
		"url_check_target",
		"url_check_interval",
	}

	supportFlag := false
	for _, sk := range supportedKey {
		if key == sk {
			supportFlag = true
			break
		}
	}
	if !supportFlag {
		err := errors.New("not supported configuration key")
		ark.Error().Str("key", key).Err(err).Msg("Failed to set task profile")
		return false
	}

	ark.Info().Str("key", key).Str("value", value).Msg("Setting task profile")

	// read profile
	profileContent, err := os.ReadFile(taskConfigDir + "/" + task)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read task profile")
	}
	ark.Debug().Str("content", string(profileContent)).Msg("Read task profile")
	profile := certark.TaskProfile{}
	err = json.Unmarshal(profileContent, &profile)
	if err != nil {
		ark.Error().Err(err).Str("task", task).Msg("Failed to parse task profile")
	}

	switch key {
	case "domain":
		profile.Domain = []string{value}
	case "acme_user":
		profile.AcmeUser = value
	case "enable":
		if value == "true" {
			profile.Enabled = true
		} else {
			profile.Enabled = false
		}
	case "dns_provider":
		profile.DNSProvider = value
	case "dns_authuser":
		profile.AcmeUser = value
	case "dns_authkey":
		profile.DNSAuthKey = value
	case "dns_authtoken":
		profile.DNSAuthToken = value
	case "dns_zonetoken":
		profile.DNSZoneToken = value
	case "dns_ttl":
		v, e := strconv.Atoi(value)
		if e != nil {
			ark.Error().Err(e).Msg("Set dns ttl failed")
		}
		profile.DNSTTL = int64(v)
	case "dns_propagation_timeout":
		v, e := strconv.Atoi(value)
		if e != nil {
			ark.Error().Err(e).Msg("Set dns propagation timeout failed")
		}
		profile.DNSPropagationTimeout = int64(v)
	case "dns_polling_interval":
		v, e := strconv.Atoi(value)
		if e != nil {
			ark.Error().Err(e).Msg("Set dns polling interval failed")
		}
		profile.DNSPollingInterval = int64(v)
	case "url_check_enable":
		if value == "true" {
			profile.UrlCheckEnable = true
		} else {
			profile.UrlCheckEnable = false
		}
	case "url_check_target":
		profile.UrlCheckTarget = value
	case "url_check_interval":
		v, e := strconv.Atoi(value)
		if e != nil {
			ark.Error().Err(e).Msg("Set dns propagation timeout failed")
		}
		profile.UrlCheckInterval = int64(v)
	default:
		ark.Error().Msg("Failed to found a valid configuration key")
	}

	// write profile to file
	profileJson, _ := json.Marshal(profile)
	fp, err := os.OpenFile(taskConfigDir+"/"+task, os.O_WRONLY|os.O_TRUNC, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to open task profile")
		return false
	}
	defer fp.Close()
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Msg("Failed to change task " + task)
	} else {
		ark.Info().Msg("Task " + task + " changed")
	}

	return true
}

// Append domains in a task profile
func appendDomainTaskProfile(task string, domains []string) {
	if !checkTaskProfileExists(task) {
		err := errors.New("task does not existed")
		ark.Error().Err(err).Msg("Failed to append domains to task profile")
		return
	}

	// read profile
	profileContent, err := os.ReadFile(taskConfigDir + "/" + task)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read task profile")
		return
	}
	ark.Debug().Str("content", string(profileContent)).Msg("Read task profile")

	origDoamin := []string{}
	newDoamin := []string{}
	for _, v := range gjson.Get(string(profileContent), "domain").Array() {
		if v.String() != "" {
			origDoamin = append(origDoamin, v.String())
		}
	}

	// filter old domains
	for _, origD := range origDoamin {
		if len(newDoamin) == 0 {
			newDoamin = append(newDoamin, origD)
		} else {
			dflag := false
			for _, newD := range newDoamin {
				if origD == newD {
					dflag = true
					ark.Debug().Str("domain", origD).Msg("Skip duplicated domain")
					continue
				}
			}
			if !dflag {
				newDoamin = append(newDoamin, origD)
			}
		}
	}

	// add new domains
	for _, newD := range domains {
		dflag := false
		for _, domain := range newDoamin {
			if domain == newD {
				dflag = true
				ark.Warn().Str("domain", newD).Msg("Dulipcated domain")
				continue
			}
		}
		if dflag {
			continue
		} else {
			ark.Debug().Str("domain", newD).Msg("New domain found")
			newDoamin = append(newDoamin, newD)
		}
	}

	profile := certark.TaskProfile{}
	err = json.Unmarshal([]byte(profileContent), &profile)
	if err != nil {
		ark.Error().Err(err).Str("task", task).Msg("Failed to parse task profile")
	}
	profile.Domain = newDoamin
	profileJson, _ := json.Marshal(profile)

	// write profile to file
	fp, err := os.OpenFile(taskConfigDir+"/"+task, os.O_WRONLY|os.O_TRUNC, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to open task profile")
		return
	}
	defer fp.Close()
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Msg("Failed to change task " + task)
	} else {
		ark.Info().Msg("Task " + task + " changed")
	}
}

// Remove domains in a task profile
func subtractDomainTaskProfile(task string, domain string) {
	if !checkTaskProfileExists(task) {
		err := errors.New("task does not existed")
		ark.Error().Err(err).Msg("Failed to append domains to task profile")
		return
	}

	// read profile
	profileContent, err := os.ReadFile(taskConfigDir + "/" + task)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read task profile")
		return
	}
	ark.Debug().Str("content", string(profileContent)).Msg("Read task profile")

	origDoamin := []string{}
	newDoamin := []string{}
	for _, v := range gjson.Get(string(profileContent), "domain").Array() {
		if v.String() != "" {
			origDoamin = append(origDoamin, v.String())
		}
	}

	// filter domains
	for _, origD := range origDoamin {
		if origD == domain {
			continue
		} else {
			newDoamin = append(newDoamin, origD)
		}
	}

	profile := certark.TaskProfile{
		TaskName: gjson.Get(string(profileContent), "task_name").String(),
		Domain:   newDoamin,
		AcmeUser: gjson.Get(string(profileContent), "acme_user").String(),
		Enabled:  gjson.Get(string(profileContent), "enabled").Bool(),
	}
	profileJson, _ := json.Marshal(profile)

	// write profile to file
	fp, err := os.OpenFile(taskConfigDir+"/"+task, os.O_WRONLY|os.O_TRUNC, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to open task profile")
		return
	}
	defer fp.Close()
	_, err = fp.WriteString(string(profileJson))
	if err != nil {
		ark.Error().Msg("Failed to change task " + task)
	} else {
		ark.Info().Msg("Task " + task + " changed")
	}
}

// set acme user in a task profile
func setAcmeUserTaskProfile(task string, acme string) {
	if !checkTaskProfileExists(task) {
		err := errors.New("task does not existed")
		ark.Error().Err(err).Msg("Failed to set acme user to task profile")
		return
	}

	// check if acme user exists
	if !checkUserExists(acme) {
		err := errors.New("acme user does not existed")
		ark.Error().Err(err).Msg("Failed to set acme user to task profile")
		return
	}

	// read profile
	profileContent, err := os.ReadFile(taskConfigDir + "/" + task)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to read task profile")
		return
	}
	ark.Debug().Str("content", string(profileContent)).Msg("Read task profile")

	origDoamin := []string{}
	for _, v := range gjson.Get(string(profileContent), "domain").Array() {
		if v.String() != "" {
			origDoamin = append(origDoamin, v.String())
		}
	}

	profile := certark.TaskProfile{
		TaskName: gjson.Get(string(profileContent), "task_name").String(),
		Domain:   origDoamin,
		AcmeUser: acme,
		Enabled:  gjson.Get(string(profileContent), "enabled").Bool(),
	}
	profileJson, _ := json.Marshal(profile)

	// write profile to file
	fp, err := os.OpenFile(taskConfigDir+"/"+task, os.O_WRONLY|os.O_TRUNC, os.ModeExclusive)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to open task profile")
		return
	}
	defer fp.Close()
	_, err = fp.WriteString(string(profileJson))
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
	profileContent, err := os.ReadFile(taskConfigDir + "/" + task)
	if err != nil {
		ark.Error().Err(err).Msg("Failed to run task")
		return
	}
	ark.Debug().Str("content", string(profileContent)).Msg("Read task profile")

	profile := string(profileContent)
	acmeUser := gjson.Get(profile, "acme_user").String()

	// check if acme user exists
	if !checkUserExists(acmeUser) {
		err := errors.New("acme user does not existed")
		ark.Error().Err(err).Str("task", task).Msg("Failed to found acme user in task profile")
		return
	}

	// read acme user profile
	au, err := GetAcmeUser(acmeUser)
	if err != nil {
		ark.Error().Err(err).Str("task", task).Msg("Read acme user failed")
		return
	}
	config := lego.NewConfig(&au)
	config.CADirURL = lego.LEDirectoryStaging
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		panic(err)
	}

	//TODO - Run task
	fmt.Println(client)

}
