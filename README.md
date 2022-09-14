# jobmgr
A program to queue, manage, and run jobs that are submitted to it. For example, you might have a computer with a GPU that you want to send work to. jobmgr will queue up the work, run each job one-by-one, and make the results available over an HTTP API.

## Overview
All work must be submitted as a *job*. Each job has a name and one or more parameters. The name will be used to look up a *jobspec*, which describes how to run the job.

When a client submits a job, jobmgr will follow the instructions in your jobspec. The output of the job will be stored and made available to the client.

Sometimes, your job might generate something other than text; for example, an image, video, or sound file. These are referred to as *artifacts*. Your job can write these output files in an `artifacts` folder, which is created in a new temporary directory each time the job is run. When the job is finished, jobmgr will automatically upload all its artifacts to an Amazon S3 bucket, and provide links to the client.

## Limitations
* The manager and worker roles are combined
	* Ideally, there would be a separation between the two (so there would be a manager that runs the HTTP API and queues the jobs, and those would be separate processes on separate computers)
	* That way, you could have multiple workers processing jobs
* Jobs are stored forever
	* This is good and bad, I guess it some point all the historical data will accumulate?
* Clients must poll for job completion
	* There's no way for a client to be nontified of job status changes - they just have to poll
	* Could add WebSockets support or something like that
* No authentication - just relies on network firewalling for security
