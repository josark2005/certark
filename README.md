# CertArk

CertArk is a certificate requestor based on [lego](https://github.com/go-acme/lego).

## Usage

### Initialization

An initialization is required before running CertArk at the first time.

```bash
sudo certark init
```

Initialization process will create a `.lock` file in configuration directory, the init process WON'T execute if the lock file exists. The `--force` flag can ignore the limitation.

```bash
sudo certark init --force
```

Or remove the `.lock` file. A `--yes-i-really-mean-it` flag is required to comfirm.

```bash
sudo certark init unlock --yes-i-really-mean-it
```

> Notice: Before you execute the initialization process or remove the lock file, backup your configurations.


### Manage ACME account (user)

ACME account is required by the Let's Encrypt orgnization.

```bash
sudo certark acme add [ACME_ACCOUNT_EMAIL]
```

Then, register the account.

```bash
sudo certark acme reg [ACME_ACCOUNT_EMAIL]
```

> Now, a task is available to create and run.


### Manage task

Task is the essential element of configurations. CertArk will not run without tasks.

```bash
sudo certark task add [TASK_NAME]
```

Then add domain(s) to the task.

```bash
sudo certark task append [TASK_NAME] [DOMAIN]
```

Or remove a domain from a task.

```bash
sudo certark task subtract [TASK_NAME] [DOMAIN]
```

Or set a single domain in a task.

```bash
sudo certark task set [TASK_NAME] -d [DOMAIN]
```

Then, allocate a acme account to the task.

```bash
sudo certark task set [TASK_NAME] -u [ACME_ACCOUNT_EMAIL]
```

> Other factors required can be set with `task set` command above. Set the `-h` flag to show all available flags/settings.


## Run in all-in-one (standalone) mode

Standalone mode means a single-node mode which including the request->download->deploy process.

```bash
sudo certark server -a
# or
sudo certark server --standalone
```

A deploy command should be specified in task profile. See details in [Automatic](#automatic-deploy)


## Automatic Deploy

Automatic Deploy is not supported yet.

## DNS Provider Authentication

Different dns providers require distinctive API authentication element(s). DNS providers authentication options listed blow are supported.

### Cloudflare

- API_Email (dns_authmail)
- API_Token (dns_authtoken)
- Permission required: 
  - `Zone / Zone / Read`
  - `Zone / DNS / Edit`


## For Developers

### Before Comimits

#### git_conifig.sh

It's very recommended to execute `scripts/dev/git_conifig.sh` file to set up a development standard environment.

#### clean.sh

If you want to clean configurations, logs or any other files that certark created, execute the `scripts/dev/clean.sh`. 

**WARNING: ALL YOUR CONFIGURATION FILES WILL BE REMOVED AFTER EXECUTE THIS SHELL FILE**