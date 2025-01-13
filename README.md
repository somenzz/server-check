# Disk space monitoring tool written in golang


When the CPU, disk, and memory usage exceeds the preset value, an enterprise WeChat alarm notification is sent.

## Update log:

Added feature: 

- health check for http service, you can add http check information in config.yaml.
- When cpu exceeds the threshold, the information of the top five processes occupying the CPU is output simultaneously. 

## Configuration file

```yaml
ewechat:
  corp_id: "your enterprise wechat corp_id"
  corp_secret: "your enterprise wechat corp_secret"
  agent_id: your enterprise wechat agent_id
  receivers: "your enterprise wechat receivers, for more receiver: receiver1|receiver2"
cpu_usage_rate: 90.0  #When the CPU usage exceeds 90, an enterprise WeChat notification will be sent. 
mem_usage_rate: 90.0  #When the Mem usage exceeds 90, an enterprise WeChat notification will be sent.
disk_usage_rate: 90.0 #When the Disk usage exceeds 90, an enterprise WeChat notification will be sent.
check_url:
  - url: "https://xxxx/api/health"
    method: get  # this is default
    expect_status_code: 200 # this is default
    expect_body: "ok"  # if expect_body is in the resp.Body, it returns true.
  - url: "https://xxxxx"
    method: post
    expect_status_code: 403
    expect_body: "As long as the returned string contains expect_body, it's a success"
```

## Instructions

### First populate the configuration file

```sh
cp config-example.yaml config.yaml
# populate the configuration file config.yaml
```
### Compile and run

```sh
go build && ./server-check
```
### Use crontab for regular monitoring

```sh
*/5 * * * * cd /path/to/your && /path/to/your/server-check
```
If you use Windows, use Scheduled Task Execution.

