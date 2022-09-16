# jobmgr
A program to queue, manage, and run jobs that are submitted to it. For example, you might have a computer with a GPU that you want to send work to. jobmgr will queue up the work, run each job one-by-one, and make the results available over an HTTP API. Tested to work on Windows and Linux.

## Overview
All work must be submitted as a *job*. Each job has a name and one or more parameters. The name will be used to look up a *jobspec*, which describes how to run the job.

When a client submits a job, jobmgr will follow the instructions in your jobspec. The output of the job (the *results*) will be stored and made available to the client.

Sometimes, your job might generate something other than text; for example, an image, video, or sound file. These are referred to as *artifacts*. Your job can write these output files in an `artifacts` folder, which is created in a new temporary directory each time the job is run. When the job is finished, jobmgr will automatically upload all its artifacts to an Amazon S3 bucket, and provide links to the client.

## Limitations/future things to do
* The manager and worker roles are combined
	* Ideally, there would be a separation between the two (so there would be a manager that runs the HTTP API and queues the jobs, and those would be separate processes on separate computers)
	* That way, you could have multiple workers processing jobs
* Jobs are stored forever
	* This is good and bad, I guess it some point all the historical data will accumulate?
* Clients must poll for job completion
	* There's no way for a client to be nontified of job status changes - they just have to poll
	* Could add WebSockets support or something like that
* No authentication - just relies on network firewalling for security
* Add support for S3-compatible alternatives for storage

## Setup
### Requirements
* Go 1.19+ (older versions might work but haven't been tested)
* MySQL database
* [roamer](https://github.com/thatoddmailbox/roamer/wiki/Installation)
* Amazon S3 bucket

### Installation and configuration
> **Warning**
> jobmgr currently performs no authentication. You should use a firewall or some other mechanism to ensure the HTTP server is not exposed to the Internet.
>
> If you need to submit jobs from another machine, try SSH port tunneling, a VPN, or [Tailscale](https://tailscale.com/).

From this folder (the one with the README), run `go build`. That should generate a `jobmgr` executable. 

You'll need to decide on a working directory for jobmgr. This will be where your configuration and jobspec files will be. On Linux, this could be `/opt/jobmgr`. On Windows, this could be `C:\jobmgr`. Wherever it is, copy the executable you just built to that folder. Then, create a `config.toml` file with the following contents:

```toml
[Database]
Host = "localhost"
Username = "jobmgr"
Password = "password123"
Database = "jobmgr"

[AWS]
AccessKeyID = "AKsomething"
SecretAccessKey = "something something"
Region = "us-east-2"
ArtifactsBucket = "bucket-name-here"
```

You will need to update the `[Database]` section with your MySQL connection information, and the `[AWS]` section with your AWS credentials and bucket name.

> **Note**
> You should create a new user in AWS IAM with access keys, and delegate bucket access to that user. For more information, see the [AWS IAM documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/getting-started.html).

Once you've done that, use a terminal to start your jobmgr executable from its working directory. That's it! You can now add a jobspec and start a job.

## Jobspecs
A jobspec describes how to run a job. It's a TOML file with various configuration options, with `Command` being the only required one.

Here is a sample jobspec, along with descriptions for each parameter:
```toml
Command = "./somecoolscript.sh"
Arguments = ["--some-useful-flag"]
WorkingDirectory = "/opt/something/"
PreserveEnvVars = ["PATH"]
Timeout = "1m30s"

[[Parameter]]
Name = "message"
Type = "string"

[[Parameter]]
Name = "other-option"
Type = "string"
```

| Name | Description |
| ---- | ----------- |
| Command | The executable to run. This must ONLY be the executable name, EXCLUDING any arguments! |
| Arguments | A list of arguments to pass to the executable when it's run. |
| WorkingDirectory | The working directory to use when launching the executable. If omitted, defaults to the job's temporary directory. |
| PreserveEnvVars | A list of environment variables to preserve when running the job. By default, this is empty and so no environment variables are set, apart from those described in the "Environment" section below. |
| Timeout | The maximum runtime of a job, in [Go duration string format](https://pkg.go.dev/time#ParseDuration). If omitted, defaults to 10 seconds. |

In addition to these options, you may specify one or more parameters for the job, as shown above. Each parameter has a Name and a Type. Valid parameter types are "string" and "int". The values of these parameters are set by whoever submits the job, and are made available as environment variables. See the "Environment" section below for more details.

### Environment
When your job is started, jobmgr creates a new temporary directory. This will be the working directory for your job (unless you override it in the jobspec). The contents of this directory are cleaned up when the job is finished, so you can write any temporary files or output there.

Inside the temporary directory, jobmgr creates another directory called `artifacts`. Anything your job writes into `artifacts` will be uploaded to Amazon S3 and made available to the client who submitted the job. This way, you can produce output files (images, music, videos, etc.) in addition to text output.

For convenience, the full path to the `artifacts` directory is also made available in the environment variable `JOBMGR_ARTIFACTS_DIR`.

In addition, if you specified any parameters for the job, those parameter's values will be made available as environment variables. The parameter's name will be normalized. For example, a parameter named `message` would be available in the environment variable `JOBMGR_PARAMETER_MESSAGE`.

#### Preserving environment variables
Apart from the `JOBMGR` environment variables described above, by default no other environment variables are set. However, you sometimes want a job to access an environment variable; for example, you might need to access something in the `PATH`. You can do this by setting the `PreserveEnvVars` jobspec option to a list of environment variables to preserve.

> **Note**
> Missing environment variables can be a bigger problem on Windows systems. If your job behaves weirdly when run from jobmgr but works when run directly, there probably is a missing environment variable. You can use the `set` (Command Prompt) or `dir env:` (PowerShell) commands to list all environment variables. Then, try running your job with `PreserveEnvVars` set to the entire list, and slowly remove environment variables until your job fails again.

## HTTP API
TODO