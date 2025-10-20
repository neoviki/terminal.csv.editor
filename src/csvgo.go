/**
 * MIT License
 *
 * Copyright (c) 2025 Viki (VN - initials of my first and last name)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 *
 *
 * Contact: contact@viki.design  
 * Website: https://www.viki.design
 * 
 */

package main

/*
  Bugs:

    1. log.Printf - makes the table to look corrupted during save etc
    
*/

/*
Tasks:
    1. make 3 versions of data_prev_prev, data_prev data_curr ( z move to prev data y move to next data)
    2. Make 1 space after last row ( don´t fill any row - even if you have space ( this is for command area ))
*/
import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
    "bufio"
    "strconv"
	"strings"
    "path/filepath"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	data       [][]string
	app        *tview.Application
	table      *tview.Table
	inputField *tview.InputField

	selectedRow int
	selectedCol int

	editing bool

    colWidths map[int]int

	colOffset = 0
	maxVisibleCols = 5

    flex *tview.Flex
    //edit_label="[]"
    edit_label=":"
    inputFile string
	pages     *tview.Pages

    numCols int
    numRows int
    
    screenHeight int
    screenWidth int
)

const defaultColWidth = 10


func createNewPage(pageElementFlex *tview.Flex) {
    pageCount:=1
	pageID := strconv.Itoa(pageCount)
	pages.AddAndSwitchToPage(pageID, pageElementFlex, true)
}

func argParse(){
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <csv-file>")
		os.Exit(1)
	}
    
    inputFile=os.Args[1]
}

func uiInit(){
	app = tview.NewApplication()
    
    //Before every draw compute screenWidth and screenHeight ( flex will get the update value only after the first iteration - at very first iteration flex have screenHeight=0)
    app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
        screenWidth, screenHeight = screen.Size()
        return false
    })
}

func uiLoop(rootUIElement tview.Primitive) {

    app.SetRoot(rootUIElement, true).EnableMouse(true)

    err := app.Run();
    if err != nil {
        panic(err)
    }

}

func pageInit(){
	pages = tview.NewPages()
}

func tableInit(){
	table = tview.NewTable()
    //Freeze ( rows,cols) : 1 - row0 will be frozen, 2-row0 and row1 will be frozen
    table.SetFixed(1, 2)
    //able to select (row,col)
    table.SetSelectable(true, true)
    table.SetBorders(true)
}	

func inputTextBoxInit(){
    inputField = tview.NewInputField().
    SetLabel(edit_label).
    SetDoneFunc(onEditDone)

    app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
        screenWidth, _ := screen.Size()
        labelWidth := len(edit_label)
        inputField.SetFieldWidth(screenWidth - labelWidth - 1) // -1 for padding
        return false
    })

}

func flexInit(){
    flex = tview.NewFlex()
    flex.SetDirection(tview.FlexRow)
}

func flexAddTable(){
    //flex.AddItem(table, 0, 1, true)   // table fills available space

    //fillRow := screenHeight-5  // skip 2 rows ≈ 1 cm
    fillRow := 33  // skip 2 rows ≈ 1 cm
    fillCol := 1 // Utilize full screen width for column
    flex.AddItem(table, fillRow, fillCol, true)   // table fills available space

    bottomBorder := tview.NewBox().
    SetBorder(false).
    SetBackgroundColor(tcell.ColorWhite) // looks like a thick line


    flex.AddItem(bottomBorder, 1, 0, false) // simulate 1-line border

}

func flexAddInputTextBox(){
    flex.AddItem(inputField, 3, 0, true) // Add input field with height 3, focusable
}

func flexRemoveInputTextBox(){
    flex.RemoveItem(inputField)
}

func loadCSV() {
	f, err := os.Open(inputFile)
	if err != nil {
        fmt.Println("Error: opening csv file")
		os.Exit(1)
	}

	defer f.Close()

	r := csv.NewReader(f)
	data, err = r.ReadAll()
	if err != nil {
		log.Fatal(err)
        fmt.Println("Error: reading csv file")
		os.Exit(1)
	}

    numCols = len(data[0])
    numRows = len(data)
}

