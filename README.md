# csvgo: terminal csv editor

**csvgo** is an **experimental terminal-based CSV editor** built in **Go**, tested on **Ubuntu Linux**.
It provides a simple spreadsheet-like interface for viewing and editing CSV files directly in the terminal.
The UI is based on [`tview`](https://github.com/rivo/tview) and [`tcell`](https://github.com/gdamore/tcell).

![DEMO](demo.gif)

---

## Features

* Edit CSV files directly in the terminal
* Insert and delete rows or columns
* Copy, cut, and paste cells
* Move with arrow keys
* Edit cell contents in an input box
* Confirmation dialogs for delete and quit actions
* Automatic config file (`.config`) for column widths
* Saves automatically on exit
* Creates `.completed.csv` file when rows are deleted (for backup/reference)

---

## Requirements

* Go 1.20 or newer
* Tested on **Ubuntu Linux [ 24.04.1 LTS ]**

---

## Build and Install

```bash
git clone https://github.com/yourusername/terminal.csv.editor.git
cd terminal.csv.editor
chmod +x install.ubuntu.sh
./install.ubuntu.sh
```

---

## Run the Application

```bash
csvgo <csv-file>
```


---

## Keyboard Shortcuts

| Key            | Action                                                                          |
| -------------- | ------------------------------------------------------------------------------- |
| **↑ ↓ ← →**    | Move selection                                                                  |
| **e** or **i** | Edit selected cell                                                              |
| **Enter**      | Insert a new row below                                                          |
| **Tab**        | Insert a new column to the right                                                |
| **Backspace**  | Delete selected column (with confirmation)                                      |
| **d**          | Delete selected row (after confirmation; also copies to `<file>.completed.csv`) |
| **c**          | Copy cell to clipboard                                                          |
| **x**          | Cut cell (copy + clear)                                                         |
| **v**          | Paste clipboard into selected cell                                              |
| **n**          | Clear selected cell (set to empty)                                              |
| **q**          | Quit (with confirmation and auto-save)                                          |
| **Esc**        | Exit edit mode or cancel dialogs                                                |

---

## Config file

For each CSV file, a config file named `<filename>.config` is automatically created.
It stores column widths in a simple `col:width` format, for example:

```
0:10
1:10
2:30
```

If the config file is missing or corrupted, it is automatically regenerated.

---

## Status

 * This project is **experimental**.
 * Currently tested only on **Ubuntu Linux (terminal mode)**.
 * Expect occasional layout or redraw issues in smaller terminals.
 * Logging with `log.Printf` can distort layout during runtime (known observation).

---

## License

MIT License 

---

## Known Issues / TODO

* `log.Printf` causes display corruption — should be replaced with a non-blocking logging option
* Add undo/redo stack for edits
* Add scroll indicators when table exceeds screen size
* Improve resize behavior for small terminal windows
* Add optional autosave toggle
* Optional read-only mode
* Optional color theme configuration
* Optional “magnify cell” full-screen view on a shortcut key
* Windows and macOS terminal support not fully tested

---

Open to collaborate — feel free to fork and work on it.



