# BRB

A tool for streamers to say they'll brb with a countdown.

## Scope

- [x] Customizable messages
- [x] Twitch follower count
- [x] Twitter follower count
- [x] Set up a `~/.brb` file for reading secrets
- [ ] Prettier print

## Install

### From Source

Requires Go

```sh
go install .
```

## Run

Just run `brb` with a duration such as 5s for 5 seconds, 2m for 2 minutes, 1h for 1 hour

```sh
brb 10s
```

### Twitch follower count

To enable twitch follower count, you need to set 4 variables in a `~/.brb` file

```sh
twitchclientid: <twitchclientid>
twitchclientsecret: <twitchclientsecret>
twitchappaccesstoken: <twitchappaccesstoken>
twitchuserid: <twitchuserid>
```

### Twitter follower count

To enable twitter follower count, you need to set 2 variables in a `~/.brb` file

```sh
twitterbearertoken: <twitterbearertoken>
twitterusername: <twitterusername>
```
