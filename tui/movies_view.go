package main

import (
	"fmt"
	"strings"
)

func renderMoviesView(movies []Movie, status string) string {
	var b strings.Builder

	b.WriteString("Films (API biblio_go)\n\n")

	if len(movies) == 0 {
		b.WriteString("Aucun film.\n\n")
	} else {
		for _, m := range movies {
			b.WriteString(fmt.Sprintf(
				"#%d - %s (%d)\n   Réalisateur: %s | Genre: %s\n   %s\n\n",
				m.ID,
				m.Title,
				m.ReleaseYear,
				m.Director,
				m.Genre,
				m.Description,
			))
		}
	}

	b.WriteString(fmt.Sprintf("\n%s\n\n[ctrl+c] pour quitter\n", status))
	b.WriteString(fmt.Sprintf("\n%s\n\n[esc] Retour au menu depuis n'importe où\n", status))

	return b.String()
}
