# Expense Tracker TUI

**Expense Tracker TUI** is a terminal-based expense/income tracking application, built with Go and powered by the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework. It provides an intuitive, keyboard-driven interface in the console to log, view, and analyze your financial entries.

---

## 🛠️ Features

- Add expense or income entries (date, amount, category, notes)
- Edit and deleting existing expenses  
- View a list of all expenses  
- Filter by category  
- Summary (e.g. monthly)  
- Persistent storage (JSON)
- Export to CSV
- Responsive terminal UI — works on Linux, macOS, Windows (with ANSI support)  

---

## 🚀 Installation

### Requirements

- Go 1.xx or newer  
- A terminal that supports ANSI escape codes (e.g. xterm, iTerm2, Windows terminal / WSL)  

### Install
```bash
go install github.com/Lexv0lk/expense-tracker-tui/expense-tracker@latest
```

### Run
```bash
expense-tracker
```

---

## 🧑‍💻 Usage Examples

Below are short demonstrations of the main features of **Expense Tracker TUI**.  
*(GIFs or screenshots will be added here later.)*

---

### 1. View Expenses in a Table
<p>
    <img src="https://s14.gifyu.com/images/bwZty.gif" width="100%" alt="Table View">
</p>
Displays all recorded expenses in a clean, sortable table view.  
Use arrow keys or shortcuts to navigate between entries.

---

### 2. Create a New Expense
<p>
    <img src="https://s14.gifyu.com/images/bwZ59.gif" width="100%" alt="Creating expense">
</p>
Add a new expense with amount, category, date, and description.  
Quick and intuitive — no mouse required.

---

### 3. Edit Selected Expense  
<p>
    <img src="https://s14.gifyu.com/images/bwZ5A.gif" width="100%" alt="Editing expense">
</p> 
Select an existing expense and update any field directly from the terminal.

---

### 4. Delete an Expense  
<p>
    <img src="https://s14.gifyu.com/images/bwZ5G.gif" width="100%" alt="Deleting expense">
</p> 
Remove unwanted or incorrect entries with a single key press.

---

### 5. Filter Expenses  
<p>
    <img src="https://s14.gifyu.com/images/bwZDb.gif" width="100%" alt="Filtering expenses">
</p> 
Filter expenses by category.

---

### 6. Monthly Expense Summary  
<p>
    <img src="https://s14.gifyu.com/images/bwZDL.gif" width="100%" alt="Getting summary">
</p> 
Automatically aggregates expenses by month and displays totals in a clear summary view.

---

### 7. Export to CSV  
<p>
    <img src="https://s14.gifyu.com/images/bwZDx.gif" width="100%" alt="Exporting to CSV">
</p> 
Export all expenses to a `.csv` file for external analysis or backup.


---

## License

This project is licensed under the MIT License.

---

## Contributing

Feel free to fork the repository, submit issues, and send pull requests. Contributions are welcome!

---

## Acknowledgements

This project was developed to practice Go programming as a part of Roadmap [Task](https://roadmap.sh/projects/expense-tracker).

---

For more information and to access the source code, visit the [expense-tracker-tui GitHub repository](https://github.com/Lexv0lk/expense-tracker-tui).
