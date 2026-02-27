# Agent Team Blueprints Skill — Design Document

**Date:** 2026-02-25
**Status:** Approved
**Component:** `skills/agent-teams/`

## Summary

A Claude Code skill that guides users through selecting, configuring, and deploying Agent Team blueprints. It provides 5 pre-built orchestration patterns covering the most common multi-agent scenarios, generating project-specific configuration files through guided dialogue.

## Design Decisions

### Delivery Form: Skill (not CLI tool)

- Fits the developer kit's existing pattern (pm skill, trello skill)
- Activates naturally through conversation context
- No external dependencies — pure markdown + references

### Architecture: SKILL.md + Independent Reference Files (Option C)

- SKILL.md stays concise (~150-180 lines): trigger logic, decision tree, dialogue flow, output format
- Each blueprint lives in its own reference file (~80-120 lines): loaded on demand after user selects
- Progressive disclosure: Claude only reads the reference file for the chosen blueprint
- Easy to extend: new blueprint = new reference file + decision tree entry

### Automation Level: Guided Configuration Generation

- Skill generates 3 artifacts (CLAUDE.md section, hooks, spawn prompts) through dialogue
- User confirms before files are written
- Does NOT auto-execute TeamCreate/Task spawn — user controls when to start the team
- Rationale: Agent Teams is experimental, token cost is high (~800k for 3-person team), user should make the explicit decision to start

## File Structure

```
skills/agent-teams/
├── SKILL.md
└── references/
    ├── overview.md
    ├── parallel-specialists.md
    ├── pipeline.md
    ├── swarm.md
    ├── competing-hypotheses.md
    └── plan-first-parallel.md
```

## SKILL.md Specification

### Trigger Description

Activates when user mentions: agent team, multi-agent, team collaboration, parallel development, swarm, or describes a task suited for multiple Claude instances working together.

### Dialogue Flow (4 Steps)

1. **Understand task** — User describes what they want to accomplish
2. **Recommend blueprint** — Suggest 1-2 blueprints based on task characteristics, explain reasoning
3. **Collect parameters** — Team name, teammate count, roles, file boundaries, verification commands
4. **Generate configuration** — Output 3 artifacts for user confirmation, then write to project files

### Blueprint Selection Decision Tree

```
Task characteristics:
├── Multi-angle review of same code?         → Parallel Specialists
├── Tasks have sequential dependencies?       → Pipeline
├── Many homogeneous repeatable tasks?        → Swarm
├── Complex bug needing hypothesis testing?   → Competing Hypotheses
└── Large feature needing design-first?       → Plan-First Parallel
```

## 5 Blueprints

### 1. Parallel Specialists

**When:** Multi-perspective review — code review, security audit, performance analysis, architecture review.

**Team structure:**
- Lead (coordinator)
- 3-5 Specialists, each with a different review lens, all running in parallel
- No task dependencies

**Configurable parameters:**
- Review perspectives (default: security / performance / test-coverage)
- Review scope (directory/file globs)
- Output format (markdown report / PR comments)
- Model per teammate (default: haiku for cost efficiency)

### 2. Pipeline

**When:** Full-stack feature development with sequential stages.

**Team structure:**
- Lead (coordinator)
- Stage 1: Architect (design API schema)
- Stage 2: Backend dev (implement endpoints, blocked by Stage 1)
- Stage 3: Frontend dev (integrate UI, blocked by Stage 2)
- Stage 4: Test writer (integration tests, blocked by Stage 2)
- Stages 3 & 4 can run in parallel (both depend on Stage 2 only)

**Configurable parameters:**
- Number of stages and roles (default: 4)
- Directories per stage
- Whether Lead must approve each stage output (plan approval)

### 3. Swarm

**When:** Large batches of homogeneous tasks — API migration, bulk test writing, batch refactoring.

**Team structure:**
- Lead (task creator)
- N Workers with identical spawn prompts, self-claiming tasks from TaskList

**Configurable parameters:**
- Worker count (default: 3-5)
- Task list generation method (file glob / manual list)
- Verification command per task
- Model per worker (default: haiku)

