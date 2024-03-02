package main

import (
  "github.com/charmbracelet/bubbles/list"
  "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/key"
  "github.com/charmbracelet/bubbles/table"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

var (
  focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
  blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
  cursorStyle         = focusedStyle.Copy()
  noStyle             = lipgloss.NewStyle()
  helpStyle           = blurredStyle.Copy()
  cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

  docStyle = lipgloss.NewStyle().Width(100).Height(100).Align(lipgloss.Center)
)

type User struct {
  id      string
  help    help.Model
  KeyMap  keyMap
  table   table.Model
  decks   []*Deck // table -> decks
}

func (u *User) UpdateTable() {
  i := u.table.Cursor()
  currRows := u.table.Rows()
  
  rows := []table.Row{}
  for j, _ := range currRows {
    if j == i {
      rows = append(rows, table.Row{u.decks[i].Name(), 
                                    u.decks[i].NumNew(), 
                                    u.decks[i].NumLearning(),
                                    u.decks[i].NumReview()})
    } else {
      rows = append(rows, currRows[j])
    }
  }
  sg_user.table.SetRows(rows)
}

func NewUser() *User {
	help := help.New()
	help.ShowAll = false
	return &User{help: help, KeyMap: DefaultKeyMap(),}
}

func (u *User) Init() tea.Cmd {
  return nil
}

func (u *User) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmd tea.Cmd
  switch msg := msg.(type) {
    case tea.KeyMsg:
      switch {
        case key.Matches(msg, u.KeyMap.Quit):
          return u, tea.Quit
        case key.Matches(msg, u.KeyMap.Open):
          i := u.table.Cursor()
          u.decks[i].rdata = ReviewData{}
          return u.decks[i].Update(nil)
        case key.Matches(msg, u.KeyMap.Review):
          i := sg_user.table.Cursor()
          u.decks[i].StartReview()
          return u.decks[i].Update(nil)
        case key.Matches(msg, u.KeyMap.New):
          newDeck := NewDeck("New Deck", "new.json", []list.Item{})
          sg_user.decks = append(sg_user.decks, newDeck)
          sg_user.table.SetRows(updateRows())
        case key.Matches(msg, u.KeyMap.Back):
          return u.Update(nil)
        case key.Matches(msg, u.KeyMap.ShowFullHelp):
          fallthrough
        case key.Matches(msg, u.KeyMap.CloseFullHelp):
          u.help.ShowAll = !u.help.ShowAll
      }
    case tea.WindowSizeMsg:
      h, v := docStyle.GetFrameSize()
      docStyle = docStyle.Width(msg.Width - h).Height(msg.Height - v)
    case Form:
      i := sg_user.table.Cursor()
      card := u.decks[i].Cards.Items()[msg.index]
      if msg.edit {
        msg.EditCard(card.(*Card))
      } else {
        u.decks[i].Cards.InsertItem(0, msg.CreateCard())
        u.decks[i].NumNewInc()
      }
      return u.decks[i].Update(nil)
  }

  u.table, cmd = u.table.Update(msg)

  return u, cmd
}

func (u *User) View() string {
  logoStyle := lipgloss.NewStyle().
                Bold(true).
                Foreground(lipgloss.Color("0")).
                MarginBottom(1)
  helpStyle := lipgloss.NewStyle().Align(lipgloss.Left).Width(58)

  gokiLogo := `   ________        __    __  
  /  _____/  ____ |  | _|__|
 /   \  ___ /    \|  |/ /  |
 \    \_\  |  /\  |    <|  |
  \______  /\____/|__|_ \__|
         \/            \/   `

  pageLeft := lipgloss.JoinVertical(
    lipgloss.Center,
    u.table.View(),             // Render the table
    helpStyle.Render(u.help.View(u)),
  )

  page := lipgloss.JoinVertical(
    lipgloss.Center,            // Center page
    logoStyle.Render(gokiLogo), // Render the logo
    pageLeft,
    "",
  )
  return docStyle.Render(page)
}
