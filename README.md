### BACKSTAGE HOOK IS IN DEVELOPMENT, IT IS NOT YET FUNCTIONAL.

# backstage-hook
Backstage-hook allows Spotify Backstage plugins to execute commands on your machine. This allows for quick plugin prototyping without compromising security

## Documentation
To learn more about backstage-hook, please go to [our documentation]().

## Installation
**Install backstage-hook on your machine:**
```bash
wget https://github.com/kuberty/kuberty/releases/... -O /usr/local/backsage-hook && sudo chmod /usr/local/backstage-hook
```
Installation instructions for [Windows](), [Mac]() and [Docker]().

**Enable the Backstage plugin on your Backstage server:**
```bash
# TODO
```

## Getting started
To run the backstage hook execute (assuming backstage is running on localhost:3000):
```bash
backstage-hook start http://localhost:3000 
```

When a plugin requests to execute a command, the request will show up in your terminal where you can accept or deny it.

## Plugins
The following plugins use backstage-hook:
- None yet