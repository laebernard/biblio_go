package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screen int

var styleTitle lipgloss.Style
var styleLabel lipgloss.Style
var styleInput lipgloss.Style

const (
	screenLogin screen = iota
	screenMenu
	screenMovies
	screenSearch
	screenCreate
	screenMovieByID
	screenUpdateMovie
	screenDeleteMovie
	screenBulkDelete
	screenBulkUpdate
	screenRegister
	screenUpdateOwnProfile
	screenUpdateUserByAdmin
	screenResetDB
	screenSettings
)

type model struct {
	screen        screen
	emailInput    textinput.Model
	passwordInput textinput.Model
	status        string
	token         string
	movies        []Movie

	searchInput textinput.Model

	titleInput       textinput.Model
	directorInput    textinput.Model
	genreInput       textinput.Model
	yearInput        textinput.Model
	descriptionInput textinput.Model

	movieIDInput        textinput.Model
	updateIDInput       textinput.Model
	deleteIDInput       textinput.Model
	bulkDeleteIDsInput  textinput.Model
	bulkUpdateIDsInput  textinput.Model
	registerNameInput   textinput.Model
	registerEmailInput  textinput.Model
	registerPassInput   textinput.Model
	profileNameInput    textinput.Model
	profileEmailInput   textinput.Model
	profilePassInput    textinput.Model
	adminUserIDInput    textinput.Model
	adminUserNameInput  textinput.Model
	adminUserEmailInput textinput.Model
	adminUserPassInput  textinput.Model

	createFieldIndex     int
	updateFieldIndex     int
	bulkUpdateFieldIndex int
	registerFieldIndex   int
	profileFieldIndex    int
	adminUserFieldIndex  int

	settingsThemeInput  textinput.Model
	settingsAPIInput    textinput.Model
	settingsAccentInput textinput.Model
	settingsFieldIndex  int
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

	movieID := textinput.New()
	movieID.Placeholder = "ID du film"

	updateID := textinput.New()
	updateID.Placeholder = "ID du film"

	deleteID := textinput.New()
	deleteID.Placeholder = "ID du film"

	bulkDeleteIDs := textinput.New()
	bulkDeleteIDs.Placeholder = "1,2,3"

	bulkUpdateIDs := textinput.New()
	bulkUpdateIDs.Placeholder = "1,2,3"

	registerName := textinput.New()
	registerName.Placeholder = "Nom"
	registerEmail := textinput.New()
	registerEmail.Placeholder = "Email"
	registerPass := textinput.New()
	registerPass.Placeholder = "Mot de passe"
	registerPass.EchoMode = textinput.EchoPassword

	profileName := textinput.New()
	profileName.Placeholder = "Nouveau nom"
	profileEmail := textinput.New()
	profileEmail.Placeholder = "Nouvel email"
	profilePass := textinput.New()
	profilePass.Placeholder = "Nouveau mot de passe"
	profilePass.EchoMode = textinput.EchoPassword

	adminUserID := textinput.New()
	adminUserID.Placeholder = "ID utilisateur"
	adminUserName := textinput.New()
	adminUserName.Placeholder = "Nouveau nom"
	adminUserEmail := textinput.New()
	adminUserEmail.Placeholder = "Nouvel email"
	adminUserPass := textinput.New()
	adminUserPass.Placeholder = "Nouveau mot de passe"
	adminUserPass.EchoMode = textinput.EchoPassword

	settingsTheme := textinput.New()
	settingsTheme.Placeholder = "light / dark"
	settingsTheme.SetValue(AppConfig.Theme)

	settingsAPI := textinput.New()
	settingsAPI.Placeholder = "API URL"
	settingsAPI.SetValue(AppConfig.APIURL)

	settingsAccent := textinput.New()
	settingsAccent.Placeholder = "#FF00FF"
	settingsAccent.SetValue(AppConfig.Accent)

	applyTheme()
	applyAccent()

	return model{
		screen:               screenLogin,
		emailInput:           email,
		passwordInput:        password,
		status:               "Entrez vos identifiants puis Enter",
		searchInput:          search,
		titleInput:           title,
		directorInput:        director,
		genreInput:           genre,
		yearInput:            year,
		descriptionInput:     desc,
		movieIDInput:         movieID,
		updateIDInput:        updateID,
		deleteIDInput:        deleteID,
		bulkDeleteIDsInput:   bulkDeleteIDs,
		bulkUpdateIDsInput:   bulkUpdateIDs,
		registerNameInput:    registerName,
		registerEmailInput:   registerEmail,
		registerPassInput:    registerPass,
		profileNameInput:     profileName,
		profileEmailInput:    profileEmail,
		profilePassInput:     profilePass,
		adminUserIDInput:     adminUserID,
		adminUserNameInput:   adminUserName,
		adminUserEmailInput:  adminUserEmail,
		adminUserPassInput:   adminUserPass,
		createFieldIndex:     0,
		updateFieldIndex:     0,
		bulkUpdateFieldIndex: 0,
		registerFieldIndex:   0,
		profileFieldIndex:    0,
		adminUserFieldIndex:  0,
		settingsThemeInput:   settingsTheme,
		settingsAPIInput:     settingsAPI,
		settingsAccentInput:  settingsAccent,
		settingsFieldIndex:   0,
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

type movieMsg struct {
	movie Movie
	err   error
}

type createMsg struct {
	err error
}

type updateMsg struct {
	err error
}

type deleteMsg struct {
	err error
}

type bulkDeleteMsg struct {
	err error
}

type bulkUpdateMsg struct {
	err error
}

type registerMsg struct {
	token string
	err   error
}

type userUpdateMsg struct {
	err error
}

type adminUserUpdateMsg struct {
	err error
}

type resetMsg struct {
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

func (m *model) updateFields() []*textinput.Model {
	return []*textinput.Model{
		&m.updateIDInput,
		&m.titleInput,
		&m.directorInput,
		&m.genreInput,
		&m.yearInput,
		&m.descriptionInput,
	}
}

func (m *model) bulkUpdateFields() []*textinput.Model {
	return []*textinput.Model{
		&m.bulkUpdateIDsInput,
		&m.titleInput,
		&m.directorInput,
		&m.genreInput,
		&m.yearInput,
		&m.descriptionInput,
	}
}

func (m *model) registerFields() []*textinput.Model {
	return []*textinput.Model{
		&m.registerNameInput,
		&m.registerEmailInput,
		&m.registerPassInput,
	}
}

func (m *model) profileFields() []*textinput.Model {
	return []*textinput.Model{
		&m.profileNameInput,
		&m.profileEmailInput,
		&m.profilePassInput,
	}
}

func (m *model) adminUserFields() []*textinput.Model {
	return []*textinput.Model{
		&m.adminUserIDInput,
		&m.adminUserNameInput,
		&m.adminUserEmailInput,
		&m.adminUserPassInput,
	}
}

func (m *model) settingsFields() []*textinput.Model {
	return []*textinput.Model{
		&m.settingsThemeInput,
		&m.settingsAPIInput,
		&m.settingsAccentInput,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			if m.screen == screenRegister {
				m.screen = screenLogin
				m.status = "Entrez vos identifiants puis Enter"
				return m, nil
			}
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

		case "f2":
			if m.screen == screenLogin {
				m.status = "Inscription utilisateur"
				m.screen = screenRegister
				m.registerFieldIndex = 0
				m.registerNameInput.SetValue("")
				m.registerEmailInput.SetValue("")
				m.registerPassInput.SetValue("")
				fields := m.registerFields()
				fields[0].Focus()
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

			if m.screen == screenMovieByID {
				id := m.movieIDInput.Value()
				if id == "" {
					m.status = "Veuillez saisir un ID"
					return m, nil
				}
				m.status = "Chargement du film..."
				return m, getMovieCmd(m.token, id)
			}

			if m.screen == screenUpdateMovie {
				id := m.updateIDInput.Value()
				if id == "" {
					m.status = "Veuillez saisir un ID"
					return m, nil
				}
				m.status = "Mise à jour..."
				return m, updateMovieCmd(m.token, id, m.titleInput.Value(), m.directorInput.Value(), m.genreInput.Value(), m.yearInput.Value(), m.descriptionInput.Value())
			}

			if m.screen == screenDeleteMovie {
				id := m.deleteIDInput.Value()
				if id == "" {
					m.status = "Veuillez saisir un ID"
					return m, nil
				}
				m.status = "Suppression..."
				return m, deleteMovieCmd(m.token, id)
			}

			if m.screen == screenBulkDelete {
				ids, err := parseIDList(m.bulkDeleteIDsInput.Value())
				if err != nil {
					m.status = fmt.Sprintf("Erreur IDs: %v", err)
					return m, nil
				}
				m.status = "Suppression en lot..."
				return m, bulkDeleteMoviesCmd(m.token, ids)
			}

			if m.screen == screenBulkUpdate {
				ids, err := parseIDList(m.bulkUpdateIDsInput.Value())
				if err != nil {
					m.status = fmt.Sprintf("Erreur IDs: %v", err)
					return m, nil
				}
				m.status = "Mise à jour en lot..."
				return m, bulkUpdateMoviesCmd(m.token, ids, m.titleInput.Value(), m.directorInput.Value(), m.genreInput.Value(), m.yearInput.Value(), m.descriptionInput.Value())
			}

			if m.screen == screenRegister {
				name := m.registerNameInput.Value()
				email := m.registerEmailInput.Value()
				password := m.registerPassInput.Value()
				if name == "" || email == "" || password == "" {
					m.status = "Tous les champs sont requis"
					return m, nil
				}
				m.status = "Inscription..."
				return m, registerCmd(name, email, password)
			}

			if m.screen == screenUpdateOwnProfile {
				m.status = "Mise à jour du profil..."
				return m, updateOwnProfileCmd(m.token, m.profileNameInput.Value(), m.profileEmailInput.Value(), m.profilePassInput.Value())
			}

			if m.screen == screenUpdateUserByAdmin {
				id := m.adminUserIDInput.Value()
				if id == "" {
					m.status = "Veuillez saisir un ID utilisateur"
					return m, nil
				}
				m.status = "Mise à jour utilisateur..."
				return m, updateUserByAdminCmd(m.token, id, m.adminUserNameInput.Value(), m.adminUserEmailInput.Value(), m.adminUserPassInput.Value())
			}

			if m.screen == screenResetDB {
				m.status = "Réinitialisation en cours..."
				return m, resetDatabaseCmd(m.token)
			}
		}

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

		if m.screen == screenUpdateMovie {
			fields := m.updateFields()
			switch msg.String() {
			case "tab":
				fields[m.updateFieldIndex].Blur()
				m.updateFieldIndex = (m.updateFieldIndex + 1) % len(fields)
				fields[m.updateFieldIndex].Focus()
				return m, nil

			case "shift+tab":
				fields[m.updateFieldIndex].Blur()
				m.updateFieldIndex--
				if m.updateFieldIndex < 0 {
					m.updateFieldIndex = len(fields) - 1
				}
				fields[m.updateFieldIndex].Focus()
				return m, nil
			}
		}

		if m.screen == screenBulkUpdate {
			fields := m.bulkUpdateFields()
			switch msg.String() {
			case "tab":
				fields[m.bulkUpdateFieldIndex].Blur()
				m.bulkUpdateFieldIndex = (m.bulkUpdateFieldIndex + 1) % len(fields)
				fields[m.bulkUpdateFieldIndex].Focus()
				return m, nil

			case "shift+tab":
				fields[m.bulkUpdateFieldIndex].Blur()
				m.bulkUpdateFieldIndex--
				if m.bulkUpdateFieldIndex < 0 {
					m.bulkUpdateFieldIndex = len(fields) - 1
				}
				fields[m.bulkUpdateFieldIndex].Focus()
				return m, nil
			}
		}

		if m.screen == screenRegister {
			fields := m.registerFields()
			switch msg.String() {
			case "tab":
				fields[m.registerFieldIndex].Blur()
				m.registerFieldIndex = (m.registerFieldIndex + 1) % len(fields)
				fields[m.registerFieldIndex].Focus()
				return m, nil

			case "shift+tab":
				fields[m.registerFieldIndex].Blur()
				m.registerFieldIndex--
				if m.registerFieldIndex < 0 {
					m.registerFieldIndex = len(fields) - 1
				}
				fields[m.registerFieldIndex].Focus()
				return m, nil
			}
		}

		if m.screen == screenUpdateOwnProfile {
			fields := m.profileFields()
			switch msg.String() {
			case "tab":
				fields[m.profileFieldIndex].Blur()
				m.profileFieldIndex = (m.profileFieldIndex + 1) % len(fields)
				fields[m.profileFieldIndex].Focus()
				return m, nil

			case "shift+tab":
				fields[m.profileFieldIndex].Blur()
				m.profileFieldIndex--
				if m.profileFieldIndex < 0 {
					m.profileFieldIndex = len(fields) - 1
				}
				fields[m.profileFieldIndex].Focus()
				return m, nil
			}
		}

		if m.screen == screenUpdateUserByAdmin {
			fields := m.adminUserFields()
			switch msg.String() {
			case "tab":
				fields[m.adminUserFieldIndex].Blur()
				m.adminUserFieldIndex = (m.adminUserFieldIndex + 1) % len(fields)
				fields[m.adminUserFieldIndex].Focus()
				return m, nil

			case "shift+tab":
				fields[m.adminUserFieldIndex].Blur()
				m.adminUserFieldIndex--
				if m.adminUserFieldIndex < 0 {
					m.adminUserFieldIndex = len(fields) - 1
				}
				fields[m.adminUserFieldIndex].Focus()
				return m, nil
			}
		}

		if m.screen == screenSettings {
			fields := m.settingsFields()

			switch msg.String() {

			case "tab":
				fields[m.settingsFieldIndex].Blur()
				m.settingsFieldIndex = (m.settingsFieldIndex + 1) % len(fields)
				fields[m.settingsFieldIndex].Focus()
				return m, nil

			case "shift+tab":
				fields[m.settingsFieldIndex].Blur()
				m.settingsFieldIndex--
				if m.settingsFieldIndex < 0 {
					m.settingsFieldIndex = len(fields) - 1
				}
				fields[m.settingsFieldIndex].Focus()
				return m, nil

			case "enter":
				AppConfig.Theme = m.settingsThemeInput.Value()
				AppConfig.APIURL = m.settingsAPIInput.Value()
				AppConfig.Accent = m.settingsAccentInput.Value()

				SaveConfig(AppConfig)

				applyTheme()

				applyAccent()

				API_URL = AppConfig.APIURL

				m.status = "Paramètres enregistrés"
				m.screen = screenMenu
				return m, nil
			}
		}

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
				m.titleInput.SetValue("")
				m.directorInput.SetValue("")
				m.genreInput.SetValue("")
				m.yearInput.SetValue("")
				m.descriptionInput.SetValue("")
				fields := m.createFields()
				fields[0].Focus()
				return m, nil

			case "4":
				m.status = "ID du film à afficher"
				m.screen = screenMovieByID
				m.movieIDInput.SetValue("")
				m.movieIDInput.Focus()
				return m, nil

			case "5":
				m.status = "Mise à jour d'un film"
				m.screen = screenUpdateMovie
				m.updateFieldIndex = 0
				m.updateIDInput.SetValue("")
				m.titleInput.SetValue("")
				m.directorInput.SetValue("")
				m.genreInput.SetValue("")
				m.yearInput.SetValue("")
				m.descriptionInput.SetValue("")
				fields := m.updateFields()
				fields[0].Focus()
				return m, nil

			case "6":
				m.status = "ID du film à supprimer"
				m.screen = screenDeleteMovie
				m.deleteIDInput.SetValue("")
				m.deleteIDInput.Focus()
				return m, nil

			case "7":
				m.status = "IDs à supprimer (séparés par des virgules)"
				m.screen = screenBulkDelete
				m.bulkDeleteIDsInput.SetValue("")
				m.bulkDeleteIDsInput.Focus()
				return m, nil

			case "8":
				m.status = "Mise à jour en lot"
				m.screen = screenBulkUpdate
				m.bulkUpdateFieldIndex = 0
				m.bulkUpdateIDsInput.SetValue("")
				m.titleInput.SetValue("")
				m.directorInput.SetValue("")
				m.genreInput.SetValue("")
				m.yearInput.SetValue("")
				m.descriptionInput.SetValue("")
				fields := m.bulkUpdateFields()
				fields[0].Focus()
				return m, nil

			case "9":
				m.status = "Inscription utilisateur"
				m.screen = screenRegister
				m.registerFieldIndex = 0
				m.registerNameInput.SetValue("")
				m.registerEmailInput.SetValue("")
				m.registerPassInput.SetValue("")
				fields := m.registerFields()
				fields[0].Focus()
				return m, nil

			case "10":
				m.status = "Mettre à jour votre profil"
				m.screen = screenUpdateOwnProfile
				m.profileFieldIndex = 0
				m.profileNameInput.SetValue("")
				m.profileEmailInput.SetValue("")
				m.profilePassInput.SetValue("")
				fields := m.profileFields()
				fields[0].Focus()
				return m, nil

			case "11":
				m.status = "Mettre à jour un utilisateur (admin)"
				m.screen = screenUpdateUserByAdmin
				m.adminUserFieldIndex = 0
				m.adminUserIDInput.SetValue("")
				m.adminUserNameInput.SetValue("")
				m.adminUserEmailInput.SetValue("")
				m.adminUserPassInput.SetValue("")
				fields := m.adminUserFields()
				fields[0].Focus()
				return m, nil

			case "12":
				m.status = "Réinitialiser la base"
				m.screen = screenResetDB
				return m, nil

			case "f4":
				m.status = "Paramètres de l'application"
				m.screen = screenSettings
				m.settingsFieldIndex = 0

				fields := m.settingsFields()
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

	case movieMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur film: %v", msg.err)
			return m, nil
		}
		m.movies = []Movie{msg.movie}
		m.status = fmt.Sprintf("Film #%d : %s", msg.movie.ID, msg.movie.Title)
		m.screen = screenMovies
		return m, nil

	case createMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur création: %v", msg.err)
			return m, nil
		}
		m.status = "Film créé !"
		m.screen = screenMenu
		return m, nil

	case updateMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur mise à jour: %v", msg.err)
			return m, nil
		}
		m.status = "Film mis à jour !"
		m.screen = screenMenu
		return m, nil

	case deleteMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur suppression: %v", msg.err)
			return m, nil
		}
		m.status = "Film supprimé !"
		m.screen = screenMenu
		return m, nil

	case bulkDeleteMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur suppression en lot: %v", msg.err)
			return m, nil
		}
		m.status = "Films supprimés !"
		m.screen = screenMenu
		return m, nil

	case bulkUpdateMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur mise à jour en lot: %v", msg.err)
			return m, nil
		}
		m.status = "Films mis à jour !"
		m.screen = screenMenu
		return m, nil

	case registerMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur inscription: %v", msg.err)
			return m, nil
		}
		m.token = msg.token
		m.status = "Inscription réussie. Vous êtes connecté."
		m.screen = screenMenu
		return m, nil

	case userUpdateMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur mise à jour profil: %v", msg.err)
			return m, nil
		}
		m.status = "Profil mis à jour !"
		m.screen = screenMenu
		return m, nil

	case adminUserUpdateMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur mise à jour utilisateur: %v", msg.err)
			return m, nil
		}
		m.status = "Utilisateur mis à jour !"
		m.screen = screenMenu
		return m, nil

	case resetMsg:
		if msg.err != nil {
			m.status = fmt.Sprintf("Erreur reset BDD: %v", msg.err)
			return m, nil
		}
		m.status = "Base de données réinitialisée !"
		m.screen = screenMenu
		return m, nil
	}

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

	case screenMovieByID:
		var c tea.Cmd
		m.movieIDInput, c = m.movieIDInput.Update(msg)
		return m, c

	case screenDeleteMovie:
		var c tea.Cmd
		m.deleteIDInput, c = m.deleteIDInput.Update(msg)
		return m, c

	case screenBulkDelete:
		var c tea.Cmd
		m.bulkDeleteIDsInput, c = m.bulkDeleteIDsInput.Update(msg)
		return m, c

	case screenCreate:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		fields := m.createFields()
		*fields[m.createFieldIndex], cmd = fields[m.createFieldIndex].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case screenUpdateMovie:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		fields := m.updateFields()
		*fields[m.updateFieldIndex], cmd = fields[m.updateFieldIndex].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case screenBulkUpdate:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		fields := m.bulkUpdateFields()
		*fields[m.bulkUpdateFieldIndex], cmd = fields[m.bulkUpdateFieldIndex].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case screenRegister:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		fields := m.registerFields()
		*fields[m.registerFieldIndex], cmd = fields[m.registerFieldIndex].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case screenUpdateOwnProfile:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		fields := m.profileFields()
		*fields[m.profileFieldIndex], cmd = fields[m.profileFieldIndex].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case screenUpdateUserByAdmin:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		fields := m.adminUserFields()
		*fields[m.adminUserFieldIndex], cmd = fields[m.adminUserFieldIndex].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	case screenSettings:
		var cmds []tea.Cmd
		var cmd tea.Cmd
		fields := m.settingsFields()
		*fields[m.settingsFieldIndex], cmd = fields[m.settingsFieldIndex].Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)

	}

	return m, nil
}

