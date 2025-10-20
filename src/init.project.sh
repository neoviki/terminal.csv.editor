go clean -modcache
rm go.mod
rm go.sum
go mod init neoviki_spreadsheet

get_packages()
{
    go get github.com/mattn/go-runewidth
    go get github.com/gdamore/tcell/v2
    go get github.com/atotto/clipboard
}

get_packages_latest()
{
    go get github.com/mattn/go-runewidth@latest
    go get github.com/gdamore/tcell/v2@latest
    go get github.com/atotto/clipboard@latest
}


get_packages

#get_packages_latest
#go mod tidy
#go doc github.com/rivo/tview.TextView.ScrollToBeginning

