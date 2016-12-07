# Nuke It From Orbit

Deletes an MCP 2.0 network domain and all the resources within it.

```bash
nifo  --region=AU \
      --datacenter=AU9 \
      --networkdomain="My network domain"
```

Also supports `--verbose` (extra diagnostic output) and `--force` (don't prompt for confirmation).
