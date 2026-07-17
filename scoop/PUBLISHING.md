# Publication Scoop

Ce dépôt contient un manifest Scoop prêt à adapter pour la CLI TUI `biblio-go`.

## Ce qui est publié
- Binaire Windows de la TUI: `biblio-go.exe`
- Manifest Scoop: `scoop/biblio-go.json`

## Pré-requis
- Un dépôt GitHub public pour les releases
- Un bucket Scoop, soit ton propre bucket, soit un PR vers un bucket existant
- Go installé sur Windows pour produire le binaire
- PowerShell 5.1 ou PowerShell 7

## Build local
Depuis `biblio_go/tui`:

```powershell
$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = "0"
go build -o biblio-go.exe .
```

Créer ensuite une archive ZIP contenant `biblio-go.exe` à la racine de l'archive:

```powershell
Compress-Archive -Path .\biblio-go.exe -DestinationPath .\biblio-go-windows-amd64.zip -Force
```

Calculer le hash SHA256:

```powershell
Get-FileHash .\biblio-go-windows-amd64.zip -Algorithm SHA256
```

## Mise à jour du manifest
Remplacer dans `scoop/biblio-go.json`:
- `<OWNER>` par ton utilisateur GitHub
- `<REPO>` par le nom du dépôt
- `REPLACE_WITH_SHA256` par le hash calculé
- la version `0.1.0` par la version publiée

## Publication
1. Créer une release GitHub avec l'archive `biblio-go-windows-amd64.zip`.
2. Copier le manifest dans le bucket Scoop cible.
3. Si c'est ton bucket, pousser le commit puis l'indexer normalement.
4. Si c'est un bucket public, ouvrir une PR avec le manifest.

## Test d'installation locale
Ajouter le bucket local ou tester le manifest directement:

```powershell
scoop install .\scoop\biblio-go.json
```

Puis lancer:

```powershell
biblio-go
```

## Validation fonctionnelle
- L'écran Paramètres doit permettre de changer `theme`, `api_url` et `accent`.
- Laisser un champ vide doit conserver la valeur existante.
- Un thème invalide doit afficher une erreur et ne rien enregistrer.
- Après sauvegarde, la configuration doit être persistée dans `tui/config/config.json`.
