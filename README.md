Tezos Network Monitor
=====================

![pipeline status](https://gitlab.com/polychainlabs/tezos-network-monitor/badges/master/pipeline.svg)


Monitor interesting Tezos network events related to your addresses and bakers.

While most monitoring is unique to your infrastructure, some events are best _additionally_ monitored on-chain to detect when your infrastructure _thinks_ it is online but the network disagrees, or to detect events related, but not necessarily originated within your infrastructure.  This monitor is built to fill in that gap -- and we recommend you customize further to best improve your ops.

### Usage

After adding your alerting keys to `.env` and listing your addresses in `config.yaml`, run like so:

```shell
source .env
go run .
```

### Alerts

This monitor alerts on the following:

1. **Slashing**
   - Alerts if anyone double bakes.  Pages if it was you :(
   - Alerts if anyone endorses bakes.  Pages if it was you :(
2. **Network** 
   - Alert if network is lagging (or node is unresponsive)
   - Page if network lags 60m+
3. **Missed Endorsements**
   - Alert if endorsements are missed
   - Page if 2 or more endorsements in last 20 blocks are missed
   - Page if 5 or more endorsements are missed per cycle
4. **Missed Baking**
   - Alert if you miss baking any blocks
   - Page if miss 2 or more blocks in a cycle
5. **Transactions**
   - Alert if funds are ever sent or received to our addresses
   - Page if tx are ever sent _from_ any of our addresses to a non-whitelisted destination
6. **Delegations**
   - Alert if delegations ever received or withdrawn
