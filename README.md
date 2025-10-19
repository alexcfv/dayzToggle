## Clone the repository
```bash
git clone https://github.com/alexcfv/dayzToggle.git
cd dayzToggle
```
## Find your keyboard device
```bash
sudo cat /proc/bus/input/devices | grep -A 5 -i "keyboard"
```
You will see something like:
```bash
I: Bus=0011 Vendor=0001 Product=0001 Version=ab83
N: Name="AT Translated Set 2 keyboard"
H: Handlers=sysrq kbd event12 leds
```
Note the event number, for example `event12`.
## Grant access to the keyboard device
```bash
sudo chmod a+r /dev/input/event12
```
## Write event12 in code

## Install Go dependencies and Build the binary
```bash
go mod tidy
go build -o dayzToggle main.go
```
## Run it
```bash
sudo ./main
```

---

## Troubleshooting
| Error                                               | Cause                                      | Solution                                                  |
| --------------------------------------------------- | ------------------------------------------ | --------------------------------------------------------- |
| `permission denied`                                 | No access to `/dev/input/eventX`           | Use `sudo chmod a+r /dev/input/eventX` or run with `sudo` |
| `open /dev/input/eventX: no such file or directory` | Wrong event number                         | Check the correct one with `cat /proc/bus/input/devices`  |
| `exec: not found`                                   | Go not installed or not in PATH            | Install Go using `sudo pacman -S go`                      |
| Script does nothing                                 | Wrong input device or Wayland restrictions | Verify that `/dev/input/eventX` is the correct keyboard   |
