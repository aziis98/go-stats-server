# Go Stats Server

This Go project implements a simple TCP server that responds to custom commands over an internal network. The server performs various system-related tasks and provides information such as CPU usage, memory status, network statistics, storage details, and system uptime.

## Commands

The protocol supports the following commands: 

- `cpu` &mdash; returns the CPU usage percentage
- `memory` &mdash; returns the memory usage in MB
- `network` &mdash; returns the network usage in bytes
- `storage` &mdash; returns the storage usage in GB
- `uptime` &mdash; returns the system uptime in seconds

[TODO: check if the command docs are correct]

## Usage

1. **Build the Server:**

    ```bash shell
    $ CGO_ENABLED=0 GOARCH=<arch> go build -a -ldflags '-s -w' -o ./out/stats-server main.go
    ```

2. **Run the Server:**

    ```bash shell
    $ ./out/stats-server serve
    ```

3. **Connect to the Server:** Use a TCP client to connect to the server on port 12345. You can send commands like "cpu", "memory", "network", "storage" or
   "uptime".

    Example using `nc`:

    ```bash shell
    $ time ( echo "cpu" | nc <hostname> 12345 )
    ```

    or directly using golang:

    ```go
    conn, _ := net.Dial("tcp", "<hostname>:12345")
    conn.Write([]byte("cpu\n"))

    data, _ := io.RealAll(conn)
    log.Printf("data: %s", data)
    ```

## GitHub Actions Workflow

The included GitHub Actions workflow automates the build and release process. On each push to the main branch, the workflow builds the Go program, creates a GitHub release, and uploads the compiled binary as an artifact.

## Self-Installing Binary

Assuming you have a cluster of machines, you can use the following commands to download, install and setup systemd services for the latest version of the binary on all nodes. You must have root access to all nodes.

```bash
$ parallel --nonall --slf nodes.txt 'mkdir -p cluster'
$ parallel --nonall --slf nodes.txt 'wget -qO- https://github.com/aziis98/go-stats-server/releases/latest/download/stats-server > ./cluster/stats-server'
$ parallel --nonall --slf nodes.txt 'chmod -v +x ./cluster/stats-server'

# setup systemd services on all nodes
$ parallel --nonall --slf nodes.txt './cluster/stats-server setup'

# setup systemd service on a single node
$ ssh root@<node>
$ cd cluster
$ ./stats-server setup

# trying out the new binary using netcat
$ time ( echo '<command>' | nc <node> 12345 )

# stop and disable the systemd service on all nodes
$ parallel --nonall --slf nodes.txt 'systemctl disable --now stats-server.service'
```

## Credits

- The shell code for the commands was originally made by [@BachoSeven](https://github.com/bachoseven)

    <https://git.phc.dm.unipi.it/phc/cluster-dashboard/src/branch/main/backend/scripts>

- Most of the code was initially made with the help of ChatGPT in this conversation

    <https://chat.openai.com/share/40247aa0-76e9-4d9a-b0fa-356a5f51f208>