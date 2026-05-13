# Status Update Notifier (RSS, Linux only)

A simple RSS feed reader that checks a feed for updates and sends a desktop notification when a new matching item appears.

This project uses the freedesktop.org desktop notifications interface via libnotify-compatible tools, so it requires Linux with a running notification daemon compatible with `org.freedesktop.Notifications`.

## Config

Configuration is stored in a `.env` file with the following variables:

- `URL` — RSS feed URL to monitor
- `FLAG` — text marker used to identify the status entry, for example `Status:`
- `CHECK_INTERVAL` — polling interval in seconds
- `STATUS_INDEX` — index of the status value in the item title; for `Status: Online`, the status index is `1`

## How it works

The script uses [`github.com/mmcdole/gofeed`](https://github.com/mmcdole/gofeed) to parse the RSS feed.

It periodically checks the feed, compares new items against previously seen entries, and sends a desktop notification when a new item containing the configured flag is found.

## Usage

This tool calls the `notify-send` command-line program to send desktop notifications, so `notify-send` must be available in your `$PATH`.

1. Clone the repository and navigate to the project directory.

```bash
git clone https://github.com/iwanttobelievee/status-update-notifier.git
cd status-update-notifier
```

2. Create a `.env` file:

```text
URL=https://example.com/rss
FLAG=Status:
CHECK_INTERVAL=60
STATUS_INDEX=1
```

3. Install notification support.

### Fedora

```bash
sudo dnf install libnotify
```

### Debian / Ubuntu

```bash
sudo apt install libnotify-bin
```

4. Run the project:

```bash
go run .
```

## Notes

- This project does not require GNOME specifically.
- Notifications will work in GNOME, KDE Plasma, XFCE, and other environments that provide a compatible notification daemon.
- If notifications do not appear, make sure a notification daemon is running in your graphical session.

## License

MIT
