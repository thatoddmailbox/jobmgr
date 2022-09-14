# jobmgr
A program to queue, manage, and run jobs that are submitted to it. For example, you might have a computer with a GPU that you want to send work to. jobmgr will queue up the work, run each job one-by-one, and make the results available over an HTTP API.

## Overview
All work must be submitted as a *job*. Each job has a name and one or more parameters. The name will be used to look up a *jobspec*, which describes how to run the job.

When a client submits a job, jobmgr will follow the instructions in your jobspec. The output of the job will be stored and made available to the client.

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
From this folder (the one with the README), run `go build`. That should generate a `jobmgr` executable. 

You'll need to decide on a working directory for jobmgr. This will be where your configuration and jobspec files will be. On Linux, this could be `/opt/jobmgr`; on Windows, this could be `C:\jobmgr`. Wherever it is, copy the executable you just build to that folder. Then, create a `config.toml` file with the following contents:

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

You will need to update the [Database] section with your MySQL connection information, and the [AWS} section with your AWS credentials and bucket name.

> **Info:**
> You should create a new user in AWS IAM with access keys, and delegate bucket access to that user. For more information, see the [AWS IAM documentation](https://docs.aws.amazon.com/IAM/latest/UserGuide/getting-started.html).

Once you've done that, use a terminal to start your jobmgr executable from its working directory. THat's it! You can now add a jobspec and start a job.