func (m model) View() string {
	switch m.screen {

	case screenLogin:
		return fmt.Sprintf(
			"Login\n\nEmail: %s\nPassword: %s\n\n[enter] pour se connecter\n[f2] pour s'inscrire\n\n%s\n",
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
				"4. Afficher un film par ID\n"+
				"5. Mettre à jour un film\n"+
				"6. Supprimer un film\n"+
				"7. Supprimer plusieurs films\n"+
				"8. Mettre à jour plusieurs films\n"+
				"9. Mettre à jour votre profil\n"+
				"10. Mettre à jour un utilisateur (admin)\n"+
				"11. Réinitialiser la BDD\n"+
				"f4. Paramètres\n"+
				"ctrl+c. Quitter\n\n"+
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

	case screenMovieByID:
		return fmt.Sprintf(
			"Afficher un film\n\nID: %s\n\n[enter] pour charger\n[esc] Menu\n\n%s\n",
			m.movieIDInput.View(),
			m.status,
		)

	case screenUpdateMovie:
		return fmt.Sprintf(
			"Mise à jour d'un film\n\nID: %s\nTitre: %s\nRéalisateur: %s\nGenre: %s\nAnnée: %s\nDescription: %s\n\n[tab] champ suivant, [shift+tab] champ précédent\n[enter] pour mettre à jour\n[esc] Menu\n\n%s\n",
			m.updateIDInput.View(),
			m.titleInput.View(),
			m.directorInput.View(),
			m.genreInput.View(),
			m.yearInput.View(),
			m.descriptionInput.View(),
			m.status,
		)

	case screenDeleteMovie:
		return fmt.Sprintf(
			"Suppression d'un film\n\nID: %s\n\n[enter] pour supprimer\n[esc] Menu\n\n%s\n",
			m.deleteIDInput.View(),
			m.status,
		)

	case screenBulkDelete:
		return fmt.Sprintf(
			"Suppression en lot\n\nIDs: %s\n\n[enter] pour supprimer\n[esc] Menu\n\n%s\n",
			m.bulkDeleteIDsInput.View(),
			m.status,
		)

	case screenBulkUpdate:
		return fmt.Sprintf(
			"Mise à jour en lot\n\nIDs: %s\nTitre: %s\nRéalisateur: %s\nGenre: %s\nAnnée: %s\nDescription: %s\n\n[tab] champ suivant, [shift+tab] champ précédent\n[enter] pour mettre à jour\n[esc] Menu\n\n%s\n",
			m.bulkUpdateIDsInput.View(),
			m.titleInput.View(),
			m.directorInput.View(),
			m.genreInput.View(),
			m.yearInput.View(),
			m.descriptionInput.View(),
			m.status,
		)

	case screenRegister:
		return fmt.Sprintf(
			"Inscription\n\nNom: %s\nEmail: %s\nMot de passe: %s\n\n[tab] champ suivant, [shift+tab] champ précédent\n[enter] pour s'inscrire\n[esc] Menu\n\n%s\n",
			m.registerNameInput.View(),
			m.registerEmailInput.View(),
			m.registerPassInput.View(),
			m.status,
		)

	case screenUpdateOwnProfile:
		return fmt.Sprintf(
			"Mettre à jour votre profil\n\nNom: %s\nEmail: %s\nMot de passe: %s\n\n[tab] champ suivant, [shift+tab] champ précédent\n[enter] pour mettre à jour\n[esc] Menu\n\n%s\n",
			m.profileNameInput.View(),
			m.profileEmailInput.View(),
			m.profilePassInput.View(),
			m.status,
		)

	case screenUpdateUserByAdmin:
		return fmt.Sprintf(
			"Mise à jour d'un utilisateur (admin)\n\nID utilisateur: %s\nNom: %s\nEmail: %s\nMot de passe: %s\n\n[tab] champ suivant, [shift+tab] champ précédent\n[enter] pour mettre à jour\n[esc] Menu\n\n%s\n",
			m.adminUserIDInput.View(),
			m.adminUserNameInput.View(),
			m.adminUserEmailInput.View(),
			m.adminUserPassInput.View(),
			m.status,
		)

	case screenResetDB:
		return fmt.Sprintf(
			"Réinitialisation de la base\n\n[enter] pour confirmer la réinitialisation\n[esc] Menu\n\n%s\n",
			m.status,
		)

	case screenSettings:
		fields := m.settingsFields()
		return fmt.Sprintf(
			styleTitle.Render("=== Paramètres ===")+"\n\n"+
				styleLabel.Render("Thème (light/dark):")+"\n%s\n\n"+
				styleLabel.Render("API URL:")+"\n%s\n\n"+
				styleLabel.Render("Couleur d'accent:")+"\n%s\n\n"+
				"[tab] champ suivant, [shift+tab] champ précédent\n"+
				"[enter] pour enregistrer\n[esc] Menu\n\n%s\n",
			fields[0].View(),
			fields[1].View(),
			fields[2].View(),
			m.status,
		)

	}

	return "Écran inconnu"
}

var AccentColor = "#FF00FF"

func applyTheme() {
	if AppConfig.Theme == "dark" {
		styleTitle = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor)).Bold(true)
		styleLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#CCCCCC"))
		styleInput = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	} else {
		styleTitle = lipgloss.NewStyle().Foreground(lipgloss.Color(AccentColor)).Bold(true)
		styleLabel = lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))
		styleInput = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
	}
}

func applyAccent() {
	AccentColor = AppConfig.Accent
}
