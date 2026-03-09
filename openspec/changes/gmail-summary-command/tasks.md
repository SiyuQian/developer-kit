## 1. Core summary command

- [ ] 1.1 Add `summaryCmd` to `internal/gmail/commands.go` with flags: `--channel`, `--dm`, `--no-mark-read`
- [ ] 1.2 Implement email fetching logic: query `is:unread after:YYYY/MM/DD` using `ListAllMessageIDs`, then concurrently fetch each message via `GetMessage` (bounded semaphore, 10 goroutines), with body truncation at 1000 chars
- [ ] 1.3 Build the prompt string from fetched emails (headers + truncated body) and invoke `claude -p` via `exec.Command`
- [ ] 1.4 Implement Slack delivery: if `--channel` or `--dm` is set, call `devpilot slack send` via exec
- [ ] 1.5 Implement mark-read: call `BatchModify` on all processed message IDs unless `--no-mark-read` is set or `claude -p` failed
- [ ] 1.6 Add error handling: check `claude` on PATH, handle empty output, gate mark-read on summary success

## 2. Tests

- [ ] 2.1 Unit test for prompt building (email concatenation + truncation)
- [ ] 2.2 Unit test for date query construction (today's date in YYYY/MM/DD format)
- [ ] 2.3 Integration test: summary with mock HTTP server (Gmail API) and mock claude command

## 3. Cleanup

- [ ] 3.1 Delete `.claude/skills/email-assistant/` directory
- [ ] 3.2 Update `CLAUDE.md` CLI commands section to include `devpilot gmail summary` and its flags
