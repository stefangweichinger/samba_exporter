# samba_exporter

A prometheus exporter for statistic data of the samba file server

**Still in development, but basic metrics are working. Tested only on Ubuntu 20.04**

## Concept

Since the tool [smbstatus](https://www.samba.org/samba/docs/current/man-html/smbstatus.1.html) provided by the [samba](https://www.samba.org/) project can only run with elevated rights, and a [prometheus exporter](https://prometheus.io/docs/instrumenting/writing_exporters/) always exposes a public network endpoint, the samba_exporter package contains two services.

- **samba_exporter**: The prometheus exporter service that exposes the public network endpoint for the prometheus server to collect the data running as non privileged user
- **samba_statusd**: The service uses `smbstatus` to collect data and return it when requested running as privileged user

Both services can communicate using a named pipe owned by a common group.

## Build and manual install

To build the project you need [Go](https://golang.org/) Version 1.14.x and [Java](https://java.com/) Version 11 on your development machine. 
On your target machine, the samba server you want to monitor, you need [samba](https://www.samba.org/) and [systemd](https://www.freedesktop.org/wiki/Software/systemd/) installed.

To build the software change to the repositories directory and run:

```sh
./gradlew build
```

For manual install on the `target` machine do the following copies:

```sh
scp ./bin/samba_exporter <target>:/usr/bin/samba_exporter
scp ./bin/samba_statusd <target>:/usr/bin/samba_statusd 
scp ./install/usr/bin/start_samba_statusd.sh <target>:/usr/bin/start_samba_statusd.sh
scp ./install/etc/systemd/system/samba_statusd.service <target>:/etc/systemd/system/samba_statusd.service
scp ./install/etc/systemd/system/samba_exporter.service <target>:/etc/systemd/system/samba_exporter.service
```

Now login to your target machine and run the commands below to enable the services and create the needed user and group:

```sh
systemctl daemon-reload
systemctl enable samba_statusd.service
systemctl enable samba_exporter.service
addgroup --system samba-exporter
adduser --system --no-create-home --disabled-login samba-exporter
adduser samba-exporter samba-exporter
```

Finally you are abel to start the services:

```sh
systemctl start samba_statusd.service
systemctl start samba_exporter.service
```

Test the `samba_exporter` by requesting metrics with `curl`:

```sh
curl http://127.0.0.1:9922/metrics 
```
