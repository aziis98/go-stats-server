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

## Downloading the Artifact

TODO

<!--
To download the compiled binary on multiple machines using GNU Parallel, use the following `wget` command:

```bash
cat machines.txt | parallel --slf - 'wget -qO- https://github.com/aziis98/go-stats-server/releases/latest/download/stats-server | tar -xz -C /path/to/destination/'
```

Replace "your-username" and "your-repo" with your GitHub username and repository name. Adjust "/path/to/destination/" to the desired destination on the target machines.

Feel free to explore and customize the server to suit your needs. If you encounter any issues or have suggestions for improvements, please open an issue on this repository. Contributions are welcome!
-->
