package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenLogin screen = iota
	screenMenu
	screenMovies
	screenSearch
	screenCreate
)

type model struct {
	screen        screen
	emailInput    textinput.Model
	passwordInput textinput.Model
	status        string
	token         string
	movies        []Movie

	// Recherche
	searchInput textinput.Model

	// Création
	titleInput       textinput.Model
	directorInput    textinput.Model
	genreInput       textinput.Model
	yearInput        textinput.Model
	descriptionInput textinput.Model

	// Index du champ actif dans la création
	createFieldIndex int
}

func initialModel() model {
	email := textinput.New()
	email.Placeholder = "email"
	email.Focus()

	password := textinput.New()
	password.Placeholder = "password"
	password.EchoMode = textinput.EchoPassword

	search := textinput.New()
	search.Placeholder = "mot-clé"

	title := textinput.New()
	title.Placeholder = "Titre"

	director := textinput.New()
	director.Placeholder = "Réalisateur"

	genre := textinput.New()
	genre.Placeholder = "Genre"

	year := textinput.New()
	year.Placeholder = "Année"

	desc := textinput.New()
	desc.Placeholder = "Description"

	return model{
		screen:           screenLogin,
		emailInput:       email,
		passwordInput:    password,
		status:           "Entrez vos identifiants puis Enter",
		searchInput:      search,
		titleInput:       title,
		directorInput:    director,
		genreInput:       genre,
		yearInput:        year,
		descriptionInput: desc,
		createFieldIndex: 0,
	}
}

type loginMsg struct {
	token string
	err   error
}

type moviesMsg struct {
	movies []Movie
	err    error
}

type createMsg struct {
	err error
}

func (m *model) createFields() []*textinput.Model {
	return []*textinput.Model{
		&m.titleInput,
		&m.directorInput,
		&m.genreInput,
		&m.yearInput,
		&m.descriptionInput,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			if m.screen != screenLogin && m.token != "" {
				m.screen = screenMenu
				m.status = "Menu principal"
			}
			return m, nil

		case "tab":
			if m.screen == screenLogin {
				if m.emailInput.Focused() {
					m.emailInput.Blur()
					m.passwordInput.Focus()
				} else {
					m.passwordInput.Blur()
					m.emailInput.Focus()
				}
				return m, nil
			}

		case "enter":
			if m.screen == screenLogin {
				email := m.emailInput.Value()
				password := m.passwordInput.Value()
				m.status = "Connexion..."
				return m, loginCmd(email, password)
			}

			if m.screen == screenSearch {
				query := m.searchInput.Value()
				m.status = "Recherche..."
				return m, searchMoviesCmd(m.token, query)
			}

			if m.screen == screenCreate {
				return m, createMovieCmd(
					m.token,
					m.titleInput.Value(),
					m.directorInput.Value(),
					m.genreInput.Value(),
					m.yearInput.Value(),
					m.descriptionInput.Value(),
				)
			}
		}

		// Navigation dans le formulaire de création
		if m.screen == screenCreate {
			fields := m.createFields()

			switch msg.String() {
			case "tab":
				fields[m.createFieldIndex].Blur()
				m.createFieldIndex = (m.createFieldIndex + 1) % len(fields)
				fields[m.createFieldIndex].Focus()
				return m, nil

			case "shift+tab":
				fields[m.createFieldIndex].Blur()
				m.createFieldIndex--
				if m.createFieldIndex < 0 {
					m.createFieldIndex = len(fields) - 1
				}
				fields[m.createFieldIndex].Focus()
				return m, nil
			}
		}

		// MENU
		if m.screen == screenMenu {
			switch msg.String() {
			case "1":
				m.status = "Chargement..."
				m.screen = screenMovies
				return m, getMoviesCmd(m.token)

			case "2":
				m.status = "Entrez un mot-clé"
				m.screen = screenSearch
				m.searchInput.Focus()
				return m, nil

			case "3":
				m.status = "Formulaire de création"
				m.screen = screenCreate
				m.createFieldIndex = 0
				fields := m.createFields()
				fields[0].Focus()
				return m, nil
			}
		}

	case loginMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur login: %v", msg.err)
			return m, nil
		}
		m.token = msg.token
		m.status = "Connecté. Choisissez une option."
		m.screen = screenMenu
		return m, nil

	case moviesMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur films: %v", msg.err)
			return m, nil
		}
		m.movies = msg.movies
		m.status = fmt.Sprintf("Films chargés (%d)", len(m.movies))
		return m, nil

	case createMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur création: %v", msg.err)
			return m, nil
		}
		m.status = "Film créé !"
		m.screen = screenMenu
		return m, nil
	}

	// Update inputs selon l'écran
	switch m.screen {
	case screenLogin:
		var c1, c2 tea.Cmd
		m.emailInput, c1 = m.emailInput.Update(msg)
		m.passwordInput, c2 = m.passwordInput.Update(msg)
		return m, tea.Batch(c1, c2)

	case screenSearch:
		var c tea.Cmd
		m.searchInput, c = m.searchInput.Update(msg)
		return m, c

	case screenCreate:
		var cmds []tea.Cmd
		var cmd tea.Cmd

		fields := m.createFields()

		*fields[m.createFieldIndex], cmd = fields[m.createFieldIndex].Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	return m, nil
}

func (m model) View() string {
	switch m.screen {

	case screenLogin:
		return fmt.Sprintf(
			"Login\n\nEmail: %s\nPassword: %s\n\n[enter] pour se connecter\n\n%s\n",
			m.emailInput.View(),
			m.passwordInput.View(),
			m.status,
		)

	case screenMenu:
		return fmt.Sprintf(
			"=== Menu ===\n\n"+
				"1. Voir les films\n"+
				"2. Rechercher un film\n"+
				"3. Ajouter un film\n"+
				"q. Quitter\n\n"+
				"[esc] Retour au menu depuis n'importe où\n\n"+
				"%s\n",
			m.status,
		)

	case screenMovies:
		return renderMoviesView(m.movies, m.status)

	case screenSearch:
		return fmt.Sprintf(
			"Recherche de films\n\nMot-clé: %s\n\n[enter] pour rechercher\n[esc] Menu\n\n%s\n",
			m.searchInput.View(),
			m.status,
		)

	case screenCreate:
		return fmt.Sprintf(
			"Création d'un film\n\nTitre: %s\nRéalisateur: %s\nGenre: %s\nAnnée: %s\nDescription: %s\n\n[tab] champ suivant, [shift+tab] champ précédent\n[enter] pour créer\n[esc] Menu\n\n%s\n",
			m.titleInput.View(),
			m.directorInput.View(),
			m.genreInput.View(),
			m.yearInput.View(),
			m.descriptionInput.View(),
			m.status,
		)
	}

	return "Écran inconnu"
}
