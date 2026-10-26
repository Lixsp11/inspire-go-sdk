# Training jobs

## Build and validate the request

Resolve the actual workspace, project, compute group, quota, image, framework, resources and priority from the request or authorized queries. Ask for missing values that available queries cannot establish. Placeholder IDs below are not runnable user instructions.

```json
{
  "name": "training-placeholder",
  "workspace_id": "workspace-placeholder",
  "project_id": "project-placeholder",
  "logic_compute_group_id": "group-placeholder",
  "command": "python train.py",
  "framework": "pytorch",
  "task_priority": 1,
  "enable_notification": false,
  "framework_config": [{
    "cpu": 4,
    "gpu_count": 0,
    "mem_gi": 16,
    "shm_gi": 2,
    "image": "registry.example/training:tag",
    "image_type": "SOURCE_PRIVATE",
    "instance_count": 1,
    "resource_spec_price": {
      "cpu_type": "cpu-placeholder",
      "cpu_count": 4,
      "memory_size_gib": 16,
      "logic_compute_group_id": "group-placeholder",
      "quota_id": "quota-placeholder"
    }
  }]
}
```

```shell
inspire train CreateJobConsole --body-file request.json --cli-dry-run-json
```

The dry run validates the SDK model and prints the native request without loading credentials or accessing the platform. It does not authorize submission. Resource counts are per instance; multiply by InstanceCount. Shared memory cannot exceed memory. GPU configurations require a confirmed GPU type. Priority must be explicit (1..10) and permitted by the selected space/project; the example value is not a universal default.

## Submit and observe

```shell
inspire train CreateJobConsole --body-file request.json
inspire train GetJob --JobId returned-job-id
inspire train GetJob --JobId returned-job-id --waiter expr=Result.status to=RUNNING timeout=180 interval=5
```

Submit only the validated request covered by the user's authorization. `Result.job_id` confirms creation; query the job to verify state and configuration. Queued or image-pulling is not running. Use the actual raw state spelling returned by the platform (for example RUNNING or job_running). The waiter only supports `Result.status`; timeout/interval are seconds. Failed terminal states return an error.

If the command reports MutationUncertain or `recovery.action=reconcile`, preserve the submission time, scope and safe request ID, then ListJobs/GetJob to find and verify candidates. Do not submit again merely because a request timed out or a same-name job was not immediately visible.

## Stop

```shell
inspire train StopJob --JobId confirmed-job-id
inspire train GetJob --JobId confirmed-job-id --waiter expr=Result.status to=STOPPED timeout=180 interval=5
```

Confirm the target from current state and authorization before stopping. An accepted stop response does not prove resources are released; query until the appropriate observed terminal state. A waiter timeout does not stop the job.
