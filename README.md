# nds-quota
> OpenNDS (formerly nodogsplash) preauth to assign per-user volume quota

It uses pre-generated vouchers stored as flat files in a directory. The voucher codes can be handed out to users of the network. Traffic consumption is tracked in regular invervals via a cron job. Users that have exceeded their quota are then deauthenticated.

The system was developed for an offsite location that shares an expensive and traffic-constrained satellite internet link with its short-term guests.

## Build
nds-quota builds a single statically-linked binary, which should run without dependencies on any system targeted by the go (cross-) compiler.
```sh
# Build for host OS/arch
make build

# Example: Cross-compile for an OpenWRT device with MIPS architecture
GOMIPS=softfloat GOOS=linux GOARCH=mipsle make build
```

## Installation
1. Install [OpenNDS](https://github.com/openNDS/openNDS)
2. Create a directory `/opt/nds-quota`
3. Copy the elf (`nds-quota`) into that directory
4. Create a `config.yml` file in the same folder (refer to `config.example.yml`)
5. Copy the `templates` directory somewhere and set its path in `config.yml`
6. Use `generate.py` to generate vouchers (`--help` for usage)
7. (optional) Adjust the template to your needs

## Configuration
1. Tell OpenNDS (nodogsplash) to use `nsd-quota` as its preauth command:
   In OpenWRT, for example, add the following line to `/etc/config/nodogsplash`
   ```
   option preauth '/opt/nds-quota/nds-quota'
   ```

2. Set up a cron job to regularly check quota and deauth clients that have exceeded theirs
   ```cron
   */2 * * * * /opt/nds-quota/nds-quota check-deauth
   ```

## FLASH WEAR WARNING!
When running nds-quota on an embedded device with only flash storage (basically every OpenWRT router), you will wear out flash, which will ultimately lead to physical failure of the device.

In such a scenario, it is recommended to store the database (`dataDirectory` in `config.yml`) on an external storage device (USB).

## nodogsplash support
The original version was written in python and is compatible with nodogsplash v4.1.0 (check `v1` branch). It was reimplemented in go due to problems running python on a router with little flash. In the process of rewriting, i discovered that nodogsplash dropped support for preauth, which is now part of OpenNDS, a fork. The preauth API of OpenNDS is slightly different from the original nodogsplash preauth API, so the current version is incompatible with old nodogsplash versions supporting preauth.