func getColWidth(col_nr int) int {
    w, ok := colWidths[col_nr]
	if ok{
        return w
    }
    return defaultColWidth
}

func renderTableHeader(){
    // Header row
    for c := 0; c < numCols; c++ {
        w := getColWidth(c)
        text:=wrapText(data[0][c], w)

        // For last column don´t wrap
        if c == numCols-1 {
            text=data[0][c]
        }
        cell := tview.NewTableCell(text)
        cell.SetTextColor(tcell.ColorYellow)
        cell.SetSelectable(true)
        cell.SetMaxWidth(w)
        cell.SetExpansion(0)

        // Only last column expands
        if c == numCols-1 {
            cell.SetExpansion(1)
        }

        //Add the cell content to table
        table.SetCell(0, c, cell)
        //log.Printf("Col [%d] Width : %d", c, w)
    }
}

func renderTableBody(){
    // Data rows
    for r := 1; r < numRows; r++ {
        for c := 0; c < numCols; c++ {
            w := getColWidth(c)
            
            text:=wrapText(data[r][c], w)
            
            // For last column don´t wrap
            if c == numCols-1{
                text=data[r][c]
            }

            cell := tview.NewTableCell(text)
            cell.SetMaxWidth(w)
            cell.SetExpansion(0)

            // Only last column expands
            if c == numCols-1 {
                cell.SetExpansion(1)
            }
            
            //Add the cell content to table
            table.SetCell(r, c, cell)
        }
    }
}

func renderTable() {
    table.Clear()
    renderTableHeader()
    renderTableBody()
    table.Select(selectedRow, selectedCol)
}

func wrapText(text string, width int) string {
    if width <= 0 {
        return text
    }
    // If text is too long, truncate and append ellipsis
    if len([]rune(text)) > width {
        if width > 3 {
            return string([]rune(text)[:width-3]) + "..."
        }
        return string([]rune(text)[:width])
    }
    // If text is too short, pad with spaces
    return text + strings.Repeat(" ", width-len([]rune(text)))
}


func setupKeybindings() {
    table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if editing {
            // Ignore keys while editing, inputField handles
            return event
        }

        switch event.Key() {
        case tcell.KeyRight:
            if selectedCol < len(data[0])-1 {
                selectedCol++
                table.Select(selectedRow, selectedCol)
            }
            return nil
        case tcell.KeyLeft:
            if selectedCol > 0 {
                selectedCol--
                table.Select(selectedRow, selectedCol)
            }
            return nil
        case tcell.KeyDown:
            if selectedRow < len(data)-1 {
                selectedRow++
                table.Select(selectedRow, selectedCol)
            }
            return nil
        case tcell.KeyUp:
            if selectedRow==0{
                table.Select(selectedRow, selectedCol)
                return nil
            }

            if selectedRow > 0 {
                selectedRow--
                table.Select(selectedRow, selectedCol)
            }
            return nil
        case tcell.KeyTab:
            insertColumnRight()
            refreshTable()
            return nil
        case tcell.KeyEnter:
            insertRowBelow()
            return nil
        case tcell.KeyEscape:
            confirmQuit()
            //app.Stop()
            return nil
        case tcell.KeyBackspace2, tcell.KeyBackspace:
            deleteColAfterConfirmation()
            //deleteSelectedCol()
            return nil
        }

        switch event.Rune() {
        case 'd':
            fCopyAndDelete:=func(){
                copySelectedRowToCompleted()
                deleteSelectedRow()
            }
            getUserConfirmation("Do you want to delete selected row?", fCopyAndDelete)
            return nil
        case 'e', 'i':
            startEditing()
            return nil
        case 'c':
            copyCellToClipboard()
            return nil
        case 'v':
            pasteClipboardToCell()
            return nil
        case 'x':
            cutCell()
            return nil
        case 'q':
            confirmQuit()
            //app.Stop()
            return nil
        case 'n': //n-null
            clearCell()
            return nil
        }

        return event
    })
}

