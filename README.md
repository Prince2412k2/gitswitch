# gitswitch

Manage multiple git identities. No more committing with your work email on personal projects.

## what it does

- keeps git profiles (name, email, SSH key) in `~/.config/git_conf/*.toml`
- lets you pick one interactively and applies it to the current repo's local `.git/config`
- sets `core.sshCommand` so the right SSH key is used per repo

## install

```bash
git clone https://github.com/Prince2412k2/gitswitch
cd gitswitch
go build -o gitswitch .
mv gitswitch /usr/local/bin/   # or anywhere in your PATH
```

## usage

```bash
# inside any git repo — pick a profile and apply it
gitswitch

# create a new profile (interactive form)
gitswitch new

# view all saved profiles
gitswitch view
```

## profile format

Profiles live at `~/.config/git_conf/<name>.toml`:

```toml
name        = "Priya Sharma"
email       = "priya@work.com"
ssh_key     = "~/.ssh/id_ed25519_work"
github_user = "priya-work"           # optional
github_url  = "https://github.com/priya-work"  # optional
notes       = "day job"              # optional
```

You can create them manually or via `gitswitch new`.

## what gets written to .git/config

```ini
[user]
    name  = Priya Sharma
    email = priya@work.com

[core]
    sshCommand = ssh -i ~/.ssh/id_ed25519_work -o IdentitiesOnly=yes
```

Only the local config is touched — your global `~/.gitconfig` is never modified.

## tips

- Run `gitswitch` right after cloning a new repo
- Pair with a shell hook to warn when no local user is set (see below)
- SSH keys: generate separate keys per account with `ssh-keygen -t ed25519 -C "work" -f ~/.ssh/id_ed25519_work`

### optional shell warning

Add to your `.zshrc` / `.bashrc` to warn when you're in a repo with no local identity:

```bash
function _gitswitch_warn() {
  if git rev-parse --git-dir &>/dev/null; then
    local name=$(git config --local user.name 2>/dev/null)
    if [[ -z "$name" ]]; then
      echo "⚠ no local git identity set — run: gitswitch"
    fi
  fi
}
chpwd_functions+=(_gitswitch_warn)   # zsh
# or: PROMPT_COMMAND="_gitswitch_warn; $PROMPT_COMMAND"  # bash
```
