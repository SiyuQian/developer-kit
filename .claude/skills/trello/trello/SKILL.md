---
name: developerkit:trello
description: Interact with Trello boards, lists, and cards directly from Claude Code. Use when the user wants to view boards, search/create/move/update Trello cards, add comments, or get a board overview. Triggers on any mention of Trello, kanban boards, task cards, or project board management.
---

# Trello

Manage Trello boards and cards. This skill works with the `trello-mcp-server` MCP server which provides structured tools for all Trello operations.

## Setup

The `trello-mcp-server` MCP must be configured in Claude Code settings. It requires two environment variables:

- `TRELLO_API_KEY` — Get from https://trello.com/power-ups/admin → New → API Key
- `TRELLO_TOKEN` — Click "Token" link on that same page

Claude Code MCP config (`.claude/settings.json`):

```json
{
  "mcpServers": {
    "trello": {
      "command": "node",
      "args": ["/path/to/trello-mcp-server/dist/index.js"],
      "env": {
        "TRELLO_API_KEY": "your-key",
        "TRELLO_TOKEN": "your-token"
      }
    }
  }
}
```

## Available MCP Tools

| Tool | Action | Key params |
|------|--------|------------|
| `trello_list_boards` | List all open boards | — |
| `trello_get_board` | Board overview with lists and cards | `board_id` or `board_name` |
| `trello_list_cards` | List cards in a list | `list_id` |
| `trello_search_cards` | Search cards by keyword | `query`, optional `board_id` |
| `trello_get_card` | Full card details | `card_id` |
| `trello_create_card` | Create a new card | `list_id`, `name`, optional `desc`, `due`, `label_ids` |
| `trello_move_card` | Move card to another list | `card_id`, `list_id` |
| `trello_add_comment` | Add comment to a card | `card_id`, `text` |
| `trello_get_board_labels` | Get labels on a board | `board_id` |
| `trello_get_board_members` | Get board members | `board_id` |

## Workflows

### "Show me my boards"

Call `trello_list_boards`.

### "What's on the Sprint board?"

Call `trello_get_board` with `board_name="Sprint"`. It returns all lists, cards, and labels in one call.

### "Find cards about authentication"

Call `trello_search_cards` with `query="authentication"`.

### "Create a bug card on the Backend board in To Do"

1. `trello_get_board` with `board_name="Backend"` — note the list IDs
2. Find the "To Do" list ID from the response
3. `trello_create_card` with `list_id`, `name`, and `desc`

### "Move the login fix card to Done"

1. `trello_search_cards` with `query="login fix"` — get card ID
2. `trello_get_board` to find the "Done" list ID
3. `trello_move_card` with `card_id` and `list_id`

### "Add a comment on the deploy card: PR merged"

1. `trello_search_cards` with `query="deploy"` — get card ID
2. `trello_add_comment` with `card_id` and `text="PR merged"`

## Name Resolution

Users refer to boards, lists, and cards by name. The `trello_get_board` tool accepts `board_name` directly. For cards and lists, use search or board overview to resolve IDs first.
