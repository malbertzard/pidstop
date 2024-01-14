# pidstop

`pidstop` is a simple command-line tool for monitoring the VMRSS and other information of a process and its children.

## Features

- Monitor VMRSS (Virtual Memory Resident Set Size) of a specified process and its children.
- Display additional information such as PID, process name, state, parent PID, command, and user.
- Option to show only the entry process, excluding its children.

## Usage

```bash
pidstop -p <PID>
pidstop -n <ProcessName>
pidstop -c <CommandToRun>
pidstop -s
```

- **-p, --pid**: Specify the PID of the process to monitor.
- **-n, --name**: Specify the name of the process to monitor.
- **-c, --command**: Specify the command to run and monitor its process.
- **-s, --show-only**: Show only the information for the entry process, excluding its children.

## Examples

```bash
# Monitor process with PID 12345 and its children
pidstop -p 12345

# Monitor process with name "example" and its children
pidstop -n example

# Run a command and monitor its process
pidstop -c "your_command_here"

# Show only the information for the entry process, excluding its children
pidstop -p 12345 -s
```

## License

This project is licensed under the [MIT License](LICENSE).
