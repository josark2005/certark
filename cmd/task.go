package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/jokin1999/certark/ark"
	"github.com/jokin1999/certark/certark"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type taskProfile struct {
	TaskName              string   `json:"task_name"`
	Domain                []string `json:"domain"`
	AcmeUser              string   `json:"acme_user"`
	Enabled               bool     `json:"enabled"`
	DNSProvider           string   `json:"dns_provider"`
	DNSAuthUser           string   `json:"dns_authuser"`
	DNSAuthKey            string   `json:"dns_authkey"`
	DNSAuthToken          string   `json:"dns_authtoken"`
	DNSZoneToken          string   `json:"dns_zonetoken"`
	DNSTTL                int64    `json:"dns_ttl"`                 // ttl 120 is recommanded
	DNSPropagationTimeout int64    `json:"dns_propagation_timeout"` // in millisecond, 60*1000 is recommanded
	DNSPollingInterval    int64    `json:"dns_polling_interval"`    // in millisecond, 5 *1000 is recommanded
}

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
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
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
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
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
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
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
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
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
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
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
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
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
	)

	c := &cobra.Command{
		Use:   "set [TASK]",
		Short: "Set config values in a task profile",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
			}
			if len(args) == 1 {
				task := args[0]
				// set domain
				ok := setTaskProfile(task, "doamin", domain)
				if !ok {
					ark.Fatal().Msg("Set domain failed")
				}
			} else {
				cmd.Help()
			}
		},
	}

	c.Flags().StringVarP(&domain, "domain", "d", "", "set domains, separated by commas")
	c.Flags().StringVarP(&acmeuser, "user", "u", "", "set acme user")
	c.Flags().BoolVar(&enabled, "enable", true, "enable task")
	c.Flags().BoolVar(&enabled, "disable", false, "disable task")
	c.Flags().StringVar(&dns_provider, "provider", "", "set dns provider")
	c.Flags().StringVar(&dns_authuser, "authuser", "", "set dns auth user or email")
	c.Flags().StringVar(&dns_authkey, "authkey", "", "set dns auth key")
	c.Flags().StringVar(&dns_authtoken, "authtoken", "", "set dns auth token")
	c.Flags().StringVar(&dns_zonetoken, "zonetoken", "", "set dns zone token")
	c.Flags().Int64VarP(&dns_ttl, "ttl", "t", 120, "set dns record ttl")
	c.Flags().Int64Var(&dns_propagation_timeout, "propagation", 60, "set propagation timeout in seconds")
	c.Flags().Int64Var(&dns_polling_interval, "interval", 5, "set polling interval in seconds")
	return c
}

// task run command
func cmdTaskRun() *cobra.Command {
	return &cobra.Command{
		Use:   "run [TASK]",
		Short: "Run a task",
		Run: func(cmd *cobra.Command, args []string) {
			if !CheckRunCondition() {
				ark.Fatal().Msg("Run condition check failed, try to run 'certark init' first")
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

	profile := taskProfile{
		TaskName: task,
		Domain:   []string{""},
		AcmeUser: "",
		Enabled:  true,
	}
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

	profile := taskProfile{
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

	profile := taskProfile{
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

	profile := taskProfile{
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

	fmt.Println(client)

}