func clearCell() {
    if selectedRow < 0 || selectedCol < 0 || selectedRow >= len(data) || selectedCol >= len(data[0]) {
        return
    }
    data[selectedRow][selectedCol] = ""
    refreshTable()
}

func startEditing() {
	editing = true
    flexAddInputTextBox() 
    inputField.SetText(data[selectedRow][selectedCol])
	//inputField.SetVisible(true)
	app.SetFocus(inputField)
}

func onEditDone(key tcell.Key) {
	if key == tcell.KeyEnter {
		data[selectedRow][selectedCol] = inputField.GetText()
		renderTable()
		editing = false
         // Clear the input field text
        inputField.SetText("")
	    flexRemoveInputTextBox()
		app.SetFocus(table)
        	
    } else if key == tcell.KeyEscape {
		editing = false
        inputField.SetText("")
	    flexRemoveInputTextBox()
		app.SetFocus(table)
	}


    //saveCSV(inputFile)

}

func copyCellToClipboard() {
	if selectedRow < len(data) && selectedCol < len(data[selectedRow]) {
		text := data[selectedRow][selectedCol]
		err := clipboard.WriteAll(text)
		if err != nil {
			log.Printf("Clipboard write failed: %v", err)
		}
	}
}

func pasteClipboardToCell() {
	if selectedRow < len(data) && selectedCol < len(data[selectedRow]) {
		text, err := clipboard.ReadAll()
		if err == nil {
			data[selectedRow][selectedCol] = text
			renderTable()
		}
	}
}

func cutCell() {
	if selectedRow < len(data) && selectedCol < len(data[selectedRow]) {
		text := data[selectedRow][selectedCol]
		err := clipboard.WriteAll(text)
		if err == nil {
			data[selectedRow][selectedCol] = ""
			renderTable()
		}
	}
}

func copySelectedRowToCompleted() {
    row, _ := table.GetSelection()
    
    if row == 0{
        //do not copy header row ( we are already it below when the completed csv file isn´t created)
        return
    }

    if row <= 0 || row >= len(data) {
        // Do not copy header or out-of-range
        return
    }

    completedFile := inputFile + ".completed.csv"

    // Check if file exists
    _, err := os.Stat(completedFile)
    fileExists := err == nil

    // Open file for appending or creating
    f, err := os.OpenFile(completedFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Error opening completed file: %v", err)
        return
    }
    defer f.Close()

    writer := csv.NewWriter(f)

    // If file does not exist, write header first
    if !fileExists {
        if err := writer.Write(data[0]); err != nil {
            log.Printf("Error writing header to completed file: %v", err)
            return
        }
    }

    // Write the selected row
    if err := writer.Write(data[row]); err != nil {
        log.Printf("Error writing to completed file: %v", err)
        return
    }
    writer.Flush()

    if err := writer.Error(); err != nil {
        log.Printf("Error flushing to completed file: %v", err)
        return
    }

	//loadCSV()
    renderTable()
    app.SetFocus(table)
}

func insertColumnRight() {
	row, col := table.GetSelection()

	// Sanity checks
	if len(data) == 0 || row < 0 || col < 0 || row >= len(data) || col >= len(data[0]) {
		return
	}

	// Determine where to insert: after current column
	insertAt := col + 1

	// Insert empty string into each row at insertAt position
	for i := range data {
		if insertAt >= len(data[i]) {
			// Append if insertAt is at the end
			data[i] = append(data[i], "")
		} else {
			// Insert in the middle
			data[i] = append(data[i][:insertAt+1], data[i][insertAt:]...)
			data[i][insertAt] = ""
		}
	}

	// Update selection to the new column
	selectedRow = row
	selectedCol = insertAt
    numCols+=1
	//saveCSV(inputFile)


	// Re-render table
	renderTable()
	app.SetFocus(table)
}

