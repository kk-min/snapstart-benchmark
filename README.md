# snapstart-benchmark
Program to perform manual benchmarking for AWS Lambda's SnapStart.

# Usage
1. `go build main.go`
2. `./main [OPTIONS]`

# Options
- `-b`: Burst count (integer)
- `-i`: Inter-arrival time (IAT) (milliseconds)
- `-o`: Output directory (dir path)
- `snapstart`: SnapStart enabled (true/false)
