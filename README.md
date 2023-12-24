# Go Stats Server

This Go project implements a simple TCP server that responds to custom commands over the network. The server performs various system-related tasks and provides
information such as CPU usage, memory status, network statistics, storage details, and system uptime.

## Usage

1. **Build the Server:**

    ```bash
    CGO_ENABLED=0 GOARCH=<arch> go build -a -ldflags '-s -w' -o ./out/stats-server main.go
    ```

2. **Run the Server:**

    ```bash
    ./out/stats-server
    ```

3. **Connect to the Server:** Use a TCP client to connect to the server on port 12345. You can send commands like "cpu," "memory," "network," "storage,"
   "uptime," and "exit."

    Example using `nc`:

    ```bash
    echo "cpu" | nc localhost 12345
    ```

    or using golang

    ```go
    import "net"

    func main() {
        conn, err := net.Dial("tcp", "localhost:12345")
        if err != nil {
            // handle error
        }
        defer conn.Close()

        conn.Write([]byte("cpu"))
    }
    ```

## GitHub Actions Workflow

The included GitHub Actions workflow automates the build and release process. On each push to the main branch, the workflow builds the Go program, creates a
GitHub release, and uploads the compiled binary as an artifact.

## Self-Installing Binary

Assuming you have a cluster of machines, you can use the following commands to download and install the latest version of the binary on all nodes. You must have root access to all nodes.

```bash
$ parallel --nonall --slf nodes.txt 'mkdir -p cluster'
$ parallel --nonall --slf nodes.txt 'wget -qO- https://github.com/aziis98/go-stats-server/releases/latest/download/stats-server > ./cluster/stats-server'
$ parallel --nonall --slf nodes.txt 'chmod -v +x ./cluster/stats-server'

# to setup on all nodes
$ parallel --nonall --slf nodes.txt './cluster/stats-server setup'

# to setup on a single node
$ ssh root@<node>
$ cd cluster
$ ./stats-server setup

# trying out the new binary using netcat
$ time ( echo '<command>' | nc <node> 12345 )
```