func insertRowBelow() {
    row, _ := table.GetSelection()

    // Prevent inserting before header row (row 0 is usually the header)
    if row < 0 || row >= len(data) {
        return
    }

    // Create a new empty row with the same number of columns as the header
    newRow := make([]string, len(data[0]))

    // Insert the new row below the selected one
    if row+1 >= len(data) {
        data = append(data, newRow)
        selectedRow = len(data) - 1
    } else {
        data = append(data[:row+1], append([][]string{newRow}, data[row+1:]...)...)
        selectedRow = row + 1
    }

    selectedCol = 0
    numRows+=1

    //saveCSV(inputFile)

    // Re-render table
    renderTable()
    app.SetFocus(table)
}


func deleteSelectedRow() {
    row, _ := table.GetSelection()

    if row == 0 {
        //do allow to delete header row
        //log.Printf("Error: Header row cannot be deleteD")
        return
    }
    // Don't delete header or if data is already minimal
    if row <= 0 || len(data) <= 1 || row >= len(data) {
        return
    }

    // Remove the row
    data = append(data[:row], data[row+1:]...)

    // Adjust selection
    if row >= len(data) {
        selectedRow = len(data) - 1
    } else {
        selectedRow = row
    }
    selectedCol = 0
    numRows -= 1
    //saveCSV(inputFile)
    
    // Re-render table
    renderTable()
    app.SetFocus(table)
}

func deleteSelectedCol() {
	row, col := table.GetSelection()

	// Sanity checks
	if len(data) == 0 || row < 0 || col < 0 || col >= len(data[0]) {
		return
	}

	// Avoid deleting if only one column left
	if len(data[0]) <= 1 {
		return
	}

	// Remove the column at index col in every row
	for i := range data {
		data[i] = append(data[i][:col], data[i][col+1:]...)
	}

	// Adjust selected column if needed
	if selectedCol >= len(data[0]) {
		selectedCol = len(data[0]) - 1
	}
	selectedRow = row

	numCols -= 1

	//saveCSV(inputFile)

	// Re-render table
	renderTable()
	app.SetFocus(table)
}

func getConfigPath(csvPath string) string {
    dir, file := filepath.Split(csvPath)
    base := strings.TrimSuffix(file, filepath.Ext(file))
    return filepath.Join(dir, base+".config")
}

func forceWriteConfig(path string) {
    log.Println("Config file corrupted. Overwriting with current data...")
    err := os.Remove(path)
    if err != nil {
        log.Printf("Error removing bad config file: %v\n", err)
        return
    }
    createConfig(path)
}

func createConfig(path string) {
    fmt.Println("Creating new config file...")
    file, err := os.Create(path)
    if err != nil {
        fmt.Printf("Error creating config file: %v\n", err)
        return
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    for col:=0;col<=10;col++{
        fmt.Fprintf(writer, "%d:%d\n", col, defaultColWidth)
    }
    writer.Flush()
}

func loadCSVConfig() {
    path := getConfigPath(inputFile)
     
    widths := make(map[int]int)
    file, err := os.Open(path)
    if err != nil {
        log.Printf("No config file found or error loading it, using default widths\n")
        createConfig(path)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        parts := strings.Split(line, ":")
        if len(parts) != 2 {
            continue
        }
        colNum, err1 := strconv.Atoi(parts[0])
        w, err2 := strconv.Atoi(parts[1])
        if err1 == nil && err2 == nil {
            widths[colNum] = w
        }
    }

    err = scanner.Err() 
    
    if err != nil {
        log.Printf("Config parsing error\n")
        forceWriteConfig(path)
        return
    }else{
        log.Printf("Loaded config [%s]\n", path)
        colWidths=widths
    }
}

func refreshTable(){
    renderTable()
    app.SetFocus(table)
}

func saveCSV(filename string) {
    tempFile := filename + ".tmp"

    f, err := os.Create(tempFile)
    if err != nil {
        log.Printf("Error creating temp CSV: %v", err)
        return
    }
    defer f.Close()

    w := csv.NewWriter(f)
    err = w.WriteAll(data)
    if err != nil {
        log.Printf("Error writing CSV data: %v", err)
        return
    }
    w.Flush()

    err = os.Rename(tempFile, filename)
    if err != nil {
        log.Printf("Error renaming temp file: %v", err)
        return
    }

    // Re-read the saved file into `data`
    file, err := os.Open(filename)
    if err != nil {
        log.Printf("Error reopening saved CSV: %v", err)
        return
    }
    defer file.Close()

    //Reloading CSV
    loadCSV()
    renderTable()
    app.SetFocus(table)
}

func magnify_cell_full_screen(app *tview.Application, content string, table *tview.Table, pages *tview.Pages) {
	textView := tview.NewTextView()
	textView.SetText(content).
		SetTextColor(tcell.ColorWhite).
		SetBackgroundColor(tcell.ColorBlack).
		SetBorder(true).
		SetTitle("Cell Content Full Screen (Press Esc to close)")

	textView.ScrollToBeginning()

	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pages.RemovePage("magnify")
			app.SetFocus(table)
			return nil
		}
		return event
	})

	modal := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(nil, 1, 0, false). // top padding: 1 row
		AddItem(
			tview.NewFlex().
				AddItem(nil, 2, 0, false).      // left padding: 2 cols
				AddItem(textView, 0, 1, true).  // main content
				AddItem(nil, 2, 0, false),      // right padding: 2 cols
			0, 1, true).
		AddItem(nil, 1, 0, false) // bottom padding: 1 row

	pages.AddPage("magnify", modal, true, true)
	app.SetFocus(textView)
}

