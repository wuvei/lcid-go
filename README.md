# LCid-Go
## Introduction

This naive toy is based on [bunnyxt/lcid](https://github.com/bunnyxt/lcid), and implemented in Golang for practice. They are same in program logic and static files. Thanks bunnyxt!

I deployed this toy on [lcid.yczheng.top](https://lcid.yczheng.top). 

- For Leetcode problems on the global site, type `lcid.yczheng.top/<problem_id>` in your browser.
- For problems on the Chinese site, type `lcid.yczheng.top/cn/<problem_id>`.
- Or you can enter [lcid.yczheng.top](https://lcid.yczheng.top) to have fun.

## Development

### Setup

Only go is needed.

```bash
# Mac
brew install go
# Debian, Ubuntu
apt install go
# CentOS, Fedora
dnf install go
```

### Fetch All Problems

To fetch all problems, execute `go run ./problem.go`. Then, you will see the logs below.

```bash
csrftoken = rqzxXv6RW9wwrkHoiyegudOgl92V9yPdSmEWfBjfIUnJQtfbbGZEnHDL3DXxwbKL
Found 2137 problems in total.
Now try fetch all 2137 LeetCode problems...
All 2137 problems info saved into problems_all.json file.
```

After few seconds, all problems info have been saved into `problems_all.json` file in json format.

### Start Backend Service

To start backend service, execute `go run ./main.go <port>` (e.g. `go run ./main.go 9191`). Then, the backend server will start at `localhost:<port>`.  If no port is added, `9191` is the default.

Now you can try the following endpoints. All of them should work correctly. 

- `localhost:<port>/`
- `localhost:<port>/1`
- `localhost:<port>/cn/1`
- `localhost:<port>/info/1`

### Deployment on Heroic

Please refer to [bunnyxt/lcid#deployment](https://github.com/bunnyxt/lcid#deployment).

### Deployment on Custom Server

#### Nohup

Just start the backend server and setup your nginx config.

Optional: `crontab -e`  add the task of fetching provblem list to scheduled lists, eg.

 ```
 5 2 * * * cd /root/lcid-go && /usr/local/go/bin/go run problem.go > ~/fetchlog.txt # Fetch Problem.
 12 2 * * * cd /root/lcid-go && ./lcid_restart.sh
 ```

Note: in `lcid_restart.sh`, I assume that there's only 1 go program executing. 

#### Systemd Service (Recommended)

Refer to [lcid-go.service](./lcid-go.service) 

Suppose you are user `root`, you could copy the `lcid-go.service` to `/usr/lib/systemd/system/`.

Then with `systemctl start lcid-go.serrvice`, the program is up.

`systemctl start lcid-go.serrvice` will add lcid-go to startup programs. Then cron jobs can be write as:

 ```
 5 2 * * * cd /root/lcid-go && /usr/local/go/bin/go run problem.go > ~/fetchlog.txt # Fetch Problem.
 12 2 * * * systemctl restart lcid-go
 ```