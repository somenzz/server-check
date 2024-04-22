# Disk space monitoring tool written in golang


When the CPU, disk, and memory usage exceeds the preset value, an enterprise WeChat alarm notification is sent.

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
*/5 * * * * /path/to/your/server-check
```
If you use Windows, use Scheduled Task Execution.