func simulateRightArrowKeyPressEvent(){
    go func() {
        app.QueueEvent(tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone))
    }()
}

func getUserConfirmation(text string, callback func()) {
    modal := tview.NewModal().
        SetText(text).
        AddButtons([]string{"Yes", "No"}).
        SetButtonBackgroundColor(tcell.ColorDarkCyan).
        SetButtonStyle(tcell.StyleDefault.
            Foreground(tcell.ColorWhite).
            Background(tcell.ColorDarkCyan)).
        SetButtonActivatedStyle(tcell.StyleDefault.
            Foreground(tcell.ColorYellow).
            Background(tcell.ColorDarkCyan).
            Bold(true)).
        SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
                callback()
            } 
            pages.RemovePage("confirm")
        })

    pages.AddPage("confirm", modal, true, true)

    // Optional: preselect "No"
    simulateRightArrowKeyPressEvent()
}

func deleteColAfterConfirmation() {
	modal := tview.NewModal().
		SetText("Do you want to delete selected col?").
		AddButtons([]string{"Yes", "No"}).
		SetButtonBackgroundColor(tcell.ColorDarkCyan).
		SetButtonStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDarkCyan)).
		SetButtonActivatedStyle(tcell.StyleDefault.
			Foreground(tcell.ColorYellow).
			Background(tcell.ColorDarkCyan).
			Bold(true)).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
                deleteSelectedCol()
			}
			pages.RemovePage("confirm")
		})

	pages.AddPage("confirm", modal, true, true)

	// Simulate right arrow key to preselect "No"
    simulateRightArrowKeyPressEvent()

}



func confirmQuit() {
	modal := tview.NewModal().
		SetText("Do you want to close the application?").
		AddButtons([]string{"Yes", "No"}).
		SetButtonBackgroundColor(tcell.ColorDarkCyan).
		SetButtonStyle(tcell.StyleDefault.
			Foreground(tcell.ColorWhite).
			Background(tcell.ColorDarkCyan)).
		SetButtonActivatedStyle(tcell.StyleDefault.
			Foreground(tcell.ColorYellow).
			Background(tcell.ColorDarkCyan).
			Bold(true)).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Yes" {
	            saveCSV(inputFile)
				app.Stop()
			} else {
				pages.RemovePage("confirm")
			}
		})

	pages.AddPage("confirm", modal, true, true)

    // Optional: preselect "No"
    simulateRightArrowKeyPressEvent()
}

/*
func refresh(){
    renderTable()
    app.SetFocus(table)
    flexAddTable()
}
*/

func main() {
    argParse()
    loadCSV()
    loadCSVConfig()
    uiInit()
    pageInit()
    tableInit()
    inputTextBoxInit()
    renderTable()
    flexInit()
    flexAddTable()
    createNewPage(flex)
    setupKeybindings()
    uiLoop(pages)
}


