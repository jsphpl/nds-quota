# nds-quota
> nodogsplash preauth to assign per-user volume quota (written in go)

## Build
```sh
# Build for host OS/arch
make build

# Example: Cross-compile for an OpenWRT device with MIPS architecture
GOMIPS=softfloat GOOS=linux GOARCH=mipsle make build
```

## Installation
1. Create a directory `/opt/nds-quota`
2. Copy the two binaries (`preauth`, `check-deauth`) into that directory
3. Create a `config.yml` file in the same folder (refer to `config.example.yml`)
4. Copy the `templates` directory somewhere and set its path in `config.yml`
5. (optional) Adjust the template to your needs

## Configuration
1. Tell nodogsplash to use `preauth` as its preauth command
   In OpenWRT, for example, add the following line to `/etc/config/nodogsplash`
   ```
   option preauth '/opt/nds-quota/preauth'
   ```

2. Set up a cron job to regularly check quota and deauth clients that have exceeded theirs
   ```cron
   */2 * * * * /opt/nds-quota/check-deauth
   ```
