# termban
terminal + kanban = termban

This was a project I used to learn how to build [bubbletea TUIs](https://github.com/charmbracelet/bubbletea).

It is archived and not being worked on anymore, but if you need a simple, local Kanban board in your terminal, it works pretty well!

## How to Run


## Default Key Binds
- Left and Right Arrow Keys: Navigate between columns
- Up and Down Arrow Keys: Navigate between tasks in columns
- Space: Switch between List Mode and Move Mode
  - When in move mode, your currently selected task is highlighted blue. You can then use the left and right arrow keys to move it to a different column.
  - Once done moving a task, you can hit space again to return to list mode.
- While in List Mode:
  - A: Add a new task (it will bring up a text entry field to enter the title and description, note the description won't be visible because I never implemented that part...)
  - D: Delete a task
- Escape: Exit termban