### 4. Competing Hypotheses

**When:** Complex debugging, root cause analysis, technical decision-making.

**Team structure:**
- Lead (judge)
- 2-3 Investigators, each assigned a different hypothesis

**Configurable parameters:**
- Initial hypothesis list (user-provided or Lead-generated from symptoms)
- Investigation scope (logs, code, config)
- Whether investigators can challenge each other via SendMessage

### 5. Plan-First Parallel

**When:** Large features requiring architecture design before parallel implementation.

**Team structure (two phases):**
- Phase 1: Lead spawns Architect (plan mode, read-only). Architect outputs design. Lead approves via plan_approval_request.
- Phase 2: Lead spawns N Implementers based on approved design, each owning a module.

**Configurable parameters:**
- Architect constraints (tech stack, architecture style)
- Implementer count (default: 3)
- Whether implementers also require plan approval

## Generated Artifacts

### Artifact 1: CLAUDE.md Team Section

Appended to project CLAUDE.md (never overwrites):

```markdown
## Agent Team: {team_name}

### Blueprint: {blueprint_type}
{one-line team goal}

### Module Boundaries
| Module | Owner | Files |
|--------|-------|-------|
| {module} | {teammate_name} | {file_glob} |

### Verification Commands
- Tests: `{test_command}`
- Lint: `{lint_command}`
- Build: `{build_command}`

### Conventions
- {inherited from existing CLAUDE.md}
- {blueprint-specific, e.g. "each teammate only modifies files in their module"}
```

### Artifact 2: Hook Scripts + settings.json

Two shell scripts in `.claude/hooks/`:

**task-completed.sh** — TaskCompleted hook, validates teammate output:
```bash
#!/bin/bash
{test_command}
if [ $? -ne 0 ]; then
  echo "Tests failed. Fix before marking complete." >&2
  exit 2
fi
```

**teammate-idle.sh** — TeammateIdle hook, checks for remaining unclaimed tasks.

Plus a settings.json fragment to merge:
```json
{
  "env": { "CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS": "1" },
  "hooks": {
    "TaskCompleted": [{ "hooks": [{ "type": "command", "command": ".claude/hooks/task-completed.sh" }] }],
    "TeammateIdle": [{ "hooks": [{ "type": "command", "command": ".claude/hooks/teammate-idle.sh" }] }]
  }
}
```

### Artifact 3: Spawn Prompt List

Per-teammate spawn prompt in markdown:

```markdown
### Teammate: {name}
- **Agent type**: {general-purpose | Explore | Plan}
- **Model**: {opus | sonnet | haiku}
- **Run in background**: true

**Prompt:**
> You are {role_description}.
> Your task: {specific_task}
> Files you own: {file_list}
> Files you must NOT modify: {exclusion_list}
> Verification: Run `{verify_command}` before marking any task complete.
> When done, update your task status to completed.
```

## Write Flow

1. Display all 3 artifacts → user confirms
2. Append CLAUDE.md section (never overwrite)
3. Create `.claude/hooks/` scripts (chmod +x)
4. Merge settings.json hooks (preserve existing config)
5. Output spawn prompts as text — user decides when to execute

## Edge Cases

- **No CLAUDE.md**: Create new file
- **Existing settings.json**: Deep merge hooks config, preserve other fields
- **User unsure which blueprint**: Show overview.md comparison table
- **Agent Teams not enabled**: Prompt user to set environment variable
- **User wants to customize beyond parameters**: Direct them to the reference file for the blueprint

## Reference File Format

Each `references/{blueprint}.md` contains:
1. One-paragraph description of the pattern
2. When to use / when NOT to use
3. Team structure diagram
4. Complete CLAUDE.md template with `{placeholders}`
5. Complete spawn prompt templates with `{placeholders}`
6. Hook script templates
7. Best practices specific to this pattern (team size, model selection, task density)
8. Example: a concrete end-to-end scenario

## Non-Goals

- No auto-execution of TeamCreate/Task spawn
- No runtime monitoring or dashboard
- No cross-blueprint composition (pick one blueprint per team)
- No custom blueprint creation tooling (v1